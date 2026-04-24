# Agent Instructions — medincident-backend

gRPC command side of a CQRS split. Go 1.26 · gorm v2 · guregu/null/v6 ·
samber/oops · zerolog · dbmate · buf. No DI framework — each binary's
main.go wires constructors explicitly and drives teardown via `defer`.
Design lives in `docs/superpowers/specs/2026-04-13-command-service-simplification-design.md`.

## Hard rules — NOT NEGOTIABLE

1. **Never push.** All work lands on local commits in a feature branch.
   Never `git push` without explicit user instruction.
2. **Never hand-create migration files.** Always run
   `task migrate:new -- <name>` (which calls `go tool dbmate new`).
   Never use `Write`/`Edit` to create a file under `db/migrations/`.
3. **Never use string literals in `oops.Code(...)` calls.** Every
   aggregate-specific Code is a package-level constant (e.g.
   `ErrCodeEmployeeAlreadyHired`, not `CodeAlreadyExists`) declared
   in the same file as the code that emits it. The sole exceptions
   are the generic validation codes emitted by
   `internal/service/validation` — `string_required`,
   `string_too_short`, `string_too_long`, `uuid_required`,
   `time_required`, `int_out_of_range`, `float_out_of_range` —
   exported as `validation.CodeXxx` package constants and produced
   by the translator, never hand-written.
4. **Never install Go tools globally.** Every Go tool (buf, dbmate,
   golangci-lint, govulncheck, protoc-gen-go, protoc-gen-go-grpc,
   protoc-gen-grpc-gateway, protoc-gen-openapiv2, protoc-gen-doc) is
   pinned in `go.mod`'s `tool` directive and invoked via `go tool
   <name>`. Adding a new tool: `go get -tool <module>@latest`.
5. **Never inline invariant limits in executable code.** Every
   max/min/threshold used in runtime logic (`With()` clauses,
   comparisons, `categoryDepth(...)` bounds, etc.) is a named
   package-level constant (`incidentClassifierMaxDepth`, …), never
   the literal. The one documented exception is **struct-tag
   validators** (`validate:"min=4,max=256"`) on command structs: the
   numeric bounds live inline because struct tags are compile-time
   strings and the tag is the contract a reviewer reads next to the
   field it constrains. `With()`, `if`, and any other executable
   reference to the same bound still uses a constant.
6. **Never use `fmt.Errorf("%w", ...)` for domain errors.** Use
   `samber/oops`:
   `oops.In("pkg").Code(...).Public("…").With("field", x).Wrap(err)`.
   `errors.New` is OK in tests for sentinel comparisons.
7. **No interfaces for services and no repository layer.** Services
   are concrete struct types that take `*gorm.DB`, `*authz.Authz`,
   and `*zerolog.Logger` in the constructor. Call gorm directly,
   open transactions inline via
   `s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {...})`.
8. **Commands are `Caller + Payload` DTOs carrying primitives only.**
   Every write command is a struct with two fields: `Caller authz.Caller`
   (the trusted identity set by the authn interceptor) and `Payload
   XxxPayload` (the validated client-facing request). The Payload holds
   **primitives only** — `string` for every UUID (validated with
   `validate:"required,uuid"`), `string` for names, `*string` for
   optional strings, `time.Time` / `*time.Time` for timestamps. Rich
   types (`uuid.UUID`, `time.Time`) are parsed inside the service via
   `uuid.MustParse(...)` after validation, then handed to the domain
   model. No `null.X` anywhere outside `internal/model`; conversion at
   the service boundary is inline via `null.StringFrom(strings.TrimSpace(...))`,
   `null.ValueFrom(...)`, etc. — no local wrapper helpers. Handlers
   become pure translators: they extract `callerID` via
   `grpcmw.CallerID(ctx)`, wrap it in `authz.Caller{ZitadelUserID: ...}`,
   and pass proto strings straight through to the Payload.
9. **Each service method runs validate → parse → authorize → execute.**
   The first line of every write is `validation.Struct(cmd.Payload)`;
   the translator walks every `validator.FieldError` and returns an
   `errors.Join` of oops leaves (each with its own `Code`, `"field"`
   path, rule, and public message) — never
   `oops.Wrap(errors.Join(...))`. Immediately after validation the
   service parses Payload UUIDs into local variables, then calls
   `s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.X(scopeID))`
   (or `authz.SystemAdmin` for scope-less operations). Only then does
   business logic run. Payload fields carry their own rules inline as
   `validate:"required,min=4,max=256"` / `validate:"omitnil,..."` tags.
   Every string field is trimmed before validation so whitespace-only
   input fails `required` / `min`. No DB `CHECK` constraints — the DB
   has only NOT NULL, PK, FK.
10. **Events are built inline with named `buildXxxEvent` functions.**
    No generic mappers. Each event has one builder that knows its
    exact proto type and assembles it inline. Each proto event file
    duplicates its own `Address`/`Point` messages so aggregates evolve
    independently.
11. **Outbox writes go through `outbox.AppendEvent(tx, subject,
    envelope)`.** The caller builds the
    `*event.v1.Envelope` explicitly. The outbox table has
    `id, subject, payload, headers, created_at, published_at` on the
    DB side, but the command service only writes `subject` and
    `payload` — the `headers` column is owned by the publisher
    service (dedup keys derived from `outbox.events.id`) and
    `published_at` is the publisher's bookkeeping column.
12. **Events carry NEW state only.**
    `OrganizationDetailsChanged.Name` is the new name; consumers
    compute diffs from their own prior projection.
13. **No DI framework.** Each binary's `main.go` wires dependencies
    explicitly via direct constructor calls and tears them down with
    `defer`. Shared startup helpers (Postgres open, zerolog build,
    Zitadel authorizer) live in `internal/bootstrap/`. No samber/do,
    no wire, no injector — compile-time wiring only. Bootstrap
    helpers must not call Ping, warm-up, or readiness probes —
    construction must return immediately so an unreachable external
    during a rolling deploy cannot hang the boot sequence past the
    k8s pod-termination grace period.
14. **Every gorm model field carries an explicit field-level
    permission tag** (`<-:create`, `<-`, `-`). Primary keys, creation
    timestamps, and parent FKs (`Clinic.OrganizationID`,
    `Department.ClinicID`) are `<-:create`. Mutable body fields use
    `<-` explicitly. Parent FKs are immutable — a clinic always
    belongs to exactly one organization and never changes parents;
    a department always belongs to exactly one clinic.
15. **Build tooling split — Makefile vs Taskfile.**
    - `Makefile` owns **binary compilation only**: `make build`,
      `make build-<binary>`, `make build-all`,
      `make build-all-<binary>`, `make clean`. Never add non-build
      targets to the Makefile.
    - `Taskfile.yml` owns **everything else**: codegen, lint, format,
      tests, vuln scan, proto tooling, DB migrations, Docker image
      builds (`docker:*`). Never add Go-binary compilation tasks
      to the Taskfile — they belong in the Makefile.

## Directory layout

```
cmd/command-server/
  main.go                                — explicit constructor wiring + graceful shutdown
  config.go                              — Config struct + readConfig (package main)
