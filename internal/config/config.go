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

// Error codes emitted by Read. File-local per project convention —
// constant names carry the component prefix for grep uniqueness;
// string values drop it because oops.In("config") already conveys
// the component.
const (
	ErrCodeConfigReadFailed      = "read_failed"
	ErrCodeConfigUnmarshalFailed = "unmarshal_failed"
	ErrCodeConfigValidateFailed  = "validate_failed"
)

// validate is the single app-scope validator instance. A package-level
// var rather than a builder function — there is no hidden state to
// reason about (validator.Validate is safe for concurrent use) and
// every Read call reuses the same cached struct metadata.
var validate = validator.New(validator.WithRequiredStructEnabled())

// Config is the root application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"   validate:"required"`
	Postgres PostgresConfig `yaml:"postgres" validate:"required"`
	Nats     NatsConfig     `yaml:"nats"     validate:"required"`
	Zitadel  ZitadelConfig  `yaml:"zitadel"  validate:"required"`
	Outbox   OutboxConfig   `yaml:"outbox"   validate:"required"`
	Zerolog  ZerologConfig  `yaml:"zerolog"  validate:"required"`
}

// ServerConfig holds gRPC server listener settings.
type ServerConfig struct {
	GRPC GRPCServerConfig `yaml:"grpc" validate:"required"`
}

// GRPCServerConfig holds per-server listener settings.
type GRPCServerConfig struct {
	// Address is the TCP listen address, e.g. ":9090".
	Address string `yaml:"address" validate:"required,hostname_port"`
	// MaxRecvMsgSize is the maximum inbound message size in bytes.
	MaxRecvMsgSize int `yaml:"max_recv_msg_size" validate:"required,min=1024,max=104857600"`
}

// PostgresConfig holds the write-store connection settings.
type PostgresConfig struct {
	// DSN is a libpq-style connection string, usually injected via
	// the DATABASE_URL environment variable through ${VAR} expansion.
	DSN  string             `yaml:"dsn"  validate:"required,startswith=postgres://|startswith=postgresql://"`
	Pool PostgresPoolConfig `yaml:"pool" validate:"required"`
}

// PostgresPoolConfig holds pgxpool tuning.
type PostgresPoolConfig struct {
	// MaxConns is the maximum number of connections in the pool.
	MaxConns int32 `yaml:"max_conns" validate:"required,min=1,max=10000"`
	// MinConns is the number of connections held open at idle.
	MinConns int32 `yaml:"min_conns" validate:"min=0,ltefield=MaxConns"`
	// MaxConnLifetime caps how long any single connection is reused
	// before pgx rotates it.
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" validate:"required,min=1s"`
	// MaxConnIdleTime caps how long an idle connection stays in the
	// pool before being closed.
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" validate:"required,min=1s"`
}

// NatsConfig holds NATS JetStream connection settings.
type NatsConfig struct {
	// URL is the NATS server URL.
	URL string `yaml:"url" validate:"required,url"`
	// Stream is the JetStream stream name events are published to.
	Stream string `yaml:"stream" validate:"required"`
	// MaxReconnects caps the number of reconnect attempts on connection
	// loss. -1 means infinite (the default). 0 disables reconnection.
	MaxReconnects int `yaml:"max_reconnects" validate:"min=-1"`
	// ReconnectWait is the base delay between reconnect attempts.
	ReconnectWait time.Duration `yaml:"reconnect_wait" validate:"required,min=100ms"`
}

// ZitadelConfig holds Zitadel gRPC client and JWT validation settings.
type ZitadelConfig struct {
	GRPCAddress string `yaml:"grpc_address" validate:"required,hostname_port"`
	JWKSURL     string `yaml:"jwks_url"     validate:"required,url"`
	Issuer      string `yaml:"issuer"       validate:"required,url"`
	Audience    string `yaml:"audience"     validate:"required"`
}

// OutboxConfig holds outbox poller tuning.
type OutboxConfig struct {
	PollInterval   time.Duration `yaml:"poll_interval"   validate:"required,min=10ms,max=1m"`
	BatchSize      int           `yaml:"batch_size"      validate:"required,min=1,max=10000"`
	PublishTimeout time.Duration `yaml:"publish_timeout" validate:"required,min=100ms,max=5m"`
}

// Default values applied by defaultConfig when a field is absent from
// the loaded YAML. Every constant here is a named value referenced both
// by the config defaults and by tests, so that any tuning change is one
// edit with no magic-number drift.
const (
	defaultGRPCAddress        = ":9090"
	defaultGRPCMaxRecvMsgSize = 4 * 1024 * 1024 // 4 MiB

	defaultPostgresMaxConns        = 20
	defaultPostgresMinConns        = 2
	defaultPostgresMaxConnLifetime = 30 * time.Minute
	defaultPostgresMaxConnIdleTime = 5 * time.Minute

	defaultNatsURL           = "nats://localhost:4222"
	defaultNatsStream        = "medincident"
	defaultNatsMaxReconnects = -1 // infinite
	defaultNatsReconnectWait = 2 * time.Second

	defaultOutboxPollInterval   = 500 * time.Millisecond
	defaultOutboxBatchSize      = 100
	defaultOutboxPublishTimeout = 5 * time.Second
)

// defaultConfig returns a Config pre-populated with the package-level
// default constants above. yaml.Unmarshal then merges the loaded file
// on top, so any field the file sets wins, and anything it omits
// inherits the default. This is the single source of truth for
// "what does the service do if this field is missing from YAML".
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
				MaxConns:        defaultPostgresMaxConns,
				MinConns:        defaultPostgresMinConns,
				MaxConnLifetime: defaultPostgresMaxConnLifetime,
				MaxConnIdleTime: defaultPostgresMaxConnIdleTime,
			},
		},
		Nats: NatsConfig{
			URL:           defaultNatsURL,
			Stream:        defaultNatsStream,
			MaxReconnects: defaultNatsMaxReconnects,
			ReconnectWait: defaultNatsReconnectWait,
		},
		Outbox: OutboxConfig{
			PollInterval:   defaultOutboxPollInterval,
			BatchSize:      defaultOutboxBatchSize,
			PublishTimeout: defaultOutboxPublishTimeout,
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

// formatValidationErrors converts a validator error into a sorted slice of
// "Field.Path: tag" strings. If err is not a validator.ValidationErrors the
// raw error message is returned as a single-element slice.
func formatValidationErrors(err error) []string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return []string{err.Error()}
	}

	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		// StructNamespace gives e.g. "Config.Postgres.DSN"; drop the root
		// "Config." prefix so the path is relative to the config struct.
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
