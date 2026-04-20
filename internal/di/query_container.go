package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-command-service/internal/config"
	identityhandler "github.com/medincident/medincident-command-service/internal/handler/identity"
	classifierhandler "github.com/medincident/medincident-command-service/internal/handler/incident/classifier"
	membershiphandler "github.com/medincident/medincident-command-service/internal/handler/membership"
	orghandler "github.com/medincident/medincident-command-service/internal/handler/orgstructure"
	statshandler "github.com/medincident/medincident-command-service/internal/handler/stats"
	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	identityread "github.com/medincident/medincident-command-service/internal/service/query/identity"
	classifierread "github.com/medincident/medincident-command-service/internal/service/query/incident/classifier"
	membershipread "github.com/medincident/medincident-command-service/internal/service/query/membership"
	orgread "github.com/medincident/medincident-command-service/internal/service/query/orgstructure"
	statsread "github.com/medincident/medincident-command-service/internal/service/query/stats"
	zitadelsvc "github.com/medincident/medincident-command-service/internal/service/zitadel"
	identityqueryv1 "github.com/medincident/medincident-command-service/pkg/query/identity/v1"
	classifierqueryv1 "github.com/medincident/medincident-command-service/pkg/query/incident/classifier/v1"
	membershipqueryv1 "github.com/medincident/medincident-command-service/pkg/query/membership/v1"
	orgqueryv1 "github.com/medincident/medincident-command-service/pkg/query/orgstructure/v1"
	statsqueryv1 "github.com/medincident/medincident-command-service/pkg/query/stats/v1"
)

// NewQueryContainer wires every provider the query-server binary
// needs: logger, Postgres pool, NATS + JetStream, readers, handlers,
// the identity projector + consumer, and the gRPC server. It is a
// distinct entry point from NewContainer so the command-server binary
// never links against NATS at runtime.
func NewQueryContainer(cfg *config.QueryServerConfig) (do.Injector, error) {
	injector := do.New()

	do.ProvideValue(injector, cfg)

	// Logging
	do.Provide(injector, provideQueryLoggerWrapper)
	do.Provide(injector, provideQueryZerolog)

	// Persistence
	do.Provide(injector, provideQueryPostgresWrapper)
	do.Provide(injector, provideQueryGormDB)

	// NATS + JetStream
	do.Provide(injector, provideNATSConnWrapper)
	do.Provide(injector, provideNATSConn)
	do.Provide(injector, provideJetStream)

	// Zitadel authorizer (JWT introspection for the AuthnInterceptor).
	do.Provide(injector, provideQueryAuthorizer)

	// Readers
	do.Provide(injector, provideOrganizationReader)
	do.Provide(injector, provideClinicReader)
	do.Provide(injector, provideDepartmentReader)
	do.Provide(injector, provideEmployeeReader)
	do.Provide(injector, provideRoleReader)
	do.Provide(injector, provideClassifierReader)
	do.Provide(injector, provideStatsReader)
	do.Provide(injector, provideIdentityReader)

	// Identity projector + consumer
	do.Provide(injector, provideIdentityProjector)
	do.Provide(injector, provideIdentityConsumerWrapper)
	do.Provide(injector, provideIdentityConsumer)

	// Handlers
	do.Provide(injector, provideOrgStructureQueryHandler)
	do.Provide(injector, provideMembershipQueryHandler)
	do.Provide(injector, provideClassifierQueryHandler)
	do.Provide(injector, provideStatsQueryHandler)
	do.Provide(injector, provideIdentityQueryHandler)

	// gRPC server
	do.Provide(injector, provideQueryGRPCServerWrapper)
	do.Provide(injector, provideQueryGRPCServer)

	return injector, nil
}

// --- logger providers (query-side) ---

func provideQueryLoggerWrapper(injector do.Injector) (*loggerWrapper, error) {
	cfg, err := do.Invoke[*config.QueryServerConfig](injector)
	if err != nil {
		return nil, err
	}
	logger, cleanup, err := buildZerolog(&cfg.Zerolog)
	if err != nil {
		return nil, err
	}
	return &loggerWrapper{logger: logger, cleanup: cleanup}, nil
}

