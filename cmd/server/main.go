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
// resources during graceful termination. 10s is generous for pool
// draining and log flushing without making SIGKILL the normal path
// under real pressure.
const shutdownTimeout = 10 * time.Second

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

	// Bootstrap logger used only until the DI container provides the
	// fully-configured zerolog.Logger. Kept minimal so there is a path
	// to log config/DI errors even when the real logger never came up.
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

	// Eager-invoke the service graph to fail fast on wiring errors
	// (bad DSN, unreachable pool, malformed zerolog config, etc.). We
	// intentionally invoke the top-of-graph types; lower layers come
	// along for the ride.
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
