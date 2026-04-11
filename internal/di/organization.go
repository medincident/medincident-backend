package di

import (
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	organizationapp "github.com/medincident/medincident-command-service/internal/orgstructure/organization/app"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/clock"
	"github.com/medincident/medincident-command-service/internal/shared/tx"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
)

// ProvideOrganizationRepository is a samber/do provider for
// organizationapp.Repository.
func ProvideOrganizationRepository(injector do.Injector) (organizationapp.Repository, error) {
	pool, err := do.Invoke[*postgres.Pool](injector)
	if err != nil {
		return nil, err
	}
	return postgres.NewOrganizationRepo(pool), nil
}

// ProvideOrganizationService is a samber/do provider for
// *organizationapp.Service.
func ProvideOrganizationService(injector do.Injector) (*organizationapp.Service, error) {
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	bg, err := do.Invoke[tx.Beginner](injector)
	if err != nil {
		return nil, err
	}
	repo, err := do.Invoke[organizationapp.Repository](injector)
	if err != nil {
		return nil, err
	}
	store, err := do.Invoke[outbox.Store](injector)
	if err != nil {
		return nil, err
	}
	reg, err := do.Invoke[outbox.Registry](injector)
	if err != nil {
		return nil, err
	}
	return organizationapp.NewService(logger, bg, repo, store, reg, clock.System{}), nil
}
