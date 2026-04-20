// Command query-server is the read-side gRPC server for the
// medincident platform. It resolves a *grpc.Server from the DI
// container, listens on the configured address, and runs alongside a
// Zitadel JetStream consumer that maintains the identity projections
// asynchronously.
package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"

	"github.com/medincident/medincident-backend/internal/config"
	"github.com/medincident/medincident-backend/internal/di"
	identityread "github.com/medincident/medincident-backend/internal/service/query/identity"
)

const (
	shutdownTimeout         = 15 * time.Second
	consumerShutdownTimeout = 10 * time.Second
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

	bootLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.ReadQueryServerConfig(configPath)
	if err != nil {
		bootLogger.Fatal().Err(err).Str("config", configPath).Msg("failed to read config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := di.NewQueryContainer(cfg)
	if err != nil {
		bootLogger.Fatal().Err(err).Msg("failed to build DI container")
	}

	logger := do.MustInvoke[*zerolog.Logger](container)
	server, err := do.Invoke[*grpc.Server](container)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to resolve grpc server")
	}
	consumer, err := do.Invoke[*identityread.Consumer](container)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to resolve identity consumer")
	}

	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", cfg.Server.GRPC.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Server.GRPC.Address).Msg("failed to listen")
	}

	// Start the identity consumer before gRPC so new connections never
	// see a half-booted projection state. If the process is shutting
	// down (ctx already cancelled) treat the consumer start failure as
	// a graceful early-exit so the DI container's Shutdown path runs.
	if err := consumer.Start(ctx); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger.Info().Err(err).Msg("consumer start aborted by shutdown signal")
		} else {
			logger.Error().Err(err).Msg("failed to start identity consumer")
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = container.ShutdownWithContext(shutdownCtx)
		return
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", cfg.Server.GRPC.Address).Msg("query-server grpc starting")
		if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveErr <- err
			return
		}
		close(serveErr)
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("query-server stopping (signal)")
	case err, ok := <-serveErr:
		if ok && err != nil {
			logger.Error().Err(err).Msg("grpc serve error, shutting down")
		} else {
			logger.Info().Msg("query-server stopping (grpc server stopped)")
		}
	}

	// Stop consumer first so no new writes land during DB shutdown.
	consumerCtx, cancelConsumer := context.WithTimeout(context.Background(), consumerShutdownTimeout)
	if err := consumer.Shutdown(consumerCtx); err != nil {
		logger.Warn().Err(err).Msg("consumer shutdown error")
	}
	cancelConsumer()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := container.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("shutdown error")
	}
	logger.Info().Msg("query-server stopped")
}
