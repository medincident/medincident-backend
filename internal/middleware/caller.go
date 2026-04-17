package middleware

import "context"

type callerIDKey struct{}

// WithCallerID returns a derived context carrying the authenticated
// caller's Zitadel user ID.
func WithCallerID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, callerIDKey{}, id)
}

// CallerID extracts the Zitadel user ID set by the authn interceptor.
// Returns "" for unauthenticated (skipped) requests.
func CallerID(ctx context.Context) string {
	v, _ := ctx.Value(callerIDKey{}).(string)
	return v
}
