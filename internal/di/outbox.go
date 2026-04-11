package di

import (
	"github.com/samber/do/v2"

	organizationinfra "github.com/medincident/medincident-command-service/internal/orgstructure/organization/infra"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
)

// ProvideOutboxStore is a samber/do provider for outbox.Store.
func ProvideOutboxStore(_ do.Injector) (outbox.Store, error) {
	return postgres.NewOutboxStore(), nil
}

// ProvideOutboxRegistry is a samber/do provider for outbox.Registry.
// Each feature-infra package registers its own mappers here.
func ProvideOutboxRegistry(_ do.Injector) (outbox.Registry, error) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	return reg, nil
}
