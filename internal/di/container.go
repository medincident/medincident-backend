// Package di wires the service's dependency graph via samber/do/v2.
//
// The container is built once from a loaded *config.Config. Providers
// are lifecycle-aware: anything that holds external resources (gorm
// pool, log file handles, gRPC server) is held by a private wrapper
// that implements the do Shutdowner protocol. Consumers depend on the
// real type (e.g. *gorm.DB, *grpc.Server), not the wrapper.
//
// Provider functions are all unexported — NewContainer is the only
// public entry point. Consumers resolve dependencies via do.Invoke on
// the returned injector, keyed by the real type.
package di

import (
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
)

// NewContainer wires all providers and returns a ready-to-invoke injector.
func NewContainer(cfg *config.Config) (do.Injector, error) {
	injector := do.New()

	do.ProvideValue(injector, cfg)

	// Logging
	do.Provide(injector, provideLoggerWrapper)
	do.Provide(injector, provideZerolog)

	// Persistence
	do.Provide(injector, providePostgresDBWrapper)
	do.Provide(injector, provideGormDB)

	// Zitadel
	do.Provide(injector, provideZitadelService)

	// Services
	do.Provide(injector, provideOrganizationService)
	do.Provide(injector, provideClinicService)
	do.Provide(injector, provideDepartmentService)
	do.Provide(injector, provideEmployeeService)
	do.Provide(injector, provideIncidentCategoryService)
	do.Provide(injector, provideIncidentTypeService)

	// Handlers
	do.Provide(injector, provideOrgStructureHandler)
	do.Provide(injector, provideMembershipHandler)
	do.Provide(injector, provideIncidentClassifierHandler)

	// gRPC server
	do.Provide(injector, provideGRPCServerWrapper)
	do.Provide(injector, provideGRPCServer)

	return injector, nil
}
