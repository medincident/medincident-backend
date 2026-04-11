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

// ProvideOrganizationRepository provides organizationapp.Repository.
func ProvideOrganizationRepository(injector do.Injector) (organizationapp.Repository, error) {
	pool, err := do.Invoke[*postgres.Pool](injector)
	if err != nil {
		return nil, err
	}
	return postgres.NewOrganizationRepo(pool), nil
}

// ProvideOrganizationService provides organizationapp.Service. Body
// constructs the unexported *service via NewService and returns it as
// the Service interface.
func ProvideOrganizationService(injector do.Injector) (organizationapp.Service, error) {
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
	pub, err := do.Invoke[outbox.Publisher](injector)
	if err != nil {
		return nil, err
	}
	clk, err := do.Invoke[clock.Clock](injector)
	if err != nil {
		return nil, err
	}
	return organizationapp.NewService(logger, bg, repo, pub, clk), nil
}
