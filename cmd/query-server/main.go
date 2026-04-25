// Command query-server is the read-side gRPC server for the
// medincident platform. It wires readers, handlers, and a
// JetStream-backed identity consumer explicitly in main — no DI
// framework.
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

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/grpc"

	"github.com/medincident/medincident-backend/internal/bootstrap"
	identityhandler "github.com/medincident/medincident-backend/internal/handler/query/identity"
	classifierhandler "github.com/medincident/medincident-backend/internal/handler/query/incident/classifier"
	membershiphandler "github.com/medincident/medincident-backend/internal/handler/query/membership"
	orghandler "github.com/medincident/medincident-backend/internal/handler/query/orgstructure"
	statshandler "github.com/medincident/medincident-backend/internal/handler/query/stats"
	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	identityread "github.com/medincident/medincident-backend/internal/service/query/identity"
	classifierread "github.com/medincident/medincident-backend/internal/service/query/incident/classifier"
	membershipread "github.com/medincident/medincident-backend/internal/service/query/membership"
	orgread "github.com/medincident/medincident-backend/internal/service/query/orgstructure"
	statsread "github.com/medincident/medincident-backend/internal/service/query/stats"
	identityqueryv1 "github.com/medincident/medincident-backend/pkg/query/identity/v1"
	classifierqueryv1 "github.com/medincident/medincident-backend/pkg/query/incident/classifier/v1"
	membershipqueryv1 "github.com/medincident/medincident-backend/pkg/query/membership/v1"
	orgqueryv1 "github.com/medincident/medincident-backend/pkg/query/orgstructure/v1"
	statsqueryv1 "github.com/medincident/medincident-backend/pkg/query/stats/v1"
)

const (
	shutdownTimeout         = 15 * time.Second
	consumerShutdownTimeout = 10 * time.Second
	natsReconnectWait       = 2 * time.Second
)

const ErrCodeJetStreamInitFailed = "jetstream_init_failed"

// natsDrainTimeout matches the old DI-era shutdown window so a partially
// degraded NATS cannot slow process exit past the k8s pod termination
// grace period.
const natsDrainTimeout = 15 * time.Second

// authnSkip lists the RPC paths that bypass JWT introspection.
var authnSkip = map[string]struct{}{
	"/grpc.health.v1.Health/Check":                                   {},
	"/grpc.health.v1.Health/Watch":                                   {},
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
}

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

	db, dbCleanup, err := bootstrap.OpenPostgres(ctx, &cfg.Postgres, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open postgres")
	}
	defer dbCleanup()

	nc, err := nats.Connect(cfg.NATS.URL,
		nats.Name("medincident-query-server"),
		nats.ReconnectWait(natsReconnectWait),
		nats.MaxReconnects(-1),
		nats.DrainTimeout(natsDrainTimeout),
	)
	if err != nil {
		logger.Fatal().Err(err).Str("url", cfg.NATS.URL).Msg("failed to connect to NATS")
	}
	defer func() { _ = nc.Drain() }()
	logger.Info().Str("url", cfg.NATS.URL).Msg("nats connection established")

	js, err := jetstream.New(nc)
	if err != nil {
		logger.Fatal().
			Err(oops.In("query.bootstrap").Code(ErrCodeJetStreamInitFailed).Wrap(err)).
			Msg("failed to init JetStream")
	}

	authorizer, err := bootstrap.NewZitadelAuthorizer(ctx, &cfg.Zitadel)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to build zitadel authorizer")
	}

	az := authz.New(db)

	orgReader := orgread.NewOrganizationReader(db, logger)
	clinReader := orgread.NewClinicReader(db, az, logger)
	deptReader := orgread.NewDepartmentReader(db, az, logger)
	empReader := membershipread.NewEmployeeReader(db, az, logger)
	roleReader := membershipread.NewRoleReader(db, az, logger)
	classReader := classifierread.NewReader(db, az, logger)
	statsReader := statsread.NewReader(db, az, logger)
	identReader := identityread.NewReader(db, az, logger)

	projector := identityread.NewProjector(db, logger)
	consumer := identityread.NewConsumer(js, &cfg.NATS, projector, logger)

	orgH := orghandler.NewOrgStructureQueryHandler(orgReader, clinReader, deptReader)
	memH := membershiphandler.NewMembershipQueryHandler(empReader, roleReader)
	clsH := classifierhandler.NewIncidentClassifierQueryHandler(classReader)
	statsH := statshandler.NewStatsQueryHandler(statsReader)
	identH := identityhandler.NewIdentityQueryHandler(identReader)

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Server.GRPC.MaxRecvMsgSize),
		grpc.ChainUnaryInterceptor(
			grpcmw.ErrorInterceptor(logger),
			grpcmw.AuthnInterceptor(authorizer, authnSkip),
		),
	)
	orgqueryv1.RegisterOrgStructureQueryServiceServer(grpcServer, orgH)
	membershipqueryv1.RegisterMembershipQueryServiceServer(grpcServer, memH)
	classifierqueryv1.RegisterIncidentClassifierQueryServiceServer(grpcServer, clsH)
	statsqueryv1.RegisterStatsQueryServiceServer(grpcServer, statsH)
	identityqueryv1.RegisterIdentityQueryServiceServer(grpcServer, identH)

	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", cfg.Server.GRPC.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Server.GRPC.Address).Msg("failed to listen")
	}

	// Start the identity consumer BEFORE gRPC so new connections never
	// see a half-booted projection state. If startup is aborted by a
	// shutdown signal, exit quietly and let deferred cleanups run.
	if err := consumer.Start(ctx); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger.Info().Err(err).Msg("consumer start aborted by shutdown signal")
			return
		}
		logger.Error().Err(err).Msg("failed to start identity consumer")
		return
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", cfg.Server.GRPC.Address).Msg("query-server grpc starting")
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
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

	// Shut the consumer first so no new writes land during DB teardown.
	consumerCtx, cancelConsumer := context.WithTimeout(context.Background(), consumerShutdownTimeout)
	if err := consumer.Shutdown(consumerCtx); err != nil {
		logger.Warn().Err(err).Msg("consumer shutdown error")
	}
	cancelConsumer()

	shutdownGRPC(grpcServer, logger)
	logger.Info().Msg("query-server stopped")
}

// shutdownGRPC issues GracefulStop bounded by shutdownTimeout; on
// deadline expiry it falls back to Stop to force-close.
func shutdownGRPC(server *grpc.Server, logger *zerolog.Logger) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		logger.Warn().Dur("timeout", shutdownTimeout).Msg("graceful stop timed out, forcing")
		server.Stop()
	}
}
