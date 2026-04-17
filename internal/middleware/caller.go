package middleware

import (
	"context"

	"github.com/samber/oops"
)

type callerIDKey struct{}

// WithCallerID returns a derived context carrying the authenticated
// caller's Zitadel user ID.
func WithCallerID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, callerIDKey{}, id)
}

// CallerID extracts the Zitadel user ID set by the authn interceptor.
// Returns an ErrCodeUnauthenticated oops error when the context has no
// caller ID (i.e. the request bypassed authn, or the handler is being
// called in a context that never went through the interceptor chain).
//
// Handlers that require an authenticated caller must propagate the
// returned error directly — it is already shaped for the gRPC error
// interceptor and maps to codes.Unauthenticated.
func CallerID(ctx context.Context) (string, error) {
	v, _ := ctx.Value(callerIDKey{}).(string)
	if v == "" {
		return "", oops.In("middleware.authn").
			Code(ErrCodeUnauthenticated).
			Public("Missing or malformed authorization header.").
			Errorf("caller id is not set in context")
	}
	return v, nil
}
