package di

import (
	"context"
	"time"

	"github.com/samber/do/v2"
	"google.golang.org/grpc"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	"github.com/medincident/medincident-command-service/internal/config"
	membershiphandler "github.com/medincident/medincident-command-service/internal/handler/membership"
	orghandler "github.com/medincident/medincident-command-service/internal/handler/orgstructure"
)

// forceStopGracePeriod gives Stop() a brief window to flush after a
// GracefulStop deadline expiry before Shutdown returns.
const forceStopGracePeriod = 50 * time.Millisecond

// grpcServerWrapper owns the *grpc.Server lifecycle. Private to di —
// consumers invoke *grpc.Server directly via ProvideGRPCServer.
type grpcServerWrapper struct {
	*grpc.Server
}

// Shutdown issues GracefulStop bounded by the ctx deadline; on
// deadline expiry, it falls back to Stop to force-close.
func (g *grpcServerWrapper) Shutdown(ctx context.Context) error {
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
		time.Sleep(forceStopGracePeriod)
		return ctx.Err()
	}
}

func provideGRPCServerWrapper(injector do.Injector) (*grpcServerWrapper, error) {
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

	membershipHandler, err := do.Invoke[*membershiphandler.MembershipHandler](injector)
	if err != nil {
		return nil, err
	}
	membershipv1.RegisterMembershipServiceServer(server, membershipHandler)

	return &grpcServerWrapper{Server: server}, nil
}

// provideGRPCServer resolves the real *grpc.Server for main.go — it
// doesn't know about the wrapper. The wrapper is still registered in
// the container so samber/do invokes its Shutdown on teardown.
func provideGRPCServer(injector do.Injector) (*grpc.Server, error) {
	w, err := do.Invoke[*grpcServerWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.Server, nil
}
