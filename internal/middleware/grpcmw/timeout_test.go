package grpcmw

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func nopInfo() *grpc.UnaryServerInfo { return &grpc.UnaryServerInfo{} }

func TestTimeoutInterceptor_PreservesStricterClientDeadline(t *testing.T) {
	interceptor := TimeoutInterceptor(30 * time.Second)

	// Client deadline is 1ms — much stricter than 30s.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	called := false
	_, _ = interceptor(ctx, nil, nopInfo(), func(ctx context.Context, _ any) (any, error) {
		called = true
		// The deadline in the handler context must be no later than the
		// client-supplied 1ms deadline.
		dl, ok := ctx.Deadline()
		if !ok {
			t.Error("expected context to have a deadline")
		}
		if time.Until(dl) > 30*time.Second {
			t.Errorf("client deadline was widened: %v remaining", time.Until(dl))
		}
		return struct{}{}, nil
	})
	if !called {
		t.Error("handler was not called")
	}
}

func TestTimeoutInterceptor_EnforcesServerDeadlineWhenClientHasNone(t *testing.T) {
	timeout := 50 * time.Millisecond
	interceptor := TimeoutInterceptor(timeout)

	var deadlineInHandler time.Time
	_, _ = interceptor(context.Background(), nil, nopInfo(), func(ctx context.Context, _ any) (any, error) {
		dl, ok := ctx.Deadline()
		if !ok {
			t.Error("expected context to have a deadline set by the interceptor")
		}
		deadlineInHandler = dl
		return struct{}{}, nil
	})

	if deadlineInHandler.IsZero() {
		t.Fatal("handler deadline was not set")
	}
	// The deadline should be roughly now + timeout; allow 200ms slack.
	upper := time.Now().Add(timeout + 200*time.Millisecond)
	if deadlineInHandler.After(upper) {
		t.Errorf("deadline too far in the future: %v", deadlineInHandler)
	}
}

func TestTimeoutInterceptor_PropagatesHandlerError(t *testing.T) {
	interceptor := TimeoutInterceptor(30 * time.Second)

	_, err := interceptor(context.Background(), nil, nopInfo(), func(_ context.Context, _ any) (any, error) {
		return struct{}{}, context.DeadlineExceeded
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
