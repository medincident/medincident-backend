// Command gateway-server is the HTTP gateway binary for the
// medincident platform. It terminates HTTP, forwards the Authorization
// header into gRPC metadata, and proxies requests to either
// command-server or query-server based on the proto package of each
// generated grpc-gateway handler.
package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/medincident/medincident-backend/internal/bootstrap"
	gwhandler "github.com/medincident/medincident-backend/internal/handler/gateway"
	"github.com/medincident/medincident-backend/internal/middleware/httpmw"
	cmdclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
	cmdmembershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
	cmdorgv1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
	qidentityv1 "github.com/medincident/medincident-backend/pkg/query/identity/v1"
	qclassifierv1 "github.com/medincident/medincident-backend/pkg/query/incident/classifier/v1"
	qmembershipv1 "github.com/medincident/medincident-backend/pkg/query/membership/v1"
	qorgv1 "github.com/medincident/medincident-backend/pkg/query/orgstructure/v1"
	qstatsv1 "github.com/medincident/medincident-backend/pkg/query/stats/v1"
)

const (
	shutdownTimeout   = 15 * time.Second
	readHeaderTimeout = 5 * time.Second
)

const ErrCodeGatewayRegisterFailed = "gateway_register_failed"

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

	bootLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := readConfig(configPath)
	if err != nil {
		bootLogger.Fatal().Err(err).Str("config", configPath).Msg("failed to read config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, loggerCleanup, err := bootstrap.BuildZerolog(&cfg.Zerolog)
	if err != nil {
		bootLogger.Fatal().Err(err).Msg("failed to build zerolog")
	}
	defer func() {
		if err := loggerCleanup(); err != nil {
			bootLogger.Error().Err(err).Msg("zerolog cleanup error")
		}
	}()

	commandConn, err := grpc.NewClient(
		cfg.Upstreams.Command.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal().Err(err).Str("upstream", "command").Str("address", cfg.Upstreams.Command.Address).
			Msg("failed to dial command upstream")
	}
	defer func() { _ = commandConn.Close() }()

	queryConn, err := grpc.NewClient(
		cfg.Upstreams.Query.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal().Err(err).Str("upstream", "query").Str("address", cfg.Upstreams.Query.Address).
			Msg("failed to dial query upstream")
	}
	defer func() { _ = queryConn.Close() }()

	gwMux, err := buildGatewayMux(ctx, commandConn, queryConn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register gateway handlers")
	}

	router := http.NewServeMux()
	router.Handle("/healthz", gwhandler.Liveness())
	router.Handle("/readyz", gwhandler.Readiness(commandConn, queryConn))
	router.Handle("/", gwMux)

	var handler http.Handler = router
	handler = httpmw.AccessLog(logger)(handler)
	if cfg.Server.HTTP.CORS != nil {
		handler = httpmw.CORS(&cors.Options{
			AllowedOrigins: cfg.Server.HTTP.CORS.AllowedOrigins,
			AllowedMethods: cfg.Server.HTTP.CORS.AllowedMethods,
			AllowedHeaders: cfg.Server.HTTP.CORS.AllowedHeaders,
			MaxAge:         cfg.Server.HTTP.CORS.MaxAgeSeconds,
		})(handler)
	}

	server := &http.Server{
		Addr:              cfg.Server.HTTP.Address,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
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
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http shutdown error")
	}
	logger.Info().Msg("gateway-server stopped")
}

// gatewayIncomingHeaderMatcher forwards the Authorization header into
// gRPC metadata as the canonical lowercase "authorization" key so the
// downstream Zitadel authn interceptor can read it. Everything else
// delegates to grpc-gateway's default matcher.
func gatewayIncomingHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, "Authorization") {
		return "authorization", true
	}
	return runtime.DefaultHeaderMatcher(key)
}

// buildGatewayMux creates the grpc-gateway mux and registers every
// command- and query-side handler against the correct upstream conn.
func buildGatewayMux(ctx context.Context, commandConn, queryConn *grpc.ClientConn) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(gatewayIncomingHeaderMatcher),
	)
	errs := []error{
		// Command side.
		cmdorgv1.RegisterOrgStructureCommandServiceHandler(ctx, mux, commandConn),
		cmdmembershipv1.RegisterMembershipCommandServiceHandler(ctx, mux, commandConn),
		cmdclassifierv1.RegisterIncidentClassifierCommandServiceHandler(ctx, mux, commandConn),
		// Query side.
		qorgv1.RegisterOrgStructureQueryServiceHandler(ctx, mux, queryConn),
		qmembershipv1.RegisterMembershipQueryServiceHandler(ctx, mux, queryConn),
		qclassifierv1.RegisterIncidentClassifierQueryServiceHandler(ctx, mux, queryConn),
		qstatsv1.RegisterStatsQueryServiceHandler(ctx, mux, queryConn),
		qidentityv1.RegisterIdentityQueryServiceHandler(ctx, mux, queryConn),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, oops.In("gateway").Code(ErrCodeGatewayRegisterFailed).Wrap(err)
	}
	return mux, nil
}
