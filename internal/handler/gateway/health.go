// Package gateway holds the HTTP handlers served directly by
// gateway-server (liveness + readiness probes). The grpc-gateway mux,
// the middleware chain, and the http.Server factory all live in
// cmd/gateway-server/main.go — construction belongs there.
package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc/connectivity"
)

// garageProbeTimeout caps a single HeadBucket call so a hung Garage
// instance cannot stall a Kubernetes readiness probe past its own
// deadline. Kept short on purpose: probes run frequently.
const garageProbeTimeout = 2 * time.Second

// StateReader is the narrow slice of *grpc.ClientConn that Readiness
// needs. Defining it lets health_test.go drive state without a real
// gRPC dial.
type StateReader interface {
	GetState() connectivity.State
}

// GarageProbe is the narrow slice of *s3.Client that Readiness needs
// to verify object storage is reachable. Tests substitute a fake
// implementation; production passes the real S3 client.
type GarageProbe interface {
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
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
// of both upstream backends and (when configured) the reachability of
// the Garage S3 bucket. The response body is
//
//	{"command":"READY","query":"CONNECTING","garage":"READY"}
//
// and the status code is 200 if every probed dependency is healthy,
// otherwise 503. gRPC IDLE is treated as healthy (transitions to
// Connecting on demand). Garage is optional — pass a nil probe and an
// empty bucket to omit it from the response entirely.
func Readiness(command, query StateReader, garage GarageProbe, bucket string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		if garage != nil && bucket != "" {
			gs := probeGarage(r.Context(), garage, bucket)
			body["garage"] = gs
			if gs != "READY" {
				code = http.StatusServiceUnavailable
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(body)
	})
}

// probeGarage performs a bounded HeadBucket call. Any error (including
// timeouts) maps to "UNAVAILABLE"; success maps to "READY". The status
// strings echo the gRPC connectivity vocabulary so the readiness body
// stays uniform.
func probeGarage(ctx context.Context, garage GarageProbe, bucket string) string {
	probeCtx, cancel := context.WithTimeout(ctx, garageProbeTimeout)
	defer cancel()
	if _, err := garage.HeadBucket(probeCtx, &s3.HeadBucketInput{Bucket: &bucket}); err != nil {
		return "UNAVAILABLE"
	}
	return "READY"
}

func isHealthy(s connectivity.State) bool {
	return s == connectivity.Ready || s == connectivity.Idle
}
