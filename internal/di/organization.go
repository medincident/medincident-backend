package di

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	organizationapp "github.com/medincident/medincident-command-service/internal/orgstructure/organization/app"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
	"github.com/medincident/medincident-command-service/internal/tx"
)

// wallClock is the default Clock implementation used in production. Tests
// inject their own clock directly when constructing the Service.
type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }

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
	return organizationapp.NewService(logger, bg, repo, store, reg, wallClock{}), nil
}
