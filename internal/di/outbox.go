package di

import (
	"github.com/samber/do/v2"

	organizationinfra "github.com/medincident/medincident-command-service/internal/orgstructure/organization/infra"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/clock"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
)

// ProvideOutboxStore provides outbox.Store.
func ProvideOutboxStore(_ do.Injector) (outbox.Store, error) {
	return postgres.NewOutboxStore(), nil
}

// ProvideOutboxRegistry provides outbox.Registry with all feature-infra
// mappers registered.
func ProvideOutboxRegistry(_ do.Injector) (outbox.Registry, error) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	return reg, nil
}

// ProvideOutboxPublisher provides outbox.Publisher. The body constructs
// the unexported *publisher via NewPublisher; samber/do registers it
// under the interface type.
func ProvideOutboxPublisher(injector do.Injector) (outbox.Publisher, error) {
	store, err := do.Invoke[outbox.Store](injector)
	if err != nil {
		return nil, err
	}
	reg, err := do.Invoke[outbox.Registry](injector)
	if err != nil {
		return nil, err
	}
	clk, err := do.Invoke[clock.Clock](injector)
	if err != nil {
		return nil, err
	}
	return outbox.NewPublisher(store, reg, clk), nil
}
