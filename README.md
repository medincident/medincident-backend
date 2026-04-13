# medincident-command-service

Write-сторона CQRS-развёртки medincident. Обслуживает три агрегата
организационной структуры: **Organization**, **Clinic**, **Department**.
Каждая мутация (Create / Update) атомарно пишет в `domain.*` таблицу и
складывает событие в `outbox.events` в той же транзакции. Отдельный
publisher-сервис драйнит outbox и публикует envelope'ы в NATS JetStream.

Read-сторона — отдельный сервис `medincident-query-service`, который
подписан на JetStream subject'ы и строит собственные проекции.
Command никогда не читает проекции и не знает о query-стороне.

## Стек

- Go 1.26
- gRPC-сервис (pure gRPC, без REST gateway — аннотации `google.api.http`
  оставлены в proto для будущей REST-обёртки)
- gorm v2 + gorm.io/driver/postgres
- `github.com/guregu/null/v6` (только в `internal/model`)
- samber/do/v2 · samber/oops · zerolog
- dbmate (`go tool dbmate`) для миграций
- buf (`go tool buf`) + protoc-gen-go + protoc-gen-go-grpc для кодогенерации

## Архитектура

- **gRPC API** — `OrgStructureService` в `medincident.service.orgstructure.v1`.
  Proto-контракты живут в отдельном репо `medincident-proto`; Go-биндинги
  генерятся через `task gen` и коммитятся в [`gen/api/medincident/`](gen/api/medincident/).
- **Плоский layout, без DDD.** Три service-струт типа
  (`OrganizationService`, `ClinicService`, `DepartmentService`) в одном
  пакете `internal/services/orgstructure`. Зависимости —
  `*gorm.DB` + `*zerolog.Logger`, никаких интерфейсов и репозиториев.
- **Файл на метод.** `organization_create.go`, `clinic_update_details.go`
  и т.д. — одна команда, её валидация, её event builder, её subject
  и её транзакционный flow в одном файле.
- **Транзакционный outbox.** Доменные события пишутся в `outbox.events`
  в той же транзакции, что и мутация агрегата, через
  `orgstructure.AppendOutboxEvent`. Схема outbox-таблицы: `id, subject,
  payload, headers, created_at, published_at` — больше ничего.
  `payload` — это сериализованный `medincident.event.v1.Envelope`
  внутри которого `google.protobuf.Any` с доменным событием.
- **Events carry NEW state only.** `OrganizationDetailsChanged.Name` —
  это новое имя; consumer'ы считают diff против своих собственных
  проекций.
- **Gorm field-level permissions.** Каждое поле модели имеет явный
  `gorm:"<-:create"` / `gorm:"<-"` / `gorm:"-"` тег, фиксирующий
  immutability на уровне persistence. Первичные ключи и родительские
  FK (`Clinic.OrganizationID`, `Department.ClinicID`) — create-only,
  потому что связь с родителем никогда не меняется.
- **Валидация на входе в сервис.** `errors.Join` собирает все
  field-violations, каждая — отдельный `oops.Code`. БД хранит только
  NOT NULL + PK + FK, без CHECK-констрейнтов.

## Директории

- `cmd/server/` — точка входа, graceful shutdown
- `configs/` — YAML config пример
- `db/migrations/` — dbmate миграции (через `task migrate:new`)
- `gen/` — сгенерированный buf-кодом protobuf (коммитится)
- `internal/config/` — YAML loader + go-playground/validator
- `internal/model/` — gorm-модели (единственное место где живёт `null.X`)
- `internal/services/orgstructure/` — бизнес-логика, файл на метод
- `internal/handler/orgstructure/` — gRPC handler, файл на RPC
- `internal/di/` — samber/do/v2 фабрики (postgres, services, handler, grpc)
- `test/integration/orgstructure/` — testcontainers-backed integration suite

## Запуск локально

Требуется Postgres 17+ и переменная `DATABASE_URL`:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/medincident?sslmode=disable"
task migrate
go run ./cmd/server --config configs/config.example.yaml
```

## Тестирование

```bash
task test:unit          # unit-тесты (быстро, без Docker)
task test:integration   # testcontainers-backed integration (требует Docker)
task test               # всё вместе
```

## Генерация кода

```bash
task gen         # buf generate из medincident-proto (git_repo input)
task gen:check   # verify gen/ в синке с proto
```

## Линт и безопасность

```bash
task fmt         # форматирование через golangci-lint
task fmt:check   # dry-run
task lint        # golangci-lint
task vuln        # govulncheck
```
