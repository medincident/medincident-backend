package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc/connectivity"
)

// fakeStateReader lets tests drive GetState() deterministically.
type fakeStateReader struct{ state connectivity.State }

func (f *fakeStateReader) GetState() connectivity.State { return f.state }

// fakeGarage lets tests drive HeadBucket outcomes without an S3 server.
type fakeGarage struct{ err error }

func (f *fakeGarage) HeadBucket(_ context.Context, _ *s3.HeadBucketInput, _ ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &s3.HeadBucketOutput{}, nil
}

func TestLivenessAlways200(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", http.NoBody)
	rec := httptest.NewRecorder()
	Liveness().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestReadiness_BothReady(t *testing.T) {
	h := Readiness(
		&fakeStateReader{state: connectivity.Ready},
		&fakeStateReader{state: connectivity.Ready},
		nil, "",
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", http.NoBody))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if body["command"] != "READY" || body["query"] != "READY" {
		t.Errorf("body = %v", body)
	}
	if _, ok := body["garage"]; ok {
		t.Errorf("garage key should be omitted when probe is nil, got body = %v", body)
	}
}

func TestReadiness_IdleCountsAsReady(t *testing.T) {
	h := Readiness(
		&fakeStateReader{state: connectivity.Idle},
		&fakeStateReader{state: connectivity.Ready},
		nil, "",
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", http.NoBody))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 (IDLE is healthy pre-first-RPC)", rec.Code)
	}
}

func TestReadiness_OneBackendDown(t *testing.T) {
	h := Readiness(
		&fakeStateReader{state: connectivity.Ready},
		&fakeStateReader{state: connectivity.TransientFailure},
		nil, "",
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", http.NoBody))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rec.Code)
	}

	var body map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["query"] != "TRANSIENT_FAILURE" {
		t.Errorf("query state in body = %q, want TRANSIENT_FAILURE", body["query"])
	}
	if body["command"] != "READY" {
		t.Errorf("command state in body = %q, want READY", body["command"])
	}
}

func TestReadiness_GarageReady(t *testing.T) {
	h := Readiness(
		&fakeStateReader{state: connectivity.Ready},
		&fakeStateReader{state: connectivity.Ready},
		&fakeGarage{}, "medincident",
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", http.NoBody))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["garage"] != "READY" {
		t.Errorf("garage = %q, want READY", body["garage"])
	}
}

func TestReadiness_GarageDown(t *testing.T) {
	h := Readiness(
		&fakeStateReader{state: connectivity.Ready},
		&fakeStateReader{state: connectivity.Ready},
		&fakeGarage{err: errors.New("connection refused")}, "medincident",
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", http.NoBody))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rec.Code)
	}
	var body map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["garage"] != "UNAVAILABLE" {
		t.Errorf("garage = %q, want UNAVAILABLE", body["garage"])
	}
}
