package di

import (
	"context"

	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
	"github.com/medincident/medincident-command-service/internal/tx"
)

func providePostgres(ctx context.Context, i do.Injector) {
	do.Provide(i, func(inj do.Injector) (*postgres.Pool, error) {
		cfg := do.MustInvoke[*config.Config](inj)
		return postgres.NewPool(ctx, cfg.Postgres)
	})
	do.Provide(i, func(inj do.Injector) (tx.Beginner, error) {
		pool := do.MustInvoke[*postgres.Pool](inj)
		return postgres.NewBeginner(pool), nil
	})
}
