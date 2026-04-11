package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/di"
	organizationapp "github.com/medincident/medincident-command-service/internal/orgstructure/organization/app"
)

// shutdownTimeout caps how long the DI container has to release its
// resources during graceful termination.
const shutdownTimeout = 10 * time.Second

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

	// bootLogger handles config/DI errors before the configured logger exists.
	bootLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.Read(configPath)
	if err != nil {
		bootLogger.Fatal().Err(err).Str("config", configPath).Msg("failed to read config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := di.NewContainer(cfg)
	if err != nil {
		bootLogger.Fatal().Err(err).Msg("failed to build DI container")
	}

	// Eager-invoke the top of the service graph to fail fast on wiring errors.
	_ = do.MustInvoke[*organizationapp.Service](container)

	logger := do.MustInvoke[*zerolog.Logger](container)
	logger.Info().Str("config", configPath).Msg("command-service started")

	<-ctx.Done()
	logger.Info().Msg("command-service stopping")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := container.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("shutdown error")
	}
	logger.Info().Msg("command-service stopped")
}
