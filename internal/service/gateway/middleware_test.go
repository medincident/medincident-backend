package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/medincident/medincident-command-service/internal/config"
)

func TestAccessLogCapturesStatusAndPath(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte(`{"ok":false}`))
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/v1/things", strings.NewReader(""))
	rec := httptest.NewRecorder()
	AccessLog(&logger)(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Fatalf("inner handler status leaked: got %d", rec.Code)
	}

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log line not valid JSON: %v\n%s", err, buf.String())
	}
	if entry["method"] != http.MethodPost {
		t.Errorf("method = %v", entry["method"])
	}
	if entry["path"] != "/v1/things" {
		t.Errorf("path = %v", entry["path"])
	}
	if got, ok := entry["status"].(float64); !ok || int(got) != http.StatusTeapot {
		t.Errorf("status = %v", entry["status"])
	}
	if _, ok := entry["duration_ms"]; !ok {
		t.Errorf("missing duration_ms; body: %s", buf.String())
	}
}

func TestAccessLogDefaultsStatusTo200(t *testing.T) {
	// A handler that writes without calling WriteHeader implicitly
	// sends 200 OK. The wrapper must report that.
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "hi")
	})
	rec := httptest.NewRecorder()
	AccessLog(&logger)(inner).ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody))

	var entry map[string]any
	_ = json.Unmarshal(buf.Bytes(), &entry)
	if got, ok := entry["status"].(float64); !ok || int(got) != http.StatusOK {
		t.Errorf("status = %v, want 200", entry["status"])
	}
}

func TestCORSMiddleware_NilWhenConfigAbsent(t *testing.T) {
	mw := CORSMiddleware(nil)
	if mw != nil {
		t.Error("CORSMiddleware(nil) should return nil")
	}
}

func TestCORSMiddleware_HandlesPreflight(t *testing.T) {
	cfg := &config.GatewayCORSConfig{
		AllowedOrigins: []string{"https://app.example.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Authorization"},
		MaxAgeSeconds:  600,
	}
	mw := CORSMiddleware(cfg)
	if mw == nil {
		t.Fatal("expected non-nil middleware")
	}
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/v1/things", http.NoBody)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "authorization") // lowercase per Fetch spec
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, "POST") {
		t.Errorf("Access-Control-Allow-Methods = %q, want to contain POST", got)
	}
}
