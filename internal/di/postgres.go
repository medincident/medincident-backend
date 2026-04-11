package di

import (
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
	"github.com/medincident/medincident-command-service/internal/tx"
)

// ProvidePostgresPool is a samber/do provider for *postgres.Pool.
func ProvidePostgresPool(injector do.Injector) (*postgres.Pool, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	return postgres.NewPool(cfg.Postgres)
}

// ProvideTxBeginner is a samber/do provider for tx.Beginner.
func ProvideTxBeginner(injector do.Injector) (tx.Beginner, error) {
	pool, err := do.Invoke[*postgres.Pool](injector)
	if err != nil {
		return nil, err
	}
	return postgres.NewBeginner(pool), nil
}
