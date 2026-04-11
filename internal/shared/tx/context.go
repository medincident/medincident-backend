package tx

import "context"

type ctxKey struct{}

// WithContext stores tx in ctx so repositories can pick it up without
// threading it through every signature.
func WithContext(ctx context.Context, t Tx) context.Context {
	return context.WithValue(ctx, ctxKey{}, t)
}

// FromContext returns the transaction stored in ctx if any.
func FromContext(ctx context.Context) (Tx, bool) {
	t, ok := ctx.Value(ctxKey{}).(Tx)
	return t, ok
}
