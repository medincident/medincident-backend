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
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/medincident/medincident-backend/internal/bootstrap"
	analyticshandler "github.com/medincident/medincident-backend/internal/handler/query/analytics"
	announcementqueryhandler "github.com/medincident/medincident-backend/internal/handler/query/announcement"
	incidentqueryhandler "github.com/medincident/medincident-backend/internal/handler/query/incident"
	bufferqueryhandler "github.com/medincident/medincident-backend/internal/handler/query/incident/buffer"
	classifierhandler "github.com/medincident/medincident-backend/internal/handler/query/incident/classifier"
	membershiphandler "github.com/medincident/medincident-backend/internal/handler/query/membership"
	orghandler "github.com/medincident/medincident-backend/internal/handler/query/orgstructure"
	requestqueryhandler "github.com/medincident/medincident-backend/internal/handler/query/request"
	requestclassifierqueryhandler "github.com/medincident/medincident-backend/internal/handler/query/request/classifier"
	selfhandler "github.com/medincident/medincident-backend/internal/handler/query/self"
	statshandler "github.com/medincident/medincident-backend/internal/handler/query/stats"
	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	analyticsread "github.com/medincident/medincident-backend/internal/service/query/analytics"
	announcementread "github.com/medincident/medincident-backend/internal/service/query/announcement"
	domainread "github.com/medincident/medincident-backend/internal/service/query/domain"
	identityread "github.com/medincident/medincident-backend/internal/service/query/identity"
	incidentread "github.com/medincident/medincident-backend/internal/service/query/incident"
	bufferread "github.com/medincident/medincident-backend/internal/service/query/incident/buffer"
	classifierread "github.com/medincident/medincident-backend/internal/service/query/incident/classifier"
	membershipread "github.com/medincident/medincident-backend/internal/service/query/membership"
	orgread "github.com/medincident/medincident-backend/internal/service/query/orgstructure"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	requestread "github.com/medincident/medincident-backend/internal/service/query/request"
	requestclassifierread "github.com/medincident/medincident-backend/internal/service/query/request/classifier"
	selfread "github.com/medincident/medincident-backend/internal/service/query/self"
	statsread "github.com/medincident/medincident-backend/internal/service/query/stats"
	"github.com/medincident/medincident-backend/internal/util/urlutil"
	analyticsqueryv1 "github.com/medincident/medincident-backend/pkg/query/analytics/v1"
	announcementqueryv1 "github.com/medincident/medincident-backend/pkg/query/announcement/v1"
	classifierqueryv1 "github.com/medincident/medincident-backend/pkg/query/incident/classifier/v1"
	incidentqueryv1 "github.com/medincident/medincident-backend/pkg/query/incident/v1"
	membershipqueryv1 "github.com/medincident/medincident-backend/pkg/query/membership/v1"
	orgqueryv1 "github.com/medincident/medincident-backend/pkg/query/orgstructure/v1"
	requestclassifierqueryv1 "github.com/medincident/medincident-backend/pkg/query/request/classifier/v1"
	requestqueryv1 "github.com/medincident/medincident-backend/pkg/query/request/v1"
	selfqueryv1 "github.com/medincident/medincident-backend/pkg/query/self/v1"
	statsqueryv1 "github.com/medincident/medincident-backend/pkg/query/stats/v1"
)

