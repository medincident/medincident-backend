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
- gRPC-сервис (pure gRPC; аннотации `google.api.http` в proto — под
  будущую REST-обёртку, для неё генерятся grpc-gateway stubs)
- gorm v2 + gorm.io/driver/postgres
- `github.com/guregu/null/v6` (только в `internal/model`)
- samber/do/v2 · samber/oops · zerolog
- dbmate (`go tool dbmate`) для миграций
- buf (`go tool buf`) + protoc-gen-go + protoc-gen-go-grpc +
  protoc-gen-grpc-gateway + protoc-gen-openapiv2 + protoc-gen-doc
  для кодогенерации

## Архитектура

- **gRPC API** — `OrgStructureService` в `service.orgstructure.v1`,
  `MembershipService` в `service.membership.v1`,
  `IncidentClassifierService` в `service.incident.classifier.v1`.
  Proto-контракты живут в [`api/proto/`](api/proto/) внутри этого
  репо; Go-биндинги, merged OpenAPI и Markdown-документация генерятся
  через `task gen` и коммитятся в [`pkg/`](pkg/),
  [`api/openapi/`](api/openapi/) и [`docs/proto/`](docs/proto/).
- **Плоский layout, без DDD.** Три service-струт типа
  (`OrganizationService`, `ClinicService`, `DepartmentService`) в одном
  пакете `internal/services/orgstructure`. Зависимости —
  `*gorm.DB` + `*zerolog.Logger`, никаких интерфейсов и репозиториев.
- **Файл на метод.** `organization_create.go`, `clinic_update_details.go`
  и т.д. — одна команда, её валидация, её event builder, её subject
  и её транзакционный flow в одном файле.
- **Транзакционный outbox.** Доменные события пишутся в `outbox.events`
  в той же транзакции, что и мутация агрегата, через
  `outbox.Publish` из [`internal/service/outbox`](internal/service/outbox).
  Схема outbox-таблицы: `id, subject, payload, headers, created_at,
  published_at` — больше ничего.
  `payload` — это сериализованный `event.v1.Envelope`
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

- `api/proto/` — исходные `.proto` контракты (event/* + service/*)
- `api/openapi/command-service.swagger.json` — merged OpenAPI v2 (коммитится)
- `pkg/` — сгенерированный buf Go-код (коммитится)
- `docs/proto/command-service.md` — сгенерированная Markdown-документация (коммитится)
- `cmd/command-server/` — точка входа command-side gRPC сервера, graceful shutdown
- `cmd/query-server/` — точка входа query-side сервера (placeholder, заполняется в Plan 3)
- `configs/` — YAML config пример
- `db/migrations/` — dbmate миграции (через `task migrate:new`)
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
go run ./cmd/command-server --config configs/command-server.example.yaml
```

## Тестирование

```bash
task test:unit          # unit-тесты (быстро, без Docker)
task test:integration   # testcontainers-backed integration (требует Docker)
task test               # всё вместе
```

## Генерация кода

```bash
task gen         # buf generate (pkg/, api/openapi/) + docs template
task gen:check   # verify pkg/, api/openapi/, docs/proto/ в синке с proto

task proto:fmt          # форматирование .proto
task proto:fmt:check    # dry-run
task proto:lint         # buf lint
task proto:breaking     # buf breaking vs origin/main
```

## Линт и безопасность

```bash
task fmt         # форматирование через golangci-lint
task fmt:check   # dry-run
task lint        # golangci-lint
task vuln        # govulncheck
```
