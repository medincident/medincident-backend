package grpcmw

import (
	"context"
	"strings"

	"github.com/samber/oops"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// ErrCodeUnauthenticated is the oops error code for missing or
	// invalid bearer tokens.
	ErrCodeUnauthenticated = "unauthenticated"
)

// AuthnInterceptor returns a gRPC UnaryServerInterceptor that
// validates the Bearer token via Zitadel OAuth2 introspection and
// places the caller's Zitadel user ID into the request context.
//
// Methods whose full name appears in skip are passed through without
// authentication (health checks, reflection, etc.).
func AuthnInterceptor(
	authorizer authorization.AuthorizationChecker[*oauth.IntrospectionContext],
	skip map[string]struct{},
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := skip[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		authHeader := authorizationFromMD(ctx)
		if authHeader == "" {
			return nil, oops.In("grpcmw.authn").
				Code(ErrCodeUnauthenticated).
				Public("Missing or malformed authorization header.").
				Errorf("missing or malformed authorization header")
		}

		authCtx, err := authorizer.CheckAuthorization(ctx, authHeader)
		if err != nil {
			return nil, oops.In("grpcmw.authn").
				Code(ErrCodeUnauthenticated).
				Public("Invalid or expired token.").
				Wrap(err)
		}

		ctx = WithCallerID(ctx, authCtx.UserID())
		return handler(ctx, req)
	}
}

// authorizationFromMD extracts the Authorization header from the gRPC
// incoming metadata. Returns "" when absent or when the scheme is not
// "Bearer". The returned value includes the canonical "Bearer " prefix
// as required by authorization.AuthorizationChecker.CheckAuthorization.
func authorizationFromMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return ""
	}
	v := strings.TrimSpace(values[0])
	const prefix = "Bearer "
	if len(v) <= len(prefix) || !strings.EqualFold(v[:len(prefix)], prefix) {
		return ""
	}
	// Normalize to canonical casing: CheckAuthorization does a
	// case-sensitive strings.CutPrefix("Bearer ", ...) internally.
	return prefix + v[len(prefix):]
}
