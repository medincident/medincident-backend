package config

import "time"

// GatewayServerConfig is the YAML-backed runtime configuration for the
// gateway-server binary. The gateway is a thin HTTP → gRPC translator
// that fronts both command-server and query-server; it owns no state
// and does not share config with either backend.
type GatewayServerConfig struct {
	Server    GatewayServerNetConfig `yaml:"server"    validate:"required"`
	Upstreams GatewayUpstreamsConfig `yaml:"upstreams" validate:"required"`
	Zerolog   ZerologConfig          `yaml:"zerolog"   validate:"required"`
}

// GatewayServerNetConfig holds the HTTP listener for the gateway.
// There is no gRPC listener — gateway-server is pure HTTP.
type GatewayServerNetConfig struct {
	HTTP GatewayHTTPConfig `yaml:"http" validate:"required"`
}

// GatewayHTTPConfig is the HTTP listener plus optional CORS block.
// An absent cors block (nil pointer) disables CORS; a present block
// activates rs/cors with the configured allowlists.
type GatewayHTTPConfig struct {
	Address string             `yaml:"address"        validate:"required,hostname_port"`
	CORS    *GatewayCORSConfig `yaml:"cors,omitempty"`
}

// GatewayCORSConfig is a direct projection of the subset of
// cors.Options the gateway exposes to operators.
type GatewayCORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins" validate:"required,min=1,dive,required"`
	AllowedMethods []string `yaml:"allowed_methods" validate:"required,min=1,dive,required"`
	AllowedHeaders []string `yaml:"allowed_headers" validate:"omitempty,dive,required"`
	MaxAgeSeconds  int      `yaml:"max_age_seconds" validate:"min=0"`
}

// GatewayUpstreamsConfig fixes the two upstream gRPC backends the
// gateway proxies into. The map is closed (command + query) — routing
// is by proto package, not by operator-provided key.
type GatewayUpstreamsConfig struct {
	Command GatewayUpstreamConfig `yaml:"command" validate:"required"`
	Query   GatewayUpstreamConfig `yaml:"query"   validate:"required"`
}

// GatewayUpstreamConfig is the dial target for one upstream. Plaintext
// only — trusted-network deployment is assumed (see design spec §3).
type GatewayUpstreamConfig struct {
	Address string `yaml:"address" validate:"required,hostname_port"`
}

func defaultGatewayServerConfig() GatewayServerConfig {
	return GatewayServerConfig{
		Server: GatewayServerNetConfig{
			HTTP: GatewayHTTPConfig{
				Address: ":8080",
			},
		},
		Zerolog: ZerologConfig{
			Level:      "info",
			Timestamp:  true,
			TimeFormat: time.RFC3339,
			Outputs: []ZerologOutputConfig{
				{
					Type:       ZerologOutputTypeConsole,
					Target:     ZerologConsoleTargetStdout,
					Pretty:     true,
					TimeFormat: "15:04:05",
					PartsOrder: []string{"time", "level", "caller", "message"},
				},
			},
		},
	}
}

// ReadGatewayServerConfig loads a gateway-server YAML config file from
// path with the same env-expand + validate pipeline as
// ReadCommandServerConfig / ReadQueryServerConfig.
func ReadGatewayServerConfig(path string) (*GatewayServerConfig, error) {
	cfg := defaultGatewayServerConfig()
	if err := readAndValidate(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
