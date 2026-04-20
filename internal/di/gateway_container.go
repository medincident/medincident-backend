package di

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/medincident/medincident-backend/internal/config"
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

// Error codes emitted by this DI init.
const (
	ErrCodeGatewayDialFailed     = "gateway_dial_failed"
	ErrCodeGatewayRegisterFailed = "gateway_register_failed"
)

// gatewayReadHeaderTimeout is the maximum time allowed to read request
// headers. Limits Slowloris-style attacks (gosec G112).
const gatewayReadHeaderTimeout = 5 * time.Second

// NewGatewayContainer wires every provider the gateway-server binary
// needs: logger, two gRPC ClientConns (command + query), grpc-gateway
// ServeMux, composed HTTP handler, and the http.Server.
func NewGatewayContainer(cfg *config.GatewayServerConfig) (do.Injector, error) {
	injector := do.New()

	do.ProvideValue(injector, cfg)

	// Logging
	do.Provide(injector, provideGatewayLoggerWrapper)
	do.Provide(injector, provideGatewayZerolog)

	// Upstream gRPC connections (wrappers only — *grpc.ClientConn
	// cannot be registered directly because samber/do keys by type and
	// there are two conns).
	do.Provide(injector, provideCommandClientConnWrapper)
	do.Provide(injector, provideQueryClientConnWrapper)

	// Gateway mux + HTTP server
	do.Provide(injector, provideGatewayMux)
	do.Provide(injector, provideGatewayHTTPServer)

	return injector, nil
}

// --- logger (gateway-side) ---

func provideGatewayLoggerWrapper(injector do.Injector) (*loggerWrapper, error) {
	cfg, err := do.Invoke[*config.GatewayServerConfig](injector)
	if err != nil {
		return nil, err
	}
	logger, cleanup, err := buildZerolog(&cfg.Zerolog)
	if err != nil {
		return nil, err
	}
	return &loggerWrapper{logger: logger, cleanup: cleanup}, nil
}

func provideGatewayZerolog(injector do.Injector) (*zerolog.Logger, error) {
	w, err := do.Invoke[*loggerWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.logger, nil
}

// --- gRPC client connections ---

// commandClientConnWrapper / queryClientConnWrapper are distinct types
// so samber/do can resolve two *grpc.ClientConn instances separately
// (samber/do keys by Go type). Both implement the Shutdowner protocol
// so the container teardown closes each connection.
type commandClientConnWrapper struct {
	*grpc.ClientConn
}

type queryClientConnWrapper struct {
	*grpc.ClientConn
}

func (w *commandClientConnWrapper) Shutdown(_ context.Context) error { return w.Close() }
func (w *queryClientConnWrapper) Shutdown(_ context.Context) error   { return w.Close() }

func provideCommandClientConnWrapper(injector do.Injector) (*commandClientConnWrapper, error) {
	cfg, err := do.Invoke[*config.GatewayServerConfig](injector)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(
		cfg.Upstreams.Command.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, oops.In("di.gateway").Code(ErrCodeGatewayDialFailed).
			With("upstream", "command").
			With("address", cfg.Upstreams.Command.Address).
			Wrap(err)
	}
	return &commandClientConnWrapper{ClientConn: conn}, nil
}

func provideQueryClientConnWrapper(injector do.Injector) (*queryClientConnWrapper, error) {
	cfg, err := do.Invoke[*config.GatewayServerConfig](injector)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(
		cfg.Upstreams.Query.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, oops.In("di.gateway").Code(ErrCodeGatewayDialFailed).
			With("upstream", "query").
			With("address", cfg.Upstreams.Query.Address).
			Wrap(err)
	}
	return &queryClientConnWrapper{ClientConn: conn}, nil
}

// --- gateway mux ---

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

func provideGatewayMux(injector do.Injector) (*runtime.ServeMux, error) {
	cmdW, err := do.Invoke[*commandClientConnWrapper](injector)
	if err != nil {
		return nil, err
	}
	qryW, err := do.Invoke[*queryClientConnWrapper](injector)
	if err != nil {
		return nil, err
	}

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(gatewayIncomingHeaderMatcher),
	)
	ctx := context.Background()
	errs := []error{
		// Command side.
		cmdorgv1.RegisterOrgStructureCommandServiceHandler(ctx, mux, cmdW.ClientConn),
		cmdmembershipv1.RegisterMembershipCommandServiceHandler(ctx, mux, cmdW.ClientConn),
		cmdclassifierv1.RegisterIncidentClassifierCommandServiceHandler(ctx, mux, cmdW.ClientConn),
		// Query side.
		qorgv1.RegisterOrgStructureQueryServiceHandler(ctx, mux, qryW.ClientConn),
		qmembershipv1.RegisterMembershipQueryServiceHandler(ctx, mux, qryW.ClientConn),
		qclassifierv1.RegisterIncidentClassifierQueryServiceHandler(ctx, mux, qryW.ClientConn),
		qstatsv1.RegisterStatsQueryServiceHandler(ctx, mux, qryW.ClientConn),
		qidentityv1.RegisterIdentityQueryServiceHandler(ctx, mux, qryW.ClientConn),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, oops.In("di.gateway").Code(ErrCodeGatewayRegisterFailed).Wrap(err)
	}
	return mux, nil
}

// --- HTTP server ---
//
// *http.Server already satisfies samber/do's ShutdownerWithContextAndError
// (native Shutdown(context.Context) error), so no intermediate wrapper is
// needed. The container calls the server's own Shutdown during teardown.

func provideGatewayHTTPServer(injector do.Injector) (*http.Server, error) {
	cfg, err := do.Invoke[*config.GatewayServerConfig](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	mux, err := do.Invoke[*runtime.ServeMux](injector)
	if err != nil {
		return nil, err
	}
	cmdW, err := do.Invoke[*commandClientConnWrapper](injector)
	if err != nil {
		return nil, err
	}
	qryW, err := do.Invoke[*queryClientConnWrapper](injector)
	if err != nil {
		return nil, err
	}

	// Route: /healthz + /readyz served by handler/gateway, everything
	// else falls through to the grpc-gateway mux. Wrap with access log
	// (always) then CORS (if configured).
	router := http.NewServeMux()
	router.Handle("/healthz", gwhandler.Liveness())
	router.Handle("/readyz", gwhandler.Readiness(cmdW.ClientConn, qryW.ClientConn))
	router.Handle("/", mux)

	var handler http.Handler = router
	handler = httpmw.AccessLog(logger)(handler)
	if cors := httpmw.CORS(cfg.Server.HTTP.CORS); cors != nil {
		handler = cors(handler)
	}

	return &http.Server{
		Addr:              cfg.Server.HTTP.Address,
		Handler:           handler,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
	}, nil
}
