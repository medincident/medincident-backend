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

const (
	ErrCodePostgresOpenFailed = "postgres_open_failed"
	ErrCodePostgresTuneFailed = "postgres_tune_failed"
)

// postgresDBWrapper owns the *gorm.DB lifecycle so samber/do can Close
// the underlying sql.DB pool on injector shutdown. Private to di —
// consumers invoke *gorm.DB directly via ProvideGormDB.
type postgresDBWrapper struct {
	*gorm.DB
}

func (p *postgresDBWrapper) Shutdown(_ context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func providePostgresDBWrapper(injector do.Injector) (*postgresDBWrapper, error) {
	cfg, err := do.Invoke[*config.CommandServerConfig](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
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
	return &postgresDBWrapper{DB: db}, nil
}

// provideGormDB resolves the real *gorm.DB for services and handlers —
// they don't know about the wrapper.
func provideGormDB(injector do.Injector) (*gorm.DB, error) {
	w, err := do.Invoke[*postgresDBWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.DB, nil
}
