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

// Error codes emitted by the pool lifecycle helpers in this file.
const (
	CodePoolParseConfigFailed = "postgres_pool_parse_config_failed"
	CodePoolNewFailed         = "postgres_pool_new_failed"
	CodePoolPingFailed        = "postgres_pool_ping_failed"
)

// Pool wraps *pgxpool.Pool with a Shutdown method so it can participate
// in samber/do's lifecycle. No other behaviour is added.
type Pool struct {
	*pgxpool.Pool
}

// NewPool constructs a pgxpool from config, applies pool tuning, and
// pings the database. A failed ping closes the pool before returning.
func NewPool(ctx context.Context, cfg config.PostgresConfig) (*Pool, error) {
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

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, oops.
			In("storage.postgres").
			Code(CodePoolNewFailed).
			Wrap(err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, oops.
			In("storage.postgres").
			Code(CodePoolPingFailed).
			Wrap(err)
	}
	return &Pool{Pool: pool}, nil
}

// Shutdown closes the pool. Implements the samber/do lifecycle contract.
func (p *Pool) Shutdown(context.Context) error {
	p.Close()
	return nil
}