cmd/query-server/
  main.go                                — readers + NATS identity consumer + gRPC server
  config.go                              — Config struct + readConfig (package main)
cmd/gateway-server/
  main.go                                — two ClientConns + grpc-gateway mux + http.Server
  config.go                              — Config struct + readConfig (package main)
internal/
  bootstrap/                             — shared startup helpers (no DI framework)
    zerolog.go                           — BuildZerolog(cfg) → *zerolog.Logger + cleanup
    postgres.go                          — OpenPostgres(ctx, cfg, logger) → *gorm.DB + cleanup
    zitadel.go                           — NewZitadelAuthorizer + NewZitadelService
  config/                                — shared config subtypes only (binary-specific
                                            top-level configs live in cmd/<name>/config.go)
    config.go                            — GRPCServerConfig / PostgresConfig / ZitadelConfig
                                            + ReadAndValidate (YAML + env-expand + validator)
    zerolog.go                           — ZerologConfig subtree
    nats.go                              — NATSConfig (query-server uses it)
  model/                                 — gorm models. ONLY place `null.X` lives.
  middleware/
    grpcmw/                              — gRPC interceptors
      authn.go                           — unary authn interceptor (Zitadel JWT)
      caller.go                          — caller-ID extraction
      error.go                           — error → gRPC status mapping
    httpmw/                              — HTTP middleware
      access_log.go                      — access log via rs/zerolog/hlog
      cors.go                            — CORS wrapper over rs/cors
  service/
    authz/                               — role-based check helpers (shared)
    zitadel/                             — Zitadel client (today: user verify)
    command/orgstructure/                — write-side business logic, one file per method
      service.go                         — three service struct types + constructors
      address.go                         — shared Address/Point validators
      organization_{create,update_details,update_legal_address}.go
      clinic_{create,update_details,update_physical_address}.go
      department_{create,update_details}.go
    command/membership/                  — write-side roles/employees/vacations
    command/incident/classifier/         — write-side incident category/type
    command/outbox/                      — outbox.Publish (transactional outbox)
  handler/
    command/                             — gRPC command-side handlers (used by command-server)
      orgstructure/                      — one file per RPC: command.go + organization_*, clinic_*, department_* method files
      membership/                        — command.go + errors.go + ~25 method files
      incident/classifier/               — command.go + per-RPC files for categories and types
    query/                               — gRPC query-side handlers (used by query-server)
      orgstructure/                      — query.go + ids.go
      membership/                        — query.go + ids.go
      incident/classifier/               — query.go (self-contained parsers for Category/Type + Organization)
      identity/                          — query.go
      stats/                             — query.go
    gateway/                             — HTTP handlers served by gateway-server
      health.go                          — /healthz + /readyz
