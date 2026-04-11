// Package di wires the service's dependency graph via samber/do/v2.
//
// The container is built once from a loaded *config.Config. Providers are
// lifecycle-aware: anything that holds external resources (pgx pool, log file
// handles) implements the do shutdowner protocol so
// injector.ShutdownWithContext cleans up on exit.
package di

import (
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
)

// NewContainer wires all providers and returns a ready-to-invoke injector.
func NewContainer(cfg *config.Config) (do.Injector, error) {
	injector := do.New()

	do.ProvideValue(injector, cfg)

	do.Provide(injector, ProvideLoggerWrapper)
	do.Provide(injector, ProvideZerolog)
	do.Provide(injector, ProvidePostgresPool)
	do.Provide(injector, ProvideTxBeginner)
	do.Provide(injector, ProvideOutboxStore)
	do.Provide(injector, ProvideOutboxRegistry)
	do.Provide(injector, ProvideOrganizationRepository)
	do.Provide(injector, ProvideOrganizationService)

	return injector, nil
}
