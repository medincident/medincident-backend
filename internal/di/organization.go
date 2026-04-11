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

// wallClock is the default Clock implementation used in production.
// Tests inject their own fixedClock directly when constructing the
// Service, so wallClock never leaves this package.
type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }

func provideOrganization(i do.Injector) {
	do.Provide(i, func(inj do.Injector) (organizationapp.Repository, error) {
		pool := do.MustInvoke[*postgres.Pool](inj)
		return postgres.NewOrganizationRepo(pool), nil
	})
	do.Provide(i, func(inj do.Injector) (*organizationapp.Service, error) {
		log := do.MustInvoke[*zerolog.Logger](inj)
		bg := do.MustInvoke[tx.Beginner](inj)
		repo := do.MustInvoke[organizationapp.Repository](inj)
		store := do.MustInvoke[outbox.Store](inj)
		reg := do.MustInvoke[outbox.Registry](inj)
		return organizationapp.NewService(log, bg, repo, store, reg, wallClock{}), nil
	})
}