const (
	shutdownTimeout         = 15 * time.Second
	consumerShutdownTimeout = 10 * time.Second
	natsReconnectWait       = 2 * time.Second
	handlerTimeout          = 30 * time.Second
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

	nc, err := nats.Connect(cfg.NATSZitadel.URL,
		nats.Name("medincident-query-server"),
		nats.ReconnectWait(natsReconnectWait),
		nats.MaxReconnects(-1),
		nats.DrainTimeout(natsDrainTimeout),
	)
	if err != nil {
		logger.Fatal().Err(err).Str("url", urlutil.Redact(cfg.NATSZitadel.URL)).Msg("failed to connect to NATS")
	}
	defer func() { _ = nc.Drain() }()
	logger.Info().Str("url", urlutil.Redact(cfg.NATSZitadel.URL)).Msg("nats connection established")

	js, err := jetstream.New(nc)
	if err != nil {
		logger.Fatal().
			Err(oops.In("query.bootstrap").Code(ErrCodeJetStreamInitFailed).Wrap(err)).
			Msg("failed to init JetStream")
	}

	ncDomain, err := nats.Connect(cfg.NATSDomain.URL,
		nats.Name("medincident-query-server-domain"),
		nats.ReconnectWait(natsReconnectWait),
		nats.MaxReconnects(-1),
		nats.DrainTimeout(natsDrainTimeout),
	)
	if err != nil {
		logger.Fatal().Err(err).Str("url", urlutil.Redact(cfg.NATSDomain.URL)).Msg("failed to connect to NATS (domain)")
	}
	defer func() { _ = ncDomain.Drain() }()

	jsDomain, err := jetstream.New(ncDomain)
	if err != nil {
		logger.Fatal().
			Err(oops.In("query.bootstrap").Code(ErrCodeJetStreamInitFailed).Wrap(err)).
			Msg("failed to init JetStream (domain)")
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
	candidateReader := membershipread.NewCandidateReader(db, az, logger)
	classReader := classifierread.NewReader(db, az, logger)
	statsReader := statsread.NewReader(db, az, logger)
	analyticsReader := analyticsread.NewReader(db, az, logger)
	incidentReader := incidentread.NewReader(db, logger)
	bufferReader := bufferread.NewReader(db, logger, incidentReader)
	reqClassifierReader := requestclassifierread.NewReader(db, az, logger)
	reqReader := requestread.NewReader(db, az, logger)
	announcementReader := announcementread.NewReader(db, logger)
	selfReader := selfread.NewSelfReader(db, logger)

	projector := identityread.NewProjector(db, logger)
	consumer := identityread.NewConsumer(js, &cfg.NATSZitadel, projector, logger)

	proj := qprojector.NewProjectors()
	domainConsumer := domainread.NewConsumer(jsDomain, &cfg.NATSDomain, db, proj, logger)

	orgH := orghandler.NewOrgStructureQueryHandler(orgReader, clinReader, deptReader)
	memH := membershiphandler.NewMembershipQueryHandler(empReader, roleReader, candidateReader)
	clsH := classifierhandler.NewIncidentClassifierQueryHandler(classReader)
	statsH := statshandler.NewStatsQueryHandler(statsReader)
	analyticsH := analyticshandler.NewAnalyticsQueryHandler(analyticsReader)
	incidentQH := incidentqueryhandler.NewIncidentQueryHandler(incidentReader)
	bufferQH := bufferqueryhandler.NewBufferQueryHandler(bufferReader)
	combinedIncidentH := incidentqueryhandler.NewCombinedIncidentQueryHandler(incidentQH, bufferQH)
	reqClassifierQH := requestclassifierqueryhandler.NewRequestClassifierQueryHandler(reqClassifierReader)
	reqQH := requestqueryhandler.NewServiceRequestQueryHandler(reqReader)
	announcementQH := announcementqueryhandler.NewAnnouncementQueryHandler(announcementReader)
	selfH := selfhandler.NewSelfQueryHandler(selfReader)

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Server.GRPC.MaxRecvMsgSize),
		grpc.ChainUnaryInterceptor(
			grpcmw.ErrorInterceptor(logger),
			grpcmw.TimeoutInterceptor(handlerTimeout),
			grpcmw.AuthnInterceptor(authorizer, authnSkip),
		),
	)
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)

	orgqueryv1.RegisterOrgStructureQueryServiceServer(grpcServer, orgH)
	membershipqueryv1.RegisterMembershipQueryServiceServer(grpcServer, memH)
	classifierqueryv1.RegisterIncidentClassifierQueryServiceServer(grpcServer, clsH)
	statsqueryv1.RegisterStatsQueryServiceServer(grpcServer, statsH)
	analyticsqueryv1.RegisterAnalyticsQueryServiceServer(grpcServer, analyticsH)
	incidentqueryv1.RegisterIncidentQueryServiceServer(grpcServer, combinedIncidentH)
	requestclassifierqueryv1.RegisterRequestClassifierQueryServiceServer(grpcServer, reqClassifierQH)
	requestqueryv1.RegisterServiceRequestQueryServiceServer(grpcServer, reqQH)
	announcementqueryv1.RegisterAnnouncementQueryServiceServer(grpcServer, announcementQH)
	selfqueryv1.RegisterSelfQueryServiceServer(grpcServer, selfH)

	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", cfg.Server.GRPC.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Server.GRPC.Address).Msg("failed to listen")
	}

	// Start consumers BEFORE gRPC so new connections never see a
	// half-booted projection state. If startup is aborted by a shutdown
	// signal, exit quietly and let deferred cleanups run.
	if err := consumer.Start(ctx); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger.Info().Err(err).Msg("consumer start aborted by shutdown signal")
			return
		}
		logger.Error().Err(err).Msg("failed to start identity consumer")
		return
	}

	if err := domainConsumer.Start(ctx); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger.Info().Err(err).Msg("domain consumer start aborted by shutdown signal")
			return
		}
		logger.Error().Err(err).Msg("failed to start domain consumer")
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

	// Shut consumers first so no new writes land during DB teardown.
	consumerCtx, cancelConsumer := context.WithTimeout(context.Background(), consumerShutdownTimeout)
	if err := consumer.Shutdown(consumerCtx); err != nil {
		logger.Warn().Err(err).Msg("consumer shutdown error")
	}
	cancelConsumer()

	consumerCtx2, cancelConsumer2 := context.WithTimeout(context.Background(), consumerShutdownTimeout)
	if err := domainConsumer.Shutdown(consumerCtx2); err != nil {
		logger.Warn().Err(err).Msg("domain consumer shutdown error")
	}
	cancelConsumer2()

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
