// Package di wires the service's dependency graph via samber/do/v2.
//
// The container is built once in main.go from a loaded *config.Config.
// Providers are idempotent and lifecycle-aware: anything that holds OS
// resources (pgx pool, log file handles) implements the do shutdowner
// protocol so injector.ShutdownWithContext cleans up on exit.
package di

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
)

// NewContainer wires all providers and returns a ready-to-invoke
// injector. The ctx is the startup context — providers that need to
// dial external services use it for their initial connect.
func NewContainer(ctx context.Context, cfg *config.Config) (do.Injector, error) {
	i := do.New()

	// config is a value dependency, not a builder.
	do.ProvideValue(i, cfg)

	// Logger wiring (copied verbatim from zitadel-actions pattern).
	do.Provide(i, ProvideLoggerWrapper)
	do.Provide(i, ProvideZerolog)

	providePostgres(i, ctx)
	provideOutbox(i)
	provideOrganization(i)

	return i, nil
}
