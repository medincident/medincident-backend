package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/medincident/medincident-backend/internal/config"
)

// fakeKeyFile creates a temporary file that satisfies the "file" validator tag.
func fakeKeyFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "key-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

// validate is the same validator instance shape that ReadAndValidate uses.
var validate = validator.New(validator.WithRequiredStructEnabled())

func TestZitadelConfig_YAMLFields(t *testing.T) {
	keyFile := fakeKeyFile(t)

	// Verify both key paths unmarshal from their expected YAML field names.
	raw := `
domain: https://auth.example.com
introspection_key_path: ` + keyFile + `
management_key_path: ` + keyFile

	var cfg config.ZitadelConfig
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.IntrospectionKeyPath != keyFile {
		t.Errorf("IntrospectionKeyPath = %q, want %q", cfg.IntrospectionKeyPath, keyFile)
	}
	if cfg.ManagementKeyPath != keyFile {
		t.Errorf("ManagementKeyPath = %q, want %q", cfg.ManagementKeyPath, keyFile)
	}
}

// TestZitadelConfig_LegacyKeyPath asserts that the old "key_path" field is
// not recognised. Before the split into introspection_key_path /
// management_key_path, both bootstrap functions used a single cfg.KeyPath.
// A config that only provides the legacy field must fail validation so that
// a misconfigured server refuses to start rather than silently using an
// empty path.
func TestZitadelConfig_LegacyKeyPath(t *testing.T) {
	keyFile := fakeKeyFile(t)

	raw := `
domain: https://auth.example.com
key_path: ` + keyFile // old, removed field

	var cfg config.ZitadelConfig
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// IntrospectionKeyPath must be empty — the legacy field is not mapped.
	if cfg.IntrospectionKeyPath != "" {
		t.Errorf("expected IntrospectionKeyPath to be empty when only key_path is provided, got %q", cfg.IntrospectionKeyPath)
	}
	// Validation must fail because IntrospectionKeyPath is required.
	if err := validate.Struct(&cfg); err == nil {
		t.Error("expected validation to fail for config with only legacy key_path, but it passed")
	}
}

func TestZitadelConfig_IntrospectionRequiredManagementOptional(t *testing.T) {
	keyFile := fakeKeyFile(t)

	t.Run("both present", func(t *testing.T) {
		cfg := config.ZitadelConfig{
			Domain:               "https://auth.example.com",
			IntrospectionKeyPath: keyFile,
			ManagementKeyPath:    keyFile,
		}
		if err := validate.Struct(&cfg); err != nil {
			t.Errorf("unexpected validation error: %v", err)
		}
	})

	t.Run("only introspection (query-server profile)", func(t *testing.T) {
		cfg := config.ZitadelConfig{
			Domain:               "https://auth.example.com",
			IntrospectionKeyPath: keyFile,
		}
		if err := validate.Struct(&cfg); err != nil {
			t.Errorf("unexpected validation error: %v", err)
		}
	})

	t.Run("introspection missing", func(t *testing.T) {
		cfg := config.ZitadelConfig{
			Domain:            "https://auth.example.com",
			ManagementKeyPath: keyFile,
		}
		if err := validate.Struct(&cfg); err == nil {
			t.Error("expected validation to fail when IntrospectionKeyPath is absent")
		}
	})

	t.Run("management path points to non-existent file", func(t *testing.T) {
		cfg := config.ZitadelConfig{
			Domain:               "https://auth.example.com",
			IntrospectionKeyPath: keyFile,
			ManagementKeyPath:    filepath.Join(t.TempDir(), "nonexistent.json"),
		}
		if err := validate.Struct(&cfg); err == nil {
			t.Error("expected validation to fail when ManagementKeyPath does not exist")
		}
	})
}
