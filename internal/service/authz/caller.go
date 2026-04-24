package authz

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
)

// Caller is the authenticated principal that issues a command. Every
// command struct in internal/service/command/... embeds this value so
// the service layer has the caller's identity without reading ctx.
//
// ZitadelUserID is the only trusted fact about the caller. It is set
// by the authn interceptor (internal/middleware/grpcmw/authn.go) after
// validating the incoming JWT, so by the time a service method sees
// it the value is guaranteed to be a non-empty identity string. No
// validate tags here — validation happened at the transport boundary.
type Caller struct {
	ZitadelUserID string
}

// CallerFromCtx extracts the authenticated caller set by the authn
// interceptor and wraps it in a Caller. The returned error is the one
// from grpcmw.CallerID — already shaped for the gRPC error
// interceptor and mapping to codes.Unauthenticated — so handlers
// propagate it directly. This is the standard adapter for query
// handlers, which receive Caller from context rather than from a
// command DTO.
func CallerFromCtx(ctx context.Context) (Caller, error) {
	id, err := grpcmw.CallerID(ctx)
	if err != nil {
		return Caller{}, err
	}
	return Caller{ZitadelUserID: id}, nil
}
