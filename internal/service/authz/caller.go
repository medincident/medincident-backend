package authz

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
