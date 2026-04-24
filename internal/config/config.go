// Package config holds configuration subtypes shared across the three
// binaries (command-server, query-server, gateway-server) and the
// YAML-load + env-expand + validator pipeline that ReadAndValidate
// implements.
//
// Binary-specific top-level configs (e.g. the command-server Config
// struct that wires Server/Postgres/Zerolog/Zitadel together) live
// inside each binary's cmd/<name>/config.go as package main — they
// are not reused across binaries and do not belong here.
//
// File layout:
//   - config.go   — shared subtypes (GRPC, Postgres, Zitadel) + ReadAndValidate
//   - zerolog.go  — ZerologConfig subtree
//   - nats.go     — NATSConfig (currently only query-server uses it)
package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

// Error codes emitted by the shared ReadAndValidate pipeline.
const (
	ErrCodeConfigReadFailed      = "read_failed"
	ErrCodeConfigUnmarshalFailed = "unmarshal_failed"
	ErrCodeConfigValidateFailed  = "validate_failed"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// GRPCServerConfig describes a gRPC listener. Used by both binaries.
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

// ZitadelConfig points at the Zitadel domain and the service-user key
// used for JWT introspection. Shared: both binaries validate tokens the
// same way, so both need the same two fields.
type ZitadelConfig struct {
	Domain  string `yaml:"domain"   validate:"required,url"`
	KeyPath string `yaml:"key_path" validate:"required,file"`
}

// ReadAndValidate is the shared YAML-load + env-expand + validator
// pipeline used by both server-config loaders. It writes into the
// pointed-to struct and returns oops-wrapped errors with the path
// attached as context.
func ReadAndValidate(path string, dst any) error {
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
