// Command gateway-server is the HTTP gateway binary for the
// medincident platform. It terminates HTTP, forwards the
// Authorization header into gRPC metadata, and proxies requests to
// either command-server or query-server based on the proto package
// of each generated grpc-gateway handler.
package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/di"
)

const shutdownTimeout = 15 * time.Second

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

	bootLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.ReadGatewayServerConfig(configPath)
	if err != nil {
		bootLogger.Fatal().Err(err).Str("config", configPath).Msg("failed to read config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := di.NewGatewayContainer(cfg)
	if err != nil {
		bootLogger.Fatal().Err(err).Msg("failed to build DI container")
	}

	logger := do.MustInvoke[*zerolog.Logger](container)
	server, err := do.Invoke[*http.Server](container)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to resolve http server")
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", cfg.Server.HTTP.Address).Msg("gateway-server http starting")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		close(serveErr)
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("gateway-server stopping (signal)")
	case err, ok := <-serveErr:
		if ok && err != nil {
			logger.Error().Err(err).Msg("http serve error, shutting down")
		} else {
			logger.Info().Msg("gateway-server stopping (http server stopped)")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := container.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("shutdown error")
	}
	logger.Info().Msg("gateway-server stopped")
}
