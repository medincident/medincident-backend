package di

import (
	"context"
	"time"

	"github.com/samber/do/v2"
	"google.golang.org/grpc"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	"github.com/medincident/medincident-command-service/internal/config"
	orghandler "github.com/medincident/medincident-command-service/internal/handler/orgstructure"
)

// forceStopGracePeriod gives Stop() a brief window to flush after a
// GracefulStop deadline expiry before Shutdown returns.
const forceStopGracePeriod = 50 * time.Millisecond

// GRPCServer wraps *grpc.Server for Shutdown semantics. The type is
// exported so cmd/server/main.go can Invoke it and access the
// embedded *grpc.Server to call Serve on the listener.
type GRPCServer struct {
	*grpc.Server
}

// Shutdown issues GracefulStop bounded by the ctx deadline; on
// deadline expiry, it falls back to Stop to force-close.
func (g *GRPCServer) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		g.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		g.Stop()
		// Brief moment for forced-close to flush.
		time.Sleep(forceStopGracePeriod)
		return ctx.Err()
	}
}

// ProvideGRPCServer constructs the gRPC server and registers the
// OrgStructureService handler on it. No interceptors — error mapping
// is deferred to a future middleware spec.
func ProvideGRPCServer(injector do.Injector) (*GRPCServer, error) {
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, err
	}
	handler, err := do.Invoke[*orghandler.OrgStructureHandler](injector)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Server.GRPC.MaxRecvMsgSize),
	)
	orgstructurev1.RegisterOrgStructureServiceServer(server, handler)
	return &GRPCServer{Server: server}, nil
}
