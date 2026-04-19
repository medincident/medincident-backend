package config

import "time"

// QueryServerConfig is the YAML-backed runtime configuration for the
// query-server binary. Fields mirror the design spec § 11. Sections
// only query-server needs (NATS, grpc-gateway) live here.
type QueryServerConfig struct {
	Server   QueryServerNetConfig  `yaml:"server"   validate:"required"`
	Postgres PostgresConfig        `yaml:"postgres" validate:"required"`
	Zerolog  ZerologConfig         `yaml:"zerolog"  validate:"required"`
	NATS     QueryServerNATSConfig `yaml:"nats"     validate:"required"`
	Zitadel  ZitadelConfig         `yaml:"zitadel"  validate:"required"`
}

// QueryServerNetConfig holds the gRPC and grpc-gateway listen addresses
// for the query-server binary.
type QueryServerNetConfig struct {
	GRPC    GRPCServerConfig         `yaml:"grpc"    validate:"required"`
	Gateway QueryServerGatewayConfig `yaml:"gateway" validate:"required"`
}

// QueryServerGatewayConfig is the HTTP listener for grpc-gateway. Kept
// query-only for now because command-server does not currently expose
// a gateway endpoint via configuration.
type QueryServerGatewayConfig struct {
	Address string `yaml:"address" validate:"required,hostname_port"`
}

// QueryServerNATSConfig is the JetStream connection block. Stream and
// subjects are operator-provisioned; query-server only attaches a
// durable consumer to an already-existing stream.
type QueryServerNATSConfig struct {
	URL         string   `yaml:"url"          validate:"required"`
	Stream      string   `yaml:"stream"       validate:"required"`
	Subjects    []string `yaml:"subjects"     validate:"required,min=1,dive,required"`
	DurableName string   `yaml:"durable_name" validate:"required"`
}

func defaultQueryServerConfig() QueryServerConfig {
	return QueryServerConfig{
		Server: QueryServerNetConfig{
			GRPC: GRPCServerConfig{
				Address:        ":9091",
				MaxRecvMsgSize: 4 * 1024 * 1024,
			},
			Gateway: QueryServerGatewayConfig{
				Address: ":8082",
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

// ReadQueryServerConfig loads a query-server YAML config file from
// path with the same env-expand + validate pipeline as
// ReadCommandServerConfig.
func ReadQueryServerConfig(path string) (*QueryServerConfig, error) {
	cfg := defaultQueryServerConfig()
	if err := readAndValidate(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
