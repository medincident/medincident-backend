package main

import (
	"time"

	"github.com/medincident/medincident-backend/internal/config"
)

// Config is the YAML-backed runtime configuration for the
// command-server binary.
type Config struct {
	Server   serverConfig          `yaml:"server"   validate:"required"`
	Postgres config.PostgresConfig `yaml:"postgres" validate:"required"`
	Zerolog  config.ZerologConfig  `yaml:"zerolog"  validate:"required"`
	Zitadel  config.ZitadelConfig  `yaml:"zitadel"  validate:"required"`
	// Garage is optional — leaving the block out disables S3 wiring.
	// When present, every field inside is validated.
	Garage *config.GarageConfig `yaml:"garage,omitempty" validate:"omitempty"`
}

// serverConfig is the command-server's listener block. Command-server
// exposes gRPC only.
type serverConfig struct {
	GRPC config.GRPCServerConfig `yaml:"grpc" validate:"required"`
}

func defaultConfig() Config {
	return Config{
		Server: serverConfig{
			GRPC: config.GRPCServerConfig{
				Address:              ":8080",
				MaxRecvMsgSize:       4 * 1024 * 1024,
				MaxConcurrentStreams: 500,
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

// readConfig loads the command-server YAML config from path and expands
// ${VAR} / $VAR references against the process environment.
func readConfig(path string) (*Config, error) {
	cfg := defaultConfig()
	if err := config.ReadAndValidate(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
