package gateway

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/connectivity"
)

// StateReader is the narrow slice of *grpc.ClientConn that Readiness
// needs. Defining it lets health_test.go drive state without a real
// gRPC dial.
type StateReader interface {
	GetState() connectivity.State
}

// Liveness returns an http.Handler that always answers 200 OK. It is
// meant for a Kubernetes liveness probe: the process is up, its HTTP
// stack can serve requests. Upstream health is Readiness's concern.
func Liveness() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// Readiness returns an http.Handler that reports the ClientConn state
// of both upstream backends. The response body is
//
//	{"command":"READY","query":"CONNECTING"}
//
// and the status code is 200 if both upstreams are in Ready or Idle
// (IDLE is healthy pre-first-RPC — gRPC transitions to Connecting on
// demand), otherwise 503.
func Readiness(command, query StateReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		cs := command.GetState()
		qs := query.GetState()
		body := map[string]string{
			"command": cs.String(),
			"query":   qs.String(),
		}
		code := http.StatusOK
		if !isHealthy(cs) || !isHealthy(qs) {
			code = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(body)
	})
}

func isHealthy(s connectivity.State) bool {
	return s == connectivity.Ready || s == connectivity.Idle
}
