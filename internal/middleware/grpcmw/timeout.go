package grpcmw

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor enforces a server-side deadline on every unary
// RPC. If the incoming context already carries a stricter deadline, the
// existing one is preserved so the client's own timeout wins.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) <= timeout {
			return handler(ctx, req)
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}
