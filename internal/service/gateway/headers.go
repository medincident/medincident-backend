// Package gateway wires the HTTP surface of gateway-server: the
// runtime.ServeMux wrapper, the middleware chain (CORS, access log),
// the health endpoints, and the header-forwarding policy. Nothing in
// this package knows about samber/do — wiring lives in internal/di.
package gateway

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// IncomingHeaderMatcher is a runtime.HeaderMatcherFunc that forwards
// the Authorization header into gRPC metadata (case-insensitively),
// and otherwise delegates to grpc-gateway's default matcher.
//
// This is the minimum policy needed for JWT passthrough: the
// downstream command-server / query-server authn interceptor reads
// the bearer token off gRPC metadata, so the Authorization header
// must reach it unmangled.
func IncomingHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, "Authorization") {
		return "authorization", true
	}
	return runtime.DefaultHeaderMatcher(key)
}
