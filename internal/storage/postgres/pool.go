// Package postgres is the single place in the codebase that knows pgx.
// Everything outside this package uses the driver-agnostic abstractions
// in internal/tx plus the repository interfaces in each feature's app
// sub-package.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/config"
)

// poolInitCtx is passed to pgxpool during construction. pgx's lazy
// connection model means no connections are dialled at construction time, so
// a background context is correct — the ctx argument is essentially unused.
var poolInitCtx = context.Background()

// Error codes emitted by the pool lifecycle helpers in this file.
const (
	CodePoolParseConfigFailed = "postgres_pool_parse_config_failed"
	CodePoolNewFailed         = "postgres_pool_new_failed"
)

// Pool wraps *pgxpool.Pool with a Shutdown method so it can participate
// in samber/do's lifecycle. No other behaviour is added.
type Pool struct {
	*pgxpool.Pool
}

// NewPool constructs a pgxpool from config and applies pool tuning.
// Liveness is a healthcheck concern — it belongs in the /health handler,
// not here.
func NewPool(cfg config.PostgresConfig) (*Pool, error) {
	pcfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, oops.
			In("storage.postgres").
			Code(CodePoolParseConfigFailed).
			Wrap(err)
	}
	if cfg.Pool.MaxConns > 0 {
		pcfg.MaxConns = cfg.Pool.MaxConns
	}
	if cfg.Pool.MinConns > 0 {
		pcfg.MinConns = cfg.Pool.MinConns
	}
	if cfg.Pool.MaxConnLifetime > 0 {
		pcfg.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	}
	if cfg.Pool.MaxConnIdleTime > 0 {
		pcfg.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(poolInitCtx, pcfg)
	if err != nil {
		return nil, oops.
			In("storage.postgres").
			Code(CodePoolNewFailed).
			Wrap(err)
	}
	return &Pool{Pool: pool}, nil
}

// Shutdown closes the pool. Implements the samber/do lifecycle contract.
func (p *Pool) Shutdown(context.Context) error {
	p.Close()
	return nil
}
