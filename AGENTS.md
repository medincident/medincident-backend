# Agent Instructions — medincident-command-service

HTTP/gRPC command side of a CQRS split. Go 1.26 · pgx/v5 · sqlc ·
NATS JetStream · samber/do · samber/oops · zerolog. Full design in
`docs/superpowers/specs/2026-04-10-command-service-base-design.md`.

## Hard rules — NOT NEGOTIABLE

1. **Never push.** All work lands on local commits in a feature branch.
   Never `git push` without explicit user instruction.
2. **Never hand-create migration files.** Always run
   `task migrate:new -- <name>` (which calls `go tool dbmate new`).
   Never use `Write`/`Edit` to create a file under `db/migrations/`.
3. **Never use string literals in `oops.Code(...)` calls, and never
   use generic codes.** Every Code is a specific package-level
   constant (e.g. `ErrCodeAddressTextEmpty`, not `CodeInvalidArgument`)
   declared **in the same file as the code that emits it** — NOT in a
   separate `codes.go`. The rare exception is a code shared by two
   producers in different files of the same package; only then may it
   live in its own file. Aggregates are single-file by convention, so
   this exception almost never applies. Test assertions reference the
   same constants. Grep is the source of truth for where a code is
   emitted and matched.
4. **Never install Go tools globally.** Every Go tool (buf,
   golangci-lint, sqlc, dbmate, govulncheck, protoc-gen-go) is pinned
   in `go.mod`'s `tool` directive and invoked via `go tool <name>`.
   Adding a new tool: `go get -tool <module>@latest` and reference it
   via `go tool` in `Taskfile.yml`.
5. **Never inline invariant limits.** Every max/min/threshold is a
   named package-level constant (`maxNameLen`, `MinLongitude`, …).
   `With()` clauses also reference the constant, never the literal.
6. **Never use `fmt.Errorf("%w", ...)` or stdlib `errors.New` for
   domain errors.** Use `samber/oops`:
   `oops.In("pkg.subsystem").Code(pkg.CodeFoo).Public("…").With("field", x).Wrap(err)`.
   `errors.New` is OK in tests for sentinel comparisons.
7. **Value Objects have ONE constructor, no mutating methods.**
   Aggregates/entities have TWO constructors: `New` (validates, raises
   Created event) and `Hydrate` (trusted, no validation, no events).
8. **Application services take `XCommand` and return `XResult`
   structs.** Positional params are reserved for domain constructors.
   Never write `Service.Create(ctx, name, desc, addr)` — only
   `Service.Create(ctx, CreateCommand{…}) (CreateResult, error)`.
9. **Domain types carry NO struct tags.** Proto-tags and JSON-tags are
   transport concerns. The outbox layer marshals domain events via
   `encoding/json` using exported field names directly; this means
   renaming a Go field on a domain event is a breaking change for
   already-persisted outbox rows.
10. **Events carry NEW state only.** `Renamed.Name`, not
    `Renamed.OldName`+`Renamed.NewName`. Consumers compute diffs from
    their own previous projection.
11. **`aggregate.Root.Raise(event, now)` is the ONLY place
    `UpdatedAt` moves.** Never write `o.UpdatedAt = now` directly in
    a mutating method.
12. **Multi-error validation spans the WHOLE operation, not just one
    constructor.** Within a single command (e.g. `Service.Create`) the
    validator must run every VO and aggregate constructor it can and
    collect all failures with `errors.Join`, so the client receives
    every field violation in one response. Inside a constructor with
    multiple validators, collect with `errors.Join` as well. The only
    time a later constructor is skipped is when it genuinely cannot
    run without a prior VO (pass `nil`/zero where semantically valid
    — e.g. a nil legal address — and continue).
13. **Proto is transport-only.** Outbox stores JSONB of domain event
    structs. The only places proto is imported in the feature tree
    are `internal/orgstructure/<aggregate>/infra/outbox.go` (for
    mappers) and the future `internal/grpcapi/` (for transport).
14. **Never `oops.Wrap(errors.Join(...))` in validation code.** An
    oops-wrapped error implements `Unwrap() error` (singular), not
    `Unwrap() []error` — so the gRPC error interceptor's
    `flattenErrors` recursive unwrap treats it as a leaf and silently
    loses the individual field violations. At the VO / aggregate /
    operation level return `errors.Join(...)` raw, without oops
    wrapping. The interceptor walks the joined tree and inspects each
    leaf's `oops.Code()` to build `BadRequest.FieldViolation` entries
    (codes like `address_text_empty`, `organization_name_empty` each
    identify a distinct violation). The ONE exception is
    `internal/tx/within.go` where rollback-also-fails wraps a joined
    error with `Code(CodeRollbackFailed)` — safe only because tx
    errors never pass through validation flattening. Do not copy this
    pattern elsewhere.

## Directory layout

- `cmd/server/` — server entry point
- `configs/` — YAML config reference + `.env.example`
- `db/migrations/` — dbmate SQL migrations (NEVER hand-crafted)
- `gen/` — buf-generated proto Go bindings (committed)
- `internal/config/` — config DTOs + YAML loader (copied from
  zitadel-actions pattern)
- `internal/correlation/` — request-scoped correlation id via ctx
- `internal/di/` — samber/do providers, including the zerolog builder
  copied verbatim from zitadel-actions
- `internal/errcodes` — **does not exist**; codes live in the same
  file as the code that emits them (a shared-code file is only
  permitted when two producers in different files of the same package
  both need it — rare)
- `internal/orgstructure/` — bounded context for organizational
  structure (Organization today; Clinic, Department later)
- `internal/outbox/` — generic outbox: Store, Registry, Publish, Subject
- `internal/shared/aggregate/` — embed base with
  `CreatedAt`/`UpdatedAt`/events + `Raise`/`PullEvents`
- `internal/shared/geo/` — Point, Address value objects
- `internal/storage/postgres/` — the ONLY place `pgx` is imported
- `internal/tx/` — driver-agnostic transaction abstraction
  (`tx.Tx`, `tx.Beginner`, `tx.Within`)
- `test/integration/` — integration tests behind the `integration`
  build tag (`task test:integration`)

## Tooling

All via `Taskfile.yml` — see README.md for the list. Key commands:

- `task gen` — regenerate proto + sqlc
- `task lint` — golangci-lint
- `task fmt` / `task fmt:check` — format via golangci-lint formatters
- `task test` — unit + integration
- `task test:unit` — fast unit tests only
- `task test:integration` — testcontainers-backed integration suite
- `task vuln` — govulncheck
- `task migrate` / `task migrate:new -- <name>` — dbmate up / new
