package di

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-command-service/internal/config"
)

// Error codes emitted by ProvidePostgresDB.
const (
	ErrCodePostgresOpenFailed = "postgres_open_failed"
	ErrCodePostgresTuneFailed = "postgres_tune_failed"
)

// PostgresDB wraps *gorm.DB so samber/do can call Shutdown on it.
// The type is exported only so callers in the same package can use it;
// services and handlers invoke *gorm.DB directly (see ProvideGormDB).
type PostgresDB struct {
	*gorm.DB
}

// Shutdown closes the underlying sql.DB pool.
func (p *PostgresDB) Shutdown(_ context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// ProvidePostgresDB opens a gorm connection to Postgres, applies pool
// tunables from config, and returns a Shutdown-aware wrapper.
//
// No Ping, no warm-up, no connectivity probe — the factory only wires.
// Connection failures surface lazily on the first query from a service.
func ProvidePostgresDB(injector do.Injector) (*PostgresDB, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Error),
	})
	if err != nil {
		return nil, oops.In("di.postgres").
			Code(ErrCodePostgresOpenFailed).
			Wrap(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, oops.In("di.postgres").
			Code(ErrCodePostgresTuneFailed).
			Wrap(err)
	}
	sqlDB.SetMaxOpenConns(cfg.Postgres.Pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Postgres.Pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Postgres.Pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Postgres.Pool.ConnMaxIdleTime)

	logger.Info().Msg("postgres pool wired")
	return &PostgresDB{DB: db}, nil
}

// ProvideGormDB resolves the *gorm.DB pointer from the wrapper so
// services can Invoke it directly without knowing about the wrapper.
func ProvideGormDB(injector do.Injector) (*gorm.DB, error) {
	p, err := do.Invoke[*PostgresDB](injector)
	if err != nil {
		return nil, err
	}
	return p.DB, nil
}
