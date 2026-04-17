# Contributing

## Требования

- [Go 1.26+](https://go.dev/dl/) — `buf`, `golangci-lint`, `dbmate`, `govulncheck`, `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`, `protoc-gen-openapiv2` и `protoc-gen-doc` запинены через `go tool`, отдельная установка не нужна.
- [Task](https://taskfile.dev/installation/) — task runner, обёртка над всеми частыми командами.
- [Docker](https://docs.docker.com/get-docker/) — для `task test:integration` (testcontainers-go поднимает Postgres 16).
- `make` — для сборки бинарника (`make build`, `make build-all`).

Никакого C-toolchain или отдельно стоящего `protoc` не нужно — всё идёт через `go tool`.

## Старт

Склонируй репу:

```bash
git clone https://github.com/medincident/medincident-command-service.git
cd medincident-command-service
```

Настрой окружение и проверь тулчейн:

```bash
cp configs/.env.example .env
task migrate        # dbmate up
task lint           # golangci-lint
task test:unit      # быстрые юнит-тесты
```

Все три должны завершиться с кодом 0 на чистом чекауте.

## Частые команды

```bash
task gen              # кодген: pkg/ (go+grpc+gateway) + api/openapi/ + docs/proto/
task proto:fmt        # форматирование .proto
task proto:lint       # buf lint
task proto:breaking   # buf breaking vs origin/main
task lint             # go tool golangci-lint run
task fmt              # форматирование через golangci-lint formatters
task fmt:check        # проверка формата (dry-run, без модификаций)
task test:unit        # юнит-тесты (race, count=1)
task test:integration # интеграционные тесты (build tag, требует Docker)
task test             # unit + integration
task vuln             # govulncheck
task migrate          # dbmate up
task migrate:new -- <name>  # создать новую миграцию через dbmate new
```

`task gen` — зонтичная точка входа для всех кодгенов. Сегодня дёргает `buf generate` (default template — `pkg/` + `api/openapi/`) и `buf generate --template buf.gen.docs.yaml` (`docs/proto/`); любой новый генератор добавится сюда.

## Сборка бинарника

Makefile знает про entry point `./cmd/command-server` и умеет кросс-компилить.

```bash
make build                             # под текущую хост-платформу → ./dist/
make build-all                         # все платформы из PLATFORMS
make run ARGS="-config configs/command-server.example.yaml"
make clean                             # снести ./dist
```

Список целевых платформ редактируется в самом [`Makefile`](Makefile) в переменной `PLATFORMS`. По умолчанию: Linux (amd64/arm64/386), macOS (amd64/arm64), Windows (amd64/arm64/386). Бинарники именуются `command-server-<os>-<arch>` (с `.exe` на Windows).

## Редактирование кода — обязательные правила

Эти правила не декоративные. Часть из них отражает уже пойманные инциденты, часть — архитектурные инварианты, на которые опирается остальной код. Полный список живёт в [AGENTS.md](AGENTS.md); ниже — ключевые для новых контрибьютеров.

### Миграции

Файлы в `db/migrations/` **всегда** создаются через:

```bash
task migrate:new -- <snake_case_name>
```

Это вызывает `go tool dbmate new <name>`. Никогда не создавай migration-файлы редактором или через `touch`/`>` — dbmate владеет форматом имени и структурой заголовка.

### Ошибки

Используй `samber/oops`:

```go
return oops.In("orgstructure.organization").
    Code(ErrCodeOrganizationNameEmpty).
    Public("organization name is required").
    With("field", "name").
    Wrap(err)
```

- Никакого `fmt.Errorf("%w", ...)` и `errors.New` для доменных ошибок. `errors.New` допустим только в тестах для sentinel-сравнений.
- `oops.Code(...)` принимает **только package-level константу**, объявленную **в том же файле**, где код эмитится. Generic'и типа `CodeInvalidArgument` запрещены. Общий файл `codes.go` разрешён только когда два producer'а в разных файлах одного пакета делят один код (редко).
- **Не заворачивай `errors.Join(...)` в oops.** `oops` реализует `Unwrap() error` (singular), и gRPC-интерсептор теряет индивидуальные field-violation'ы при flatten'е. Возвращай joined error сырым.

### Value Objects, агрегаты, entities

- VO: **один** конструктор, никаких мутирующих методов. Пример: `geo.NewAddress(...)`.
- Агрегаты/entities: **два** конструктора — `New` (валидирует, raise `Created`) и `Hydrate` (trusted, без валидации и без events). Никаких экспортированных `Validate()`.
- Агрегаты embed'ят `internal/shared/aggregate.Root` для `CreatedAt`/`UpdatedAt`/`events`. `Raise(event, now)` — единственное место, где `UpdatedAt` двигается. Никогда не пиши `o.UpdatedAt = now` напрямую.

### Инварианты

Все max/min/threshold — именованные package-level константы (`maxNameLen`, `MinLongitude`, …). Никаких инлайн-литералов в валидаторах. `With()` тоже ссылается на константу, не на литерал.

### Multi-error валидация

Внутри одной операции (например `Service.Create`) валидатор прогоняет **все** VO и конструкторы, которые может, и собирает все ошибки через `errors.Join`. Клиент должен получить все нарушения в одном ответе, а не первое попавшееся. Fail-fast разрешён только когда следующий шаг физически не может запуститься без валидного предыдущего значения — и даже тогда пробуй пробросить `nil`/zero, если это семантически валидно.

### Application services

Сервисы принимают `XCommand` и возвращают `XResult` — структуры, не позиционные параметры. Пример: `Service.Create(ctx, CreateOrganizationCommand{...}) (CreateOrganizationResult, error)`, а не `Service.Create(ctx, name, desc, addr)`. Позиционные параметры оставлены для доменных конструкторов.

### Domain types и транспорт

- Доменные типы **не несут struct tag'ов**. Proto-теги и JSON-теги — это транспортный concern. Outbox marshal'ит доменные события через `encoding/json` по exported field names напрямую, так что переименование поля доменного события — breaking change для уже записанных outbox-строк.
- Proto — это transport-only. Импорт proto в feature-дереве разрешён только в `internal/<bc>/<agg>/infra/outbox.go` (мапперы) и будущем `internal/grpcapi/`.
- Domain events несут **только новое состояние**. Не `OldX` + `NewX`. Consumer'ы считают диф из своих проекций.

### Go-тулинг

Никогда не ставь Go-тулы глобально. Всё добавляется через `go get -tool <module>@latest` и вызывается как `go tool <name>` в `Taskfile.yml`.

## Тестовый layout

- **Unit-тесты** — рядом с кодом (`*_test.go` в том же пакете).
- **Интеграционные тесты** — под `test/integration/` с `//go:build integration`. Поднимают реальный Postgres 16 через testcontainers-go, применяют миграции через `go tool dbmate`, гоняют настоящий Service / Repository / outbox pipeline.

```bash
task test:unit          # быстрый цикл
task test:integration   # полный стек, требует Docker
task test               # unit + integration
```

## Git workflow

- **Никогда не `git push` без явной команды пользователя.** Всё живёт в локальных коммитах на feature-ветке до тех пор, пока тебя не попросят запушить.
- Миграции — через `task migrate:new`, никогда руками.
- Перед PR: `task lint`, `task fmt:check`, `task test`, `task vuln` — всё зелёное.
