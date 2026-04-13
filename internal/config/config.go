package config

import (
	"errors"
	"os"
	"sort"
	"strings"
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

// Config is the root application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"   validate:"required"`
	Postgres PostgresConfig `yaml:"postgres" validate:"required"`
	Zerolog  ZerologConfig  `yaml:"zerolog"  validate:"required"`
}

type ServerConfig struct {
	GRPC GRPCServerConfig `yaml:"grpc" validate:"required"`
}

type GRPCServerConfig struct {
	Address        string `yaml:"address"           validate:"required,hostname_port"`
	MaxRecvMsgSize int    `yaml:"max_recv_msg_size" validate:"required,min=1024,max=104857600"`
}

type PostgresConfig struct {
	DSN  string             `yaml:"dsn"  validate:"required,startswith=postgres://|startswith=postgresql://"`
	Pool PostgresPoolConfig `yaml:"pool" validate:"required"`
}

type PostgresPoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns"     validate:"required,min=1,max=10000"`
	MaxIdleConns    int           `yaml:"max_idle_conns"     validate:"min=0,ltefield=MaxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"  validate:"required,min=1s"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" validate:"required,min=1s"`
}

// Default values applied by defaultConfig when a field is absent from
// the loaded YAML.
const (
	defaultGRPCAddress        = ":9090"
	defaultGRPCMaxRecvMsgSize = 4 * 1024 * 1024 // 4 MiB

	defaultPostgresMaxOpenConns    = 20
	defaultPostgresMaxIdleConns    = 2
	defaultPostgresConnMaxLifetime = 30 * time.Minute
	defaultPostgresConnMaxIdleTime = 5 * time.Minute
)

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			GRPC: GRPCServerConfig{
				Address:        defaultGRPCAddress,
				MaxRecvMsgSize: defaultGRPCMaxRecvMsgSize,
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

func formatValidationErrors(err error) []string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return []string{err.Error()}
	}
	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		ns := strings.TrimPrefix(fe.StructNamespace(), "Config.")
		tag := fe.Tag()
		if param := fe.Param(); param != "" {
			tag = tag + "=" + param
		}
		msgs = append(msgs, ns+": "+tag)
	}
	sort.Strings(msgs)
	return msgs
}

// Read loads a YAML config file from path and expands ${VAR} / $VAR
// references in its content using the current process environment.
func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, oops.
			In("config").
			Code(ErrCodeConfigReadFailed).
			With("path", path).
			Wrap(err)
	}
	expanded := os.ExpandEnv(string(data))
	cfg := defaultConfig()
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, oops.
			In("config").
			Code(ErrCodeConfigUnmarshalFailed).
			With("path", path).
			Wrap(err)
	}
	if err := validate.Struct(&cfg); err != nil {
		return nil, oops.
			In("config").
			Code(ErrCodeConfigValidateFailed).
			With("path", path).
			With("violations", formatValidationErrors(err)).
			Wrap(err)
	}
	return &cfg, nil
}