func provideQueryZerolog(injector do.Injector) (*zerolog.Logger, error) {
	w, err := do.Invoke[*loggerWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.logger, nil
}

// --- postgres providers (query-side) ---

func provideQueryPostgresWrapper(injector do.Injector) (*postgresDBWrapper, error) {
	cfg, err := do.Invoke[*config.QueryServerConfig](injector)
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

	logger.Info().Msg("query-server postgres pool wired")
	return &postgresDBWrapper{DB: db, pool: pool}, nil
}

func provideQueryGormDB(injector do.Injector) (*gorm.DB, error) {
	w, err := do.Invoke[*postgresDBWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.DB, nil
}

// --- reader providers ---

func provideOrganizationReader(injector do.Injector) (*orgread.OrganizationReader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orgread.NewOrganizationReader(db, logger), nil
}

func provideClinicReader(injector do.Injector) (*orgread.ClinicReader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orgread.NewClinicReader(db, logger), nil
}

func provideDepartmentReader(injector do.Injector) (*orgread.DepartmentReader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return orgread.NewDepartmentReader(db, logger), nil
}

func provideEmployeeReader(injector do.Injector) (*membershipread.EmployeeReader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return membershipread.NewEmployeeReader(db, logger), nil
}

func provideRoleReader(injector do.Injector) (*membershipread.RoleReader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return membershipread.NewRoleReader(db, logger), nil
}

func provideClassifierReader(injector do.Injector) (*classifierread.Reader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return classifierread.NewReader(db, logger), nil
}

func provideStatsReader(injector do.Injector) (*statsread.Reader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return statsread.NewReader(db, logger), nil
}

func provideIdentityReader(injector do.Injector) (*identityread.Reader, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return identityread.NewReader(db, logger), nil
}

// --- identity projector + consumer ---

func provideIdentityProjector(injector do.Injector) (*identityread.Projector, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return identityread.NewProjector(db, logger), nil
}

// identityConsumerWrapper lets samber/do call Shutdown on the Consumer.
type identityConsumerWrapper struct {
	*identityread.Consumer
}

// Shutdown forwards to the inner Consumer shutdown.
func (w *identityConsumerWrapper) Shutdown(ctx context.Context) error {
	return w.Consumer.Shutdown(ctx)
}

func provideIdentityConsumerWrapper(injector do.Injector) (*identityConsumerWrapper, error) {
	js, err := do.Invoke[jetstream.JetStream](injector)
	if err != nil {
		return nil, err
	}
	cfg, err := do.Invoke[*config.QueryServerConfig](injector)
	if err != nil {
		return nil, err
	}
	projector, err := do.Invoke[*identityread.Projector](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	return &identityConsumerWrapper{
		Consumer: identityread.NewConsumer(js, cfg, projector, logger),
	}, nil
}

func provideIdentityConsumer(injector do.Injector) (*identityread.Consumer, error) {
	w, err := do.Invoke[*identityConsumerWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.Consumer, nil
}

// --- handlers ---

func provideOrgStructureQueryHandler(injector do.Injector) (*orghandler.OrgStructureQueryHandler, error) {
	org, err := do.Invoke[*orgread.OrganizationReader](injector)
	if err != nil {
		return nil, err
	}
	clin, err := do.Invoke[*orgread.ClinicReader](injector)
	if err != nil {
		return nil, err
	}
	dept, err := do.Invoke[*orgread.DepartmentReader](injector)
	if err != nil {
		return nil, err
	}
	return orghandler.NewOrgStructureQueryHandler(org, clin, dept), nil
}

func provideMembershipQueryHandler(injector do.Injector) (*membershiphandler.MembershipQueryHandler, error) {
	emp, err := do.Invoke[*membershipread.EmployeeReader](injector)
	if err != nil {
		return nil, err
	}
	role, err := do.Invoke[*membershipread.RoleReader](injector)
	if err != nil {
		return nil, err
	}
	return membershiphandler.NewMembershipQueryHandler(emp, role), nil
}

func provideClassifierQueryHandler(injector do.Injector) (*classifierhandler.IncidentClassifierQueryHandler, error) {
	r, err := do.Invoke[*classifierread.Reader](injector)
	if err != nil {
		return nil, err
	}
	return classifierhandler.NewIncidentClassifierQueryHandler(r), nil
}

func provideStatsQueryHandler(injector do.Injector) (*statshandler.StatsQueryHandler, error) {
	r, err := do.Invoke[*statsread.Reader](injector)
	if err != nil {
		return nil, err
	}
	return statshandler.NewStatsQueryHandler(r), nil
}

func provideIdentityQueryHandler(injector do.Injector) (*identityhandler.IdentityQueryHandler, error) {
	r, err := do.Invoke[*identityread.Reader](injector)
	if err != nil {
		return nil, err
	}
	return identityhandler.NewIdentityQueryHandler(r), nil
}

// --- gRPC server ---

// provideQueryAuthorizer wires JWT introspection against the Zitadel
// domain configured for the query-server. Same flow as the
// command-side provideAuthorizer — both binaries validate tokens the
// same way, against the same Zitadel instance.
func provideQueryAuthorizer(injector do.Injector) (*authorization.Authorizer[*oauth.IntrospectionContext], error) {
	cfg, err := do.Invoke[*config.QueryServerConfig](injector)
	if err != nil {
		return nil, err
	}

	hostname, port, tls, _, err := zitadelsvc.ParseDomain(cfg.Zitadel.Domain)
	if err != nil {
		return nil, err
	}
	opts := zitadelsvc.ZitadelOptsFromParsed(port, tls)

	ctx, cancel := context.WithTimeout(context.Background(), zitadelInitTimeout)
	defer cancel()

	return authorization.New[*oauth.IntrospectionContext](
		ctx,
		zitadel.New(hostname, opts...),
		oauth.DefaultAuthorization(cfg.Zitadel.KeyPath),
	)
}

func provideQueryGRPCServerWrapper(injector do.Injector) (*grpcServerWrapper, error) {
	cfg, err := do.Invoke[*config.QueryServerConfig](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	authorizer, err := do.Invoke[*authorization.Authorizer[*oauth.IntrospectionContext]](injector)
	if err != nil {
		return nil, err
	}
	orgHandler, err := do.Invoke[*orghandler.OrgStructureQueryHandler](injector)
	if err != nil {
		return nil, err
	}
	memHandler, err := do.Invoke[*membershiphandler.MembershipQueryHandler](injector)
	if err != nil {
		return nil, err
	}
	clsHandler, err := do.Invoke[*classifierhandler.IncidentClassifierQueryHandler](injector)
	if err != nil {
		return nil, err
	}
	statsHandler, err := do.Invoke[*statshandler.StatsQueryHandler](injector)
	if err != nil {
		return nil, err
	}
	identHandler, err := do.Invoke[*identityhandler.IdentityQueryHandler](injector)
	if err != nil {
		return nil, err
	}

	// Health + reflection are public infrastructure endpoints — no JWT
	// required. Matches command-server's skip set.
	authnSkip := map[string]struct{}{
		"/grpc.health.v1.Health/Check":                                   {},
		"/grpc.health.v1.Health/Watch":                                   {},
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
	}

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Server.GRPC.MaxRecvMsgSize),
		grpc.ChainUnaryInterceptor(
			grpcmw.ErrorInterceptor(logger),
			grpcmw.AuthnInterceptor(authorizer, authnSkip),
		),
	)
	orgqueryv1.RegisterOrgStructureQueryServiceServer(server, orgHandler)
	membershipqueryv1.RegisterMembershipQueryServiceServer(server, memHandler)
	classifierqueryv1.RegisterIncidentClassifierQueryServiceServer(server, clsHandler)
	statsqueryv1.RegisterStatsQueryServiceServer(server, statsHandler)
	identityqueryv1.RegisterIdentityQueryServiceServer(server, identHandler)
	return &grpcServerWrapper{Server: server}, nil
}

func provideQueryGRPCServer(injector do.Injector) (*grpc.Server, error) {
	w, err := do.Invoke[*grpcServerWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.Server, nil
}
