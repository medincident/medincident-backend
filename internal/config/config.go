package config

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

// Error codes emitted by Read.
const (
	ErrCodeConfigReadFailed      = "read_failed"
	ErrCodeConfigUnmarshalFailed = "unmarshal_failed"
	ErrCodeConfigValidateFailed  = "validate_failed"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// CommandServerConfig is the root application configuration for the
// command-server binary. The query-server binary uses
// QueryServerConfig (defined in query_server.go).
type CommandServerConfig struct {
	Server   ServerConfig   `yaml:"server"   validate:"required"`
	Postgres PostgresConfig `yaml:"postgres" validate:"required"`
	Zerolog  ZerologConfig  `yaml:"zerolog"  validate:"required"`
	Zitadel  ZitadelConfig  `yaml:"zitadel"  validate:"required"`
}

type ServerConfig struct {
	GRPC GRPCServerConfig `yaml:"grpc" validate:"required"`
}

type GRPCServerConfig struct {
	Address        string `yaml:"address"           validate:"required,hostname_port"`
	MaxRecvMsgSize int    `yaml:"max_recv_msg_size" validate:"required,min=1024,max=104857600"`
}

// PostgresConfig is the database connection config. Pool tuning is
// expressed in-band via DSN query parameters honored by pgxpool:
//
//   - pool_max_conns, pool_min_conns
//   - pool_max_conn_lifetime, pool_max_conn_idle_time
//   - pool_health_check_period, pool_max_conn_lifetime_jitter
//
// See configs/*.example.yaml for a worked example.
type PostgresConfig struct {
	DSN string `yaml:"dsn" validate:"required,startswith=postgres://|startswith=postgresql://"`
}

type ZitadelConfig struct {
	Domain  string `yaml:"domain"   validate:"required,url"`
	KeyPath string `yaml:"key_path" validate:"required,file"`
}

// Default values applied by defaultConfig when a field is absent from
// the loaded YAML.
const (
	defaultGRPCAddress        = ":9090"
	defaultGRPCMaxRecvMsgSize = 4 * 1024 * 1024 // 4 MiB
)

func defaultCommandServerConfig() CommandServerConfig {
	return CommandServerConfig{
		Server: ServerConfig{
			GRPC: GRPCServerConfig{
				Address:        defaultGRPCAddress,
				MaxRecvMsgSize: defaultGRPCMaxRecvMsgSize,
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

// readAndValidate is the shared YAML-load + env-expand + validator
// pipeline used by both server-config loaders. It writes into the
// pointed-to struct and returns oops-wrapped errors with the path
// attached as context.
func readAndValidate(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return oops.
			In("config").
			Code(ErrCodeConfigReadFailed).
			With("path", path).
			Wrap(err)
	}
	expanded := os.ExpandEnv(string(data))
	if err := yaml.Unmarshal([]byte(expanded), dst); err != nil {
		return oops.
			In("config").
			Code(ErrCodeConfigUnmarshalFailed).
			With("path", path).
			Wrap(err)
	}
	if err := validate.Struct(dst); err != nil {
		return oops.
			In("config").
			Code(ErrCodeConfigValidateFailed).
			With("path", path).
			With("violations", err.Error()).
			Wrap(err)
	}
	return nil
}
