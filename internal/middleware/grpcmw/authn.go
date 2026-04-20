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
	authorizer *authorization.Authorizer[*oauth.IntrospectionContext],
	skip map[string]struct{},
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := skip[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		token := bearerTokenFromMD(ctx)
		if token == "" {
			return nil, oops.In("grpcmw.authn").
				Code(ErrCodeUnauthenticated).
				Public("Missing or malformed authorization header.").
				Errorf("empty bearer token")
		}

		authCtx, err := authorizer.CheckAuthorization(ctx, token)
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

// bearerTokenFromMD extracts the Bearer token from the gRPC
// "authorization" metadata header. Returns "" when absent or when
// the scheme is not "Bearer".
func bearerTokenFromMD(ctx context.Context) string {
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
	return v[len(prefix):]
}
