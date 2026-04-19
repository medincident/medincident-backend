package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, contents string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "gateway.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestReadGatewayServerConfig_Valid(t *testing.T) {
	path := writeTempFile(t, `
server:
  http:
    address: ":8080"
upstreams:
  command:
    address: "command-server:9090"
  query:
    address: "query-server:9091"
zerolog:
  level: info
  timestamp: true
  time_format: "2006-01-02T15:04:05Z07:00"
  outputs:
    - type: console
      target: stdout
`)
	cfg, err := ReadGatewayServerConfig(path)
	if err != nil {
		t.Fatalf("ReadGatewayServerConfig: %v", err)
	}
	if cfg.Server.HTTP.Address != ":8080" {
		t.Errorf("address = %q, want :8080", cfg.Server.HTTP.Address)
	}
	if cfg.Upstreams.Command.Address != "command-server:9090" {
		t.Errorf("command upstream = %q", cfg.Upstreams.Command.Address)
	}
	if cfg.Upstreams.Query.Address != "query-server:9091" {
		t.Errorf("query upstream = %q", cfg.Upstreams.Query.Address)
	}
	if cfg.Server.HTTP.CORS != nil {
		t.Errorf("expected CORS disabled when block omitted, got %+v", cfg.Server.HTTP.CORS)
	}
}

func TestReadGatewayServerConfig_CORSBlock(t *testing.T) {
	path := writeTempFile(t, `
server:
  http:
    address: ":8080"
    cors:
      allowed_origins: ["https://app.example.com"]
      allowed_methods: ["GET", "POST"]
      allowed_headers: ["Authorization"]
      max_age_seconds: 600
upstreams:
  command: {address: "127.0.0.1:9090"}
  query:   {address: "127.0.0.1:9091"}
zerolog:
  level: info
  outputs:
    - type: console
      target: stdout
`)
	cfg, err := ReadGatewayServerConfig(path)
	if err != nil {
		t.Fatalf("ReadGatewayServerConfig: %v", err)
	}
	if cfg.Server.HTTP.CORS == nil {
		t.Fatal("expected CORS non-nil")
	}
	if got := cfg.Server.HTTP.CORS.AllowedOrigins; len(got) != 1 || got[0] != "https://app.example.com" {
		t.Errorf("allowed_origins = %v", got)
	}
	if cfg.Server.HTTP.CORS.MaxAgeSeconds != 600 {
		t.Errorf("max_age_seconds = %d, want 600", cfg.Server.HTTP.CORS.MaxAgeSeconds)
	}
}

func TestReadGatewayServerConfig_MissingUpstream(t *testing.T) {
	path := writeTempFile(t, `
server:
  http:
    address: ":8080"
upstreams:
  command: {address: "127.0.0.1:9090"}
zerolog:
  level: info
  outputs:
    - type: console
      target: stdout
`)
	_, err := ReadGatewayServerConfig(path)
	if err == nil {
		t.Fatal("expected validation error for missing query upstream")
	}
	if !strings.Contains(err.Error(), "Query") && !strings.Contains(err.Error(), "query") {
		t.Errorf("error should mention missing query upstream, got: %v", err)
	}
}

func TestReadGatewayServerConfig_BadAddress(t *testing.T) {
	path := writeTempFile(t, `
server:
  http:
    address: "not-a-host-port"
upstreams:
  command: {address: "127.0.0.1:9090"}
  query:   {address: "127.0.0.1:9091"}
zerolog:
  level: info
  outputs:
    - type: console
      target: stdout
`)
	_, err := ReadGatewayServerConfig(path)
	if err == nil {
		t.Fatal("expected validation error for bad address")
	}
}
