package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
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
)

// postgresDBWrapper owns the pgxpool lifecycle so samber/do can close
// the pool on injector shutdown. Consumers invoke *gorm.DB directly via
// ProvideGormDB; the pool itself is an implementation detail.
type postgresDBWrapper struct {
	*gorm.DB
	pool *pgxpool.Pool
}

// Shutdown tears down the connection chain in LIFO order: first
// close the *sql.DB created by stdlib.OpenDBFromPool so its
// connectionOpener goroutine exits, then drain and close the
// underlying pgxpool. The pool close runs unconditionally via
// defer — even if sql.DB close fails, the pool must not be leaked.
func (p *postgresDBWrapper) Shutdown(_ context.Context) error {
	defer p.pool.Close()
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

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		return nil, oops.In("di.postgres").
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
		return nil, oops.In("di.postgres").
			Code(ErrCodePostgresOpenFailed).
			Wrap(err)
	}

	logger.Info().Msg("postgres pool wired")
	return &postgresDBWrapper{DB: db, pool: pool}, nil
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