api/proto/                               — proto contracts (source of truth)
  buf.yaml                               — module config (lint, breaking, deps)
  event/v1/envelope.proto                — Envelope (transport wrapper)
  event/{organization,clinic,department,employee,system_admin}/v1/
  event/incident/{category,type}/v1/     — per-aggregate events
  command/{orgstructure,membership}/v1/  — gRPC command service contracts
  command/incident/classifier/v1/        — gRPC command service contract
api/openapi/medincident.swagger.json     — generated merged OpenAPI v2 covering both sides (committed)
pkg/                                     — buf-generated Go (committed)
  event/**/*.pb.go                       — event messages
  command/**/{*.pb.go,*_grpc.pb.go,*.pb.gw.go}  — gRPC + gateway stubs (command side)
docs/proto/medincident.md                — generated combined Markdown docs for all proto contracts (committed)
buf.gen.yaml                             — go + grpc + gateway + openapi generation
buf.gen.docs.yaml                        — protoc-gen-doc generation
build/{command,query,gateway}-server.Dockerfile  — multi-stage Docker builds (one per binary)
db/migrations/                           — dbmate migrations (never hand-written).
                                           Includes domain.*, outbox.*, projections.*
                                           (projections.* is unused in Plan 1, wired in Plan 2)
test/integration/orgstructure/           — testcontainers-backed integration suite
configs/                                 — command-server.example.yaml + query-server.example.yaml + gateway-server.example.yaml
```

## Ports

All three binaries default to `:8080` inside their process / container:

- `command-server` — gRPC on `:8080`
- `query-server` — gRPC on `:8080`
- `gateway-server` — HTTP on `:8080`

`Dockerfile` `EXPOSE 8080` everywhere. External port mapping is ops'
responsibility (Docker `-p`, k8s `Service`, ingress). No port variance
in the image layer.

For local multi-binary runs on the same host, override the listen
address in the binary-specific YAML (see comments in
`configs/*.example.yaml`). Never hard-code per-binary ports in Go —
the default in every `defaultXxxServerConfig()` is `:8080` without
exception.

## Subject scheme

Outbox subjects follow `medincident.event.<aggregate>.v1.<action>`:
- `medincident.event.organization.v1.{created,details_changed,legal_address_changed}`
- `medincident.event.clinic.v1.{created,details_changed,physical_address_changed}`
- `medincident.event.department.v1.{created,details_changed}`

## Tooling

Binary compilation goes through `Makefile` (see Hard Rule 15).
Everything else goes through `Taskfile.yml`.

- `make build` — compile all three binaries for the host platform
- `make build-<binary>` — e.g. `make build-gateway-server`
- `make build-all` — cross-compile all binaries for every platform in `PLATFORMS`
- `make build-all-<binary>` — single binary, all platforms
- `make clean` — remove `./dist`

Taskfile commands:

- `task gen` — `go tool buf generate` + docs template (pkg/, api/openapi/, docs/proto/)
- `task gen:check` — verify `pkg/`, `api/openapi/`, `docs/proto/` are in sync with proto
- `task proto:fmt`, `task proto:fmt:check` — buf format for .proto files
- `task proto:lint` — buf lint (STANDARD rules, except FIELD_LOWER_SNAKE_CASE)
- `task proto:breaking` — buf breaking vs origin/main (gracefully skipped if main has no api/proto/)
- `task fmt`, `task fmt:check` — gofumpt + goimports via golangci-lint
- `task lint` — golangci-lint
- `task vuln` — govulncheck
- `task test:unit` — unit tests
- `task test:integration` — testcontainers-backed integration suite
- `task test` — unit + integration
- `task migrate` — apply dbmate migrations
- `task migrate:new -- <name>` — create a new migration file pair
- `task docker:command-server`, `task docker:query-server`, `task docker:gateway-server` — Docker image builds
