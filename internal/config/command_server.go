package config

import "time"

// CommandServerConfig is the root application configuration for the
// command-server binary. The query-server binary uses QueryServerConfig.
type CommandServerConfig struct {
	Server   ServerConfig   `yaml:"server"   validate:"required"`
	Postgres PostgresConfig `yaml:"postgres" validate:"required"`
	Zerolog  ZerologConfig  `yaml:"zerolog"  validate:"required"`
	Zitadel  ZitadelConfig  `yaml:"zitadel"  validate:"required"`
}

// ServerConfig is the command-server's listener block. Command-server
// currently exposes gRPC only (no grpc-gateway), so this is a thin
// wrapper over GRPCServerConfig.
type ServerConfig struct {
	GRPC GRPCServerConfig `yaml:"grpc" validate:"required"`
}

func defaultCommandServerConfig() CommandServerConfig {
	return CommandServerConfig{
		Server: ServerConfig{
			GRPC: GRPCServerConfig{
				Address:        ":9090",
				MaxRecvMsgSize: 4 * 1024 * 1024,
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

// ReadCommandServerConfig loads a command-server YAML config file from
// path and expands ${VAR} / $VAR references in its content using the
// current process environment.
func ReadCommandServerConfig(path string) (*CommandServerConfig, error) {
	cfg := defaultCommandServerConfig()
	if err := readAndValidate(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
