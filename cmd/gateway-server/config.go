package main

import (
	"time"

	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/config"
)

// ErrCodeGatewayCORSCredentialsWildcard is emitted when a CORS block
// sets allow_credentials=true while allowed_origins contains the "*"
// wildcard. The CORS spec forbids credentialed requests with a
// wildcard Access-Control-Allow-Origin, so the gateway refuses to
// boot rather than silently producing browser errors at runtime.
const ErrCodeGatewayCORSCredentialsWildcard = "cors_credentials_wildcard"

// Config is the YAML-backed runtime configuration for the gateway-server
// binary. The gateway is a thin HTTP → gRPC translator that fronts
// both command-server and query-server; it owns no state and does
// not share config with either backend.
type Config struct {
	Server    serverConfig         `yaml:"server"    validate:"required"`
	Upstreams upstreamsConfig      `yaml:"upstreams" validate:"required"`
	Zerolog   config.ZerologConfig `yaml:"zerolog"   validate:"required"`
	// Garage is optional. When present, /readyz performs a HeadBucket
	// probe against the configured bucket and returns 503 if it fails.
	Garage *config.GarageConfig `yaml:"garage,omitempty" validate:"omitempty"`
}

type serverConfig struct {
	HTTP httpConfig `yaml:"http" validate:"required"`
}

// httpConfig is the HTTP listener plus optional CORS block. An absent
// cors block (nil pointer) disables CORS; a present block activates
// rs/cors with the configured allowlists.
type httpConfig struct {
	Address string      `yaml:"address"        validate:"required,hostname_port"`
	CORS    *corsConfig `yaml:"cors,omitempty"`
}

type corsConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"   validate:"required,min=1,dive,required"`
	AllowedMethods   []string `yaml:"allowed_methods"   validate:"required,min=1,dive,required"`
	AllowedHeaders   []string `yaml:"allowed_headers"   validate:"required,min=1,dive,required"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAgeSeconds    int      `yaml:"max_age_seconds"   validate:"min=0"`
}

// upstreamsConfig fixes the two upstream gRPC backends the gateway
// proxies into. The map is closed (command + query) — routing is by
// proto package, not by operator-provided key.
type upstreamsConfig struct {
	Command upstreamConfig `yaml:"command" validate:"required"`
	Query   upstreamConfig `yaml:"query"   validate:"required"`
}

// upstreamConfig is the dial target for one upstream. Plaintext only —
// trusted-network deployment is assumed.
type upstreamConfig struct {
	Address string `yaml:"address" validate:"required,hostname_port"`
}

func defaultConfig() Config {
	return Config{
		Server: serverConfig{
			HTTP: httpConfig{
				Address: ":8080",
			},
		},
		Zerolog: config.ZerologConfig{
			Level:      "info",
			Timestamp:  true,
			TimeFormat: time.RFC3339,
			Outputs: []config.ZerologOutputConfig{
				{
					Type:       config.ZerologOutputTypeConsole,
					Target:     config.ZerologConsoleTargetStdout,
					Pretty:     true,
					TimeFormat: "15:04:05",
					PartsOrder: []string{"time", "level", "caller", "message"},
				},
			},
		},
	}
}

func readConfig(path string) (*Config, error) {
	cfg := defaultConfig()
	if err := config.ReadAndValidate(path, &cfg); err != nil {
		return nil, err
	}
	if err := validateCORS(cfg.Server.HTTP.CORS); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validateCORS enforces CORS-spec constraints that struct-tag
// validators cannot express. Specifically: a credentialed CORS policy
// (`allow_credentials: true`) is incompatible with the "*" origin
// wildcard — browsers reject the response when both are combined.
func validateCORS(c *corsConfig) error {
	if c == nil {
		return nil
	}
	if !c.AllowCredentials {
		return nil
	}
	for _, origin := range c.AllowedOrigins {
		if origin == "*" {
			return oops.
				In("gateway").
				Code(ErrCodeGatewayCORSCredentialsWildcard).
				With("allowed_origins", c.AllowedOrigins).
				Errorf("cors: allow_credentials=true is incompatible with allowed_origins=\"*\"")
		}
	}
	return nil
}
