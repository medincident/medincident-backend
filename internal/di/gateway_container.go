package di

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/medincident/medincident-command-service/internal/config"
	"github.com/medincident/medincident-command-service/internal/service/gateway"
)

// Error codes emitted by this DI init.
const (
	ErrCodeGatewayDialFailed = "gateway_dial_failed"
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

func provideGatewayMux(injector do.Injector) (*runtime.ServeMux, error) {
	cmdW, err := do.Invoke[*commandClientConnWrapper](injector)
	if err != nil {
		return nil, err
	}
	qryW, err := do.Invoke[*queryClientConnWrapper](injector)
	if err != nil {
		return nil, err
	}
	return gateway.NewServeMux(context.Background(), cmdW.ClientConn, qryW.ClientConn)
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

	handler := gateway.BuildHandler(
		mux,
		gateway.Liveness(),
		gateway.Readiness(cmdW.ClientConn, qryW.ClientConn),
		gateway.AccessLog(logger),
		gateway.CORSMiddleware(cfg.Server.HTTP.CORS),
	)

	return &http.Server{
		Addr:              cfg.Server.HTTP.Address,
		Handler:           handler,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
	}, nil
}
