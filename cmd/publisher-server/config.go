package main

import (
	"time"

	"github.com/medincident/medincident-backend/internal/config"
)

// Config is the YAML-backed runtime configuration for the publisher-server
// binary. The publisher reads unpublished rows from outbox.events (command DB)
// and forwards them to NATS JetStream. It holds no gRPC server.
type Config struct {
	Postgres config.PostgresConfig `yaml:"postgres" validate:"required"`
	NATS     publisherNATSConfig   `yaml:"nats"     validate:"required"`
	Zerolog  config.ZerologConfig  `yaml:"zerolog"  validate:"required"`
}

type publisherNATSConfig struct {
	URL string `yaml:"url" validate:"required,url"`
}

func defaultConfig() Config {
	return Config{
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
