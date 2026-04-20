package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-command-service/internal/config"
)

// Error codes emitted by OpenPostgres.
const (
	ErrCodePostgresOpenFailed = "postgres_open_failed"
)

// OpenPostgres builds a *gorm.DB backed by pgxpool from a config.
// It returns the gorm handle and a cleanup function that tears down
// the connection chain (sql.DB → pgxpool) in the correct order.
//
// The caller is responsible for invoking cleanup, typically via
// `defer cleanup()` immediately after the error check.
func OpenPostgres(ctx context.Context, cfg *config.PostgresConfig, logger *zerolog.Logger) (*gorm.DB, func(), error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, nil, oops.In("bootstrap.postgres").
			Code(ErrCodePostgresOpenFailed).
			Wrap(err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Error),
	})
	if err != nil {
		_ = sqlDB.Close()
		pool.Close()
		return nil, nil, oops.In("bootstrap.postgres").
			Code(ErrCodePostgresOpenFailed).
			Wrap(err)
	}

	logger.Info().Msg("postgres pool wired")

	// Teardown order matches DI-era: close sql.DB first so its
	// connectionOpener goroutine exits, then drain the pool. The pool
	// close runs unconditionally via defer — even if sql.DB.Close fails,
	// the pool must not be leaked.
	cleanup := func() {
		defer pool.Close()
		if inner, err := db.DB(); err == nil {
			_ = inner.Close()
		}
	}
	return db, cleanup, nil
}
