package main

import (
	"time"

	"github.com/medincident/medincident-backend/internal/config"
)

// Config is the YAML-backed runtime configuration for the query-server
// binary. It carries two NATS blocks: nats_zitadel drives the identity
// projector's JetStream consumer, nats_domain drives the domain event
// consumer.
type Config struct {
	Server      serverConfig          `yaml:"server"       validate:"required"`
	Postgres    config.PostgresConfig `yaml:"postgres"     validate:"required"`
	Zerolog     config.ZerologConfig  `yaml:"zerolog"      validate:"required"`
	NATSZitadel config.NATSConfig     `yaml:"nats_zitadel" validate:"required"`
	NATSDomain  config.NATSConfig     `yaml:"nats_domain"  validate:"required"`
	Zitadel     config.ZitadelConfig  `yaml:"zitadel"      validate:"required"`
}

type serverConfig struct {
	GRPC config.GRPCServerConfig `yaml:"grpc" validate:"required"`
}

func defaultConfig() Config {
	return Config{
		Server: serverConfig{
			GRPC: config.GRPCServerConfig{
				Address:        ":8080",
				MaxRecvMsgSize: 4 * 1024 * 1024,
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
	return &cfg, nil
}
