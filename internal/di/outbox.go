package di

import (
	"github.com/samber/do/v2"

	organizationinfra "github.com/medincident/medincident-command-service/internal/orgstructure/organization/infra"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
)

func provideOutbox(i do.Injector) {
	do.Provide(i, func(do.Injector) (outbox.Store, error) {
		return postgres.NewOutboxStore(), nil
	})
	do.Provide(i, func(do.Injector) (outbox.Registry, error) {
		reg := outbox.NewRegistry()
		// Each feature-infra package registers its own mappers. When
		// Clinic / Department aggregates land, add a line here.
		organizationinfra.RegisterOutboxMappers(reg)
		return reg, nil
	})
}
