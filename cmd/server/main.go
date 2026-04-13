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

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/di"
)

const shutdownTimeout = 10 * time.Second

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

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

	logger := do.MustInvoke[*zerolog.Logger](container)
	grpcWrapper, err := do.Invoke[*di.GRPCServer](container)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to resolve grpc server")
	}
	// grpcWrapper embeds *grpc.Server; container.Shutdown will call
	// grpcWrapper.Shutdown for graceful teardown via the samber/do
	// Shutdowner protocol.
	server := grpcWrapper.Server

	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", cfg.Server.GRPC.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Server.GRPC.Address).Msg("failed to listen")
	}

	go func() {
		logger.Info().Str("addr", cfg.Server.GRPC.Address).Msg("grpc server starting")
		if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Error().Err(err).Msg("grpc serve error")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("command-service stopping")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := container.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("shutdown error")
	}
	logger.Info().Msg("command-service stopped")
}
