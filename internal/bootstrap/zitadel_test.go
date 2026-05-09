package bootstrap_test

import (
	"context"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/medincident/medincident-backend/internal/bootstrap"
	"github.com/medincident/medincident-backend/internal/config"
)

// TestNewZitadelService_EmptyManagementKeyPath asserts that NewZitadelService
// returns a clear configuration error when management_key_path is not set,
// rather than propagating a low-signal "open : no such file" from the SDK.
func TestNewZitadelService_EmptyManagementKeyPath(t *testing.T) {
	cfg := &config.ZitadelConfig{
		Domain:               "https://auth.example.com",
		IntrospectionKeyPath: "/tmp/introspection.json",
		ManagementKeyPath:    "",
	}
	logger := zerolog.Nop()
	_, err := bootstrap.NewZitadelService(context.Background(), cfg, &logger)
	if err == nil {
		t.Fatal("expected error when ManagementKeyPath is empty, got nil")
	}
	if !strings.Contains(err.Error(), "management_key_path") {
		t.Errorf("error should mention management_key_path, got: %v", err)
	}
}
