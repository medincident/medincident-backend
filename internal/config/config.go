// Package config defines the YAML-backed runtime configuration for the
// command-server and query-server binaries.
//
// File layout:
//   - config.go        — types and helpers shared by both binaries
//   - command_server.go — CommandServerConfig and its loader
//   - query_server.go   — QueryServerConfig and its loader
//   - zerolog.go        — zerolog config subtree (shared)
package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

// Error codes emitted by the shared readAndValidate pipeline.
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
