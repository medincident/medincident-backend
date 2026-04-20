package di

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/config"
)

// Error codes emitted by the NATS provider.
const (
	ErrCodeNATSConnectFailed   = "nats_connect_failed"
	ErrCodeJetStreamInitFailed = "jetstream_init_failed"
)

// natsConnWrapper owns the *nats.Conn lifecycle so samber/do can Close
// the connection on injector shutdown. Private to di — consumers
// invoke *nats.Conn directly via provideNATSConn.
type natsConnWrapper struct {
	*nats.Conn
}

// Shutdown drains the underlying nats.Conn so any buffered
// publishes/acks are flushed before the connection goes away. Drain
// closes the connection after draining completes.
func (w *natsConnWrapper) Shutdown(_ context.Context) error {
	return w.Drain()
}

func provideNATSConnWrapper(injector do.Injector) (*natsConnWrapper, error) {
	cfg, err := do.Invoke[*config.QueryServerConfig](injector)
	if err != nil {
		return nil, err
	}
	logger, err := do.Invoke[*zerolog.Logger](injector)
	if err != nil {
		return nil, err
	}
	nc, err := nats.Connect(cfg.NATS.URL,
		nats.Name("medincident-query-server"),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, oops.In("di.nats").
			Code(ErrCodeNATSConnectFailed).
			With("url", cfg.NATS.URL).
			Wrap(err)
	}
	logger.Info().Str("url", cfg.NATS.URL).Msg("nats connection established")
	return &natsConnWrapper{Conn: nc}, nil
}

// provideNATSConn returns the underlying *nats.Conn.
func provideNATSConn(injector do.Injector) (*nats.Conn, error) {
	w, err := do.Invoke[*natsConnWrapper](injector)
	if err != nil {
		return nil, err
	}
	return w.Conn, nil
}

// provideJetStream returns the JetStream client bound to the shared
// *nats.Conn.
func provideJetStream(injector do.Injector) (jetstream.JetStream, error) {
	nc, err := do.Invoke[*nats.Conn](injector)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, oops.In("di.nats").
			Code(ErrCodeJetStreamInitFailed).
			Wrap(err)
	}
	return js, nil
}
