// Command command-server is the write-side gRPC server for the
// medincident platform. It wires services, handlers, and middleware
// explicitly in main — no DI framework.
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/medincident/medincident-backend/internal/bootstrap"
	announcementhandler "github.com/medincident/medincident-backend/internal/handler/command/announcement"
	incidenthandler "github.com/medincident/medincident-backend/internal/handler/command/incident"
	bufferhandler "github.com/medincident/medincident-backend/internal/handler/command/incident/buffer"
	classifierhandler "github.com/medincident/medincident-backend/internal/handler/command/incident/classifier"
	membershiphandler "github.com/medincident/medincident-backend/internal/handler/command/membership"
	orghandler "github.com/medincident/medincident-backend/internal/handler/command/orgstructure"
	requesthandler "github.com/medincident/medincident-backend/internal/handler/command/request"
	requestclassifierhandler "github.com/medincident/medincident-backend/internal/handler/command/request/classifier"
	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	membershipsvc "github.com/medincident/medincident-backend/internal/service/command/membership"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestclassifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
	requestclassifierv1 "github.com/medincident/medincident-backend/pkg/command/request/classifier/v1"
	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

const (
	shutdownTimeout = 10 * time.Second
	handlerTimeout  = 30 * time.Second
)

// authnSkip lists the RPC paths that bypass JWT introspection: health
// and reflection endpoints need to answer before anyone is authed.
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

	zitadelService, err := bootstrap.NewZitadelService(ctx, &cfg.Zitadel, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to build zitadel service")
	}

	authorizer, err := bootstrap.NewZitadelAuthorizer(ctx, &cfg.Zitadel)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to build zitadel authorizer")
	}

	az := authz.New(db)

	orgSvc := orgsvc.NewOrganizationService(db, az, logger)
	clinSvc := orgsvc.NewClinicService(db, az, logger)
	deptSvc := orgsvc.NewDepartmentService(db, az, logger)
	empSvc := membershipsvc.NewEmployeeService(db, az, zitadelService, logger)
	catSvc := classifiersvc.NewIncidentCategoryService(db, az, logger)
	typSvc := classifiersvc.NewIncidentTypeService(db, az, logger)
	incidentSvc := incidentsvc.NewIncidentService(db, az, logger)
	bufferSvc := buffersvc.NewBufferService(db, az, logger)
	reqTypeSvc := requestclassifiersvc.NewRequestTypeService(db, az, logger)
	reqSvc := requestsvc.NewServiceRequestService(db, az, logger)
	announcementSvc := announcementsvc.NewAnnouncementService(db, az, logger)

	orgStructureHandler := orghandler.NewOrgStructureHandler(orgSvc, clinSvc, deptSvc)
	membershipH := membershiphandler.NewMembershipHandler(empSvc)
	classifierH := classifierhandler.NewIncidentClassifierHandler(catSvc, typSvc)
	incidentH := incidenthandler.NewIncidentHandler(incidentSvc)
	bufferH := bufferhandler.NewBufferHandler(bufferSvc)
	reqClassifierH := requestclassifierhandler.NewRequestClassifierHandler(reqTypeSvc)
	reqH := requesthandler.NewServiceRequestHandler(reqSvc)
	announcementH := announcementhandler.NewAnnouncementHandler(announcementSvc)

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

	orgstructurev1.RegisterOrgStructureCommandServiceServer(grpcServer, orgStructureHandler)
	membershipv1.RegisterMembershipCommandServiceServer(grpcServer, membershipH)
	incidentclassifierv1.RegisterIncidentClassifierCommandServiceServer(grpcServer, classifierH)
	incidentv1.RegisterIncidentCommandServiceServer(grpcServer, incidentH)
	bufferv1.RegisterIncidentBufferCommandServiceServer(grpcServer, bufferH)
	requestclassifierv1.RegisterRequestClassifierCommandServiceServer(grpcServer, reqClassifierH)
	requestv1.RegisterServiceRequestCommandServiceServer(grpcServer, reqH)
	announcementv1.RegisterAnnouncementCommandServiceServer(grpcServer, announcementH)

	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", cfg.Server.GRPC.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Server.GRPC.Address).Msg("failed to listen")
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", cfg.Server.GRPC.Address).Msg("command-server grpc starting")
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveErr <- err
			return
		}
		close(serveErr)
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("command-server stopping (signal)")
	case err, ok := <-serveErr:
		if ok && err != nil {
			logger.Error().Err(err).Msg("grpc serve error, shutting down")
		} else {
			logger.Info().Msg("command-server stopping (grpc server stopped)")
		}
	}

	shutdownGRPC(grpcServer, logger)
	logger.Info().Msg("command-server stopped")
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
