# medincident-command-service

Write-сторона CQRS-развертки medincident. Хостит всю доменную бизнес-логику
платформы: создание и изменение агрегатов, инварианты, публикацию доменных
событий. Identity живёт в Zitadel IdP; этот сервис хранит только тонкий
app-level слой поверх пользователей (роли, флаги) плюс свои собственные
агрегаты (сейчас — organization; далее — clinic, department, patient).

Read-сторона — отдельный сервис `medincident-query-service`, который
подписан на JetStream subject'ы и строит собственные проекции. Command
никогда не читает проекции и не знает о query-стороне.

## Архитектура

- **gRPC API** — транспорт команд. Proto-контракты живут в
  `medincident-proto`, Go-биндинги генерятся через `task gen` и
  коммитятся в [`gen/`](gen/).
- **Package-by-feature** — каждый агрегат владеет директорией
  `internal/<bounded-context>/<aggregate>/` с доменом (root, events,
  валидаторы), `app/` (Service, commands, Repository interface) и
  `infra/` (outbox proto мапперы). Никаких layer-by-layer слоёв.
- **Транзакционный outbox** — доменные события пишутся в
  `outbox.events` в той же транзакции, что и мутация агрегата, через
  `tx.Within`. Отдельный outbox-publisher (будет позже) читает строки
  через LISTEN/NOTIFY + polling fallback и публикует их в NATS
  JetStream, обёрнутыми в `medincident.events.v1.Envelope`.
- **Domain events carry NEW state only.** `Renamed.Name`, не
  `Renamed.OldName`+`Renamed.NewName`. Диф считают consumer'ы из своих
  собственных проекций.
- **`aggregate.Root`** embed даёт `CreatedAt`/`UpdatedAt`/pending events
  и метод `Raise(event, now)` — единственное легальное место, где
  `UpdatedAt` двигается.
- **Correlation id** — request-scoped строка в `context.Context` (пакет
  `internal/correlation`), пишется в каждую outbox-строку, чтобы
  publisher мог проставить `envelope.correlation_id`.
- **samber/do/v2 DI** — весь граф собирается в `internal/di/` и
  инвокается из [`cmd/server/main.go`](cmd/server/main.go). Eager
  invocation при старте падает быстро на любой ошибке проводки.

Полный дизайн — в
[`docs/superpowers/specs/2026-04-10-command-service-base-design.md`](docs/superpowers/specs/2026-04-10-command-service-base-design.md).

## Раскладка

```
cmd/server/                            — entry point, graceful shutdown
configs/                               — YAML config reference + .env.example
db/migrations/                         — dbmate SQL миграции (НИКОГДА не руками)
gen/                                   — buf-generated Go биндинги (закоммичены)
internal/
  config/                              — конфиг DTO + YAML loader
  correlation/                         — request-scoped correlation id
  di/                                  — samber/do providers
  orgstructure/                        — bounded context org-структуры
    organization/
      domain/                          — root, events, валидаторы, VO
      app/                             — Service, commands, Repository iface
      infra/                           — outbox proto мапперы
  outbox/                              — generic outbox: Store, Registry, Publish
  shared/
    aggregate/                         — base Root с CreatedAt/UpdatedAt/events
    geo/                               — Point, Address value objects
  storage/postgres/                    — единственное место импорта pgx
  tx/                                  — tx.Tx, tx.Beginner, tx.Within
test/integration/                      — интеграционные тесты (build tag)
```

## Стек

Go 1.26 · gRPC · pgx/v5 · sqlc · Masterminds/squirrel · dbmate ·
NATS JetStream · samber/do/v2 · samber/oops · zerolog · google/uuid (v7) ·
testcontainers-go.

Все Go-тулы (buf, golangci-lint, sqlc, dbmate, govulncheck,
protoc-gen-go) запинены через `tool` директиву в `go.mod` и вызываются
через `go tool <name>`. Никаких `brew install` кроме самого Go, Docker
(для интеграционных тестов) и `task`.

## Quickstart

```bash
cp configs/.env.example .env
task migrate           # dbmate up
task test              # unit + integration
task run               # go run ./cmd/server -config configs/config.example.yaml
```

## Сборка

Для сборки бинарника используется Makefile — он знает про `./cmd/server`
и умеет кросс-компилить под все платформы из списка `PLATFORMS`:

```bash
make build             # сборка под текущую хост-платформу в ./dist/
make build-all         # кросс-компиляция под все платформы из PLATFORMS
make run ARGS="-config configs/config.example.yaml"
make clean             # снести ./dist
```

Список платформ редактируется прямо в [`Makefile`](Makefile) (`PLATFORMS`).
По умолчанию собирает Linux/macOS/Windows под amd64, arm64 и 386.

## Contributing

См. [CONTRIBUTING.md](CONTRIBUTING.md) — требования, установка, правила
редактирования кода (миграции, oops, valueobjects) и доступные task'и.
