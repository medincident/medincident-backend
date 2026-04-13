# Agent Instructions — medincident-command-service

gRPC command side of a CQRS split. Go 1.26 · gorm v2 · guregu/null/v6 ·
samber/do/v2 · samber/oops · zerolog · dbmate · buf. Design lives in
`docs/superpowers/specs/2026-04-13-command-service-simplification-design.md`.

## Hard rules — NOT NEGOTIABLE

1. **Never push.** All work lands on local commits in a feature branch.
   Never `git push` without explicit user instruction.
2. **Never hand-create migration files.** Always run
   `task migrate:new -- <name>` (which calls `go tool dbmate new`).
   Never use `Write`/`Edit` to create a file under `db/migrations/`.
3. **Never use string literals in `oops.Code(...)` calls, and never
   use generic codes.** Every Code is a specific package-level
   constant (e.g. `ErrCodeAddressTextEmpty`, not `CodeInvalidArgument`)
   declared in the same file as the code that emits it.
4. **Never install Go tools globally.** Every Go tool (buf, dbmate,
   golangci-lint, govulncheck, protoc-gen-go, protoc-gen-go-grpc) is
   pinned in `go.mod`'s `tool` directive and invoked via `go tool
   <name>`. Adding a new tool: `go get -tool <module>@latest`.
5. **Never inline invariant limits.** Every max/min/threshold is a
   named package-level constant (`organizationMaxNameLen`,
   `minLongitude`, …). `With()` clauses also reference the constant,
   never the literal.
6. **Never use `fmt.Errorf("%w", ...)` for domain errors.** Use
   `samber/oops`:
   `oops.In("pkg").Code(...).Public("…").With("field", x).Wrap(err)`.
   `errors.New` is OK in tests for sentinel comparisons.
7. **No interfaces for services and no repository layer.** Services
   are concrete struct types that take `*gorm.DB` and `*zerolog.Logger`
   in the constructor. Call gorm directly, open transactions inline
   via `s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error
   {...})`.
8. **Commands and results use plain Go types only.** No `null.X`
   anywhere outside `internal/model`. Optional fields in commands are
   `*string` / `*float64` / `*PointInput`. Conversion to `null.X` at
   the service boundary is inline via `null.StringFromPtr`,
   `null.FloatFrom`, etc. — no local wrapper helpers.
9. **Validation at the service boundary, multi-error.** Every command
   method collects field errors via `errors.Join(errs...)` raw (never
   `oops.Wrap(errors.Join(...))`). Validators are pure functions in
   the same package; each failure leaf carries its own oops `Code`.
   No DB `CHECK` constraints — the DB has only NOT NULL, PK, FK.
10. **Events are built inline with named `buildXxxEvent` functions.**
    No generic mappers. Each event has one builder that knows its
    exact proto type and assembles it inline. Each proto event file
    duplicates its own `Address`/`Point` messages so aggregates evolve
    independently.
11. **Outbox writes go through `AppendOutboxEvent(tx, subject,
    envelope, headers)`.** The caller builds the
    `*medincident.event.v1.Envelope` explicitly. The outbox table has
    exactly `id, subject, payload, headers, created_at, published_at`
    — nothing more. Command-service never writes or reads
    `published_at`; the publisher service owns that column.
12. **Events carry NEW state only.**
    `OrganizationDetailsChanged.Name` is the new name; consumers
    compute diffs from their own prior projection.
13. **DI factories no healthcheck.** Providers only wire; no Ping,
    warm-up, or probe inside the provide function.
14. **Every gorm model field carries an explicit field-level
    permission tag** (`<-:create`, `<-`, `-`). Primary keys, creation
    timestamps, and parent FKs (`Clinic.OrganizationID`,
    `Department.ClinicID`) are `<-:create`. Mutable body fields use
    `<-` explicitly. Parent FKs are immutable — a clinic always
    belongs to exactly one organization and never changes parents;
    a department always belongs to exactly one clinic.

## Directory layout

```
cmd/server/main.go                       — entry point, graceful shutdown
internal/
  config/                                — YAML + go-playground/validator
  di/                                    — samber/do/v2 providers (all factories)
    container.go                         — NewContainer + do.Provide wiring
    zerolog.go                           — logger construction
    postgres.go                          — *gorm.DB provider + Shutdown hook
    grpc.go                              — *grpc.Server wrapper + GracefulStop hook
    services.go                          — three service providers
    handler.go                           — OrgStructureHandler provider
  model/                                 — gorm models. ONLY place `null.X` lives.
  services/orgstructure/                 — business logic, one file per method
    service.go                           — three service struct types + constructors
    outbox.go                            — AppendOutboxEvent helper
    address.go                           — shared Address/Point validators
    organization_{create,update_details,update_legal_address}.go
    clinic_{create,update_details,update_physical_address}.go
    department_{create,update_details}.go
  handler/orgstructure/                  — gRPC handler, one file per RPC
    handler.go                           — OrgStructureHandler struct + proto conversion helpers
    organization_{create,update_details,update_legal_address}.go
    clinic_{create,update_details,update_physical_address}.go
    department_{create,update_details}.go
gen/api/medincident/                     — buf-generated proto (committed)
  event/v1/                              — Envelope
  event/organization/v1/                 — Organization events + own Address/Point
  event/clinic/v1/                       — Clinic events + own Address/Point
  event/department/v1/                   — Department events (no Address)
  service/orgstructure/v1/               — OrgStructureService gRPC + Request/Response
db/migrations/                           — dbmate migrations (never hand-written)
test/integration/orgstructure/           — testcontainers-backed integration suite
configs/                                 — config.example.yaml
```

## Subject scheme

Outbox subjects follow `medincident.event.<aggregate>.v1.<action>`:
- `medincident.event.organization.v1.{created,details_changed,legal_address_changed}`
- `medincident.event.clinic.v1.{created,details_changed,physical_address_changed}`
- `medincident.event.department.v1.{created,details_changed}`

## Tooling

All via `Taskfile.yml`. Key commands:

- `task gen` — `go tool buf generate` (no sqlc)
- `task gen:check` — verify `gen/` is in sync with proto
- `task fmt`, `task fmt:check` — gofumpt + goimports via golangci-lint
- `task lint` — golangci-lint
- `task vuln` — govulncheck
- `task test:unit` — unit tests
- `task test:integration` — testcontainers-backed integration suite
- `task test` — unit + integration
- `task migrate` — apply dbmate migrations
- `task migrate:new -- <name>` — create a new migration file pair
