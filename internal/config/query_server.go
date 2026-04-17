package config

import "time"

// QueryServerConfig is the YAML-backed runtime configuration for the
// query-server binary. The placeholder cmd/query-server/main.go added
// in Plan 1 does not load it yet; the loader and full set of fields
// are wired in Plan 3 alongside the actual readers and the Zitadel
// NATS consumer.
//
// Fields mirror the design spec § 11. Sections that command-server
// has but query-server does not (Zitadel key-based machine auth) are
// intentionally absent. Sections only query-server needs (NATS) are
// only here.
type QueryServerConfig struct {
	Server   QueryServerNetConfig   `yaml:"server"   validate:"required"`
	Postgres PostgresConfig         `yaml:"postgres" validate:"required"`
	Zerolog  ZerologConfig          `yaml:"zerolog"  validate:"required"`
	NATS     QueryServerNATSConfig  `yaml:"nats"     validate:"required"`
	Zitadel  QueryServerZitadelInfo `yaml:"zitadel"  validate:"required"`
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

// QueryServerZitadelInfo is the Zitadel base info needed by readers
// that join projection.users with live Zitadel data, if any.
type QueryServerZitadelInfo struct {
	Domain string `yaml:"domain" validate:"required,url"`
}

// Default values applied when a query-server YAML omits a field.
const (
	defaultQueryServerGRPCAddress    = ":9091"
	defaultQueryServerGatewayAddress = ":8082"
)

func defaultQueryServerConfig() QueryServerConfig {
	return QueryServerConfig{
		Server: QueryServerNetConfig{
			GRPC: GRPCServerConfig{
				Address:        defaultQueryServerGRPCAddress,
				MaxRecvMsgSize: defaultGRPCMaxRecvMsgSize,
			},
			Gateway: QueryServerGatewayConfig{
				Address: defaultQueryServerGatewayAddress,
			},
		},
		Postgres: PostgresConfig{
			Pool: PostgresPoolConfig{
				MaxOpenConns:    defaultPostgresMaxOpenConns,
				MaxIdleConns:    defaultPostgresMaxIdleConns,
				ConnMaxLifetime: defaultPostgresConnMaxLifetime,
				ConnMaxIdleTime: defaultPostgresConnMaxIdleTime,
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
