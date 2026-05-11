// Command publisher-server reads unpublished rows from outbox.events and
// publishes them to NATS JetStream in strict seq order. It acquires a
// Postgres advisory lock so that only one replica is active at a time;
// a second replica blocks until the first exits and the lock is released.
//
// Wiring: no DI framework. Dependencies are wired inline via direct
// constructor calls; teardown via defer. Bootstrap helpers must not
// call Ping or warm-up probes — see AGENTS.md rule 13.
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/bootstrap"
	"github.com/medincident/medincident-backend/internal/util/urlutil"
)

const (
	// outboxPublisherLockID is the Postgres advisory lock id that ensures
	// only one publisher replica is active at a time.
	outboxPublisherLockID = int64(7246490831)

	// fetchBatchSize is the maximum number of unpublished outbox rows
	// fetched per publish cycle.
	fetchBatchSize = 100

	// pollInterval is the fallback polling period used when no LISTEN
	// notification arrives (covers reconnect windows and cold starts).
	pollInterval = 5 * time.Second

	// natsReconnectWait is the interval between NATS reconnect attempts.
	natsReconnectWait = 2 * time.Second

	// natsDrainTimeout bounds how long Drain() is allowed to wait on
	// shutdown.
	natsDrainTimeout = 15 * time.Second

	// ErrCodeJetStreamInitFailed is emitted when the JetStream context
	// cannot be initialised from the NATS connection.
	ErrCodeJetStreamInitFailed = "jetstream_init_failed"

	// ErrCodeAdvisoryLockFailed is emitted when acquiring the Postgres
	// advisory lock fails.
	ErrCodeAdvisoryLockFailed = "advisory_lock_failed"

	// ErrCodeListenFailed is emitted when LISTEN on outbox_events fails.
	ErrCodeListenFailed = "listen_failed"

	// ErrCodeOutboxFetchFailed is emitted when fetching unpublished rows
	// from outbox.events fails.
	ErrCodeOutboxFetchFailed = "outbox_fetch_failed"

	// ErrCodeOutboxMarkFailed is emitted when marking a row as published
	// fails after successful NATS publish.
	ErrCodeOutboxMarkFailed = "outbox_mark_failed"

	// ErrCodeOutboxPublishFailed is emitted when publishing a row to NATS
	// JetStream fails.
	ErrCodeOutboxPublishFailed = "outbox_publish_failed"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to the YAML configuration file")
	flag.Parse()

	bootLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := readConfig(configPath)
	if err != nil {
		bootLogger.Fatal().Err(err).Str("config", configPath).Msg("failed to read config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, loggerCleanup, err := bootstrap.BuildZerolog(&cfg.Zerolog)
	if err != nil {
		bootLogger.Fatal().Err(err).Msg("failed to build zerolog")
	}
	defer func() {
		if err := loggerCleanup(); err != nil {
			bootLogger.Error().Err(err).Msg("zerolog cleanup error")
		}
	}()

	// Open pgxpool for the lock connection and the fetch/mark queries.
	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open postgres pool")
	}
	defer pool.Close()

	// Connect to NATS.
	nc, err := nats.Connect(cfg.NATS.URL,
		nats.Name("medincident-publisher-server"),
		nats.ReconnectWait(natsReconnectWait),
		nats.MaxReconnects(-1),
		nats.DrainTimeout(natsDrainTimeout),
	)
	if err != nil {
		logger.Fatal().Err(err).Str("url", urlutil.Redact(cfg.NATS.URL)).Msg("failed to connect to NATS")
	}
	defer func() { _ = nc.Drain() }()
	logger.Info().Str("url", urlutil.Redact(cfg.NATS.URL)).Msg("nats connection established")

	js, err := jetstream.New(nc)
	if err != nil {
		logger.Fatal().
			Err(oops.In("publisher.bootstrap").Code(ErrCodeJetStreamInitFailed).Wrap(err)).
			Msg("failed to init JetStream")
	}

	p := &publisher{pool: pool, js: js, logger: logger}
	p.run(ctx)
	logger.Info().Msg("publisher-server stopped")
}

type publisher struct {
	pool   *pgxpool.Pool
	js     jetstream.JetStream
	logger *zerolog.Logger
}

// run acquires the advisory lock (blocking), then enters the publish loop.
// Returns when ctx is cancelled.
func (p *publisher) run(ctx context.Context) {
	// Acquire advisory lock on a dedicated connection that lives for the
	// lifetime of this process. The lock is released automatically when
	// the connection is closed (process death / normal shutdown).
	lockConn, err := p.pool.Acquire(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return // shutdown before we got started
		}
		p.logger.Fatal().
			Err(oops.In("publisher").Code(ErrCodeAdvisoryLockFailed).Wrap(err)).
			Msg("failed to acquire pool connection for advisory lock")
		return
	}
	defer lockConn.Release()

	p.logger.Info().Msg("waiting for advisory lock (blocking)")
	if _, err := lockConn.Exec(ctx, "SELECT pg_advisory_lock($1)", outboxPublisherLockID); err != nil {
		if ctx.Err() != nil {
			return
		}
		p.logger.Fatal().
			Err(oops.In("publisher").Code(ErrCodeAdvisoryLockFailed).Wrap(err)).
			Msg("failed to acquire advisory lock")
		return
	}
	p.logger.Info().Msg("advisory lock acquired — this replica is the active publisher")

	// Acquire a second connection for LISTEN/NOTIFY.
	listenConn, err := p.pool.Acquire(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		p.logger.Fatal().Err(err).Msg("failed to acquire listen connection")
		return
	}
	defer listenConn.Release()

	if _, err := listenConn.Exec(ctx, "LISTEN outbox_events"); err != nil {
		if ctx.Err() != nil {
			return
		}
		p.logger.Fatal().
			Err(oops.In("publisher").Code(ErrCodeListenFailed).Wrap(err)).
			Msg("failed to LISTEN on outbox_events")
		return
	}
	p.logger.Info().Msg("LISTEN outbox_events active")

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Initial cycle on startup (catch rows inserted while we were locking).
	p.publishCycle(ctx)

	notifyCh := make(chan struct{}, 1)
	go func() {
		for {
			_, err := listenConn.Conn().WaitForNotification(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				p.logger.Warn().Err(err).Msg("LISTEN error; retrying in 5s")
				time.Sleep(5 * time.Second)
				continue
			}
			select {
			case notifyCh <- struct{}{}:
			default:
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-notifyCh:
			p.publishCycle(ctx)
		case <-ticker.C:
			p.publishCycle(ctx)
		}
	}
}

type outboxRow struct {
	Seq     int64
	Subject string
	Payload []byte
}

// publishCycle fetches up to fetchBatchSize unpublished rows in seq order
// and publishes each one. Stops on the first publish error to preserve
// ordering.
func (p *publisher) publishCycle(ctx context.Context) {
	rows, err := p.fetchUnpublished(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		p.logger.Error().Err(err).Msg("outbox fetch failed")
		return
	}
	for _, row := range rows {
		if ctx.Err() != nil {
			return
		}
		if err := p.publishRow(ctx, row); err != nil {
			p.logger.Error().
				Err(err).
				Int64("seq", row.Seq).
				Str("subject", row.Subject).
				Msg("outbox publish failed; will retry on next cycle")
			return
		}
	}
}

func (p *publisher) fetchUnpublished(ctx context.Context) ([]outboxRow, error) {
	const q = `
		SELECT seq, subject, payload
		  FROM outbox.events
		 WHERE published_at IS NULL
		 ORDER BY seq
		 LIMIT $1`
	pgRows, err := p.pool.Query(ctx, q, fetchBatchSize)
	if err != nil {
		return nil, oops.In("publisher").Code(ErrCodeOutboxFetchFailed).Wrap(err)
	}
	defer pgRows.Close()

	var out []outboxRow
	for pgRows.Next() {
		var r outboxRow
		if err := pgRows.Scan(&r.Seq, &r.Subject, &r.Payload); err != nil {
			return nil, oops.In("publisher").Code(ErrCodeOutboxFetchFailed).Wrap(err)
		}
		out = append(out, r)
	}
	return out, pgRows.Err()
}

func (p *publisher) publishRow(ctx context.Context, row outboxRow) error {
	// js.Publish waits for server ack by default — safe to mark as published
	// immediately after this returns nil.
	if _, err := p.js.Publish(ctx, row.Subject, row.Payload); err != nil {
		return oops.In("publisher").
			Code(ErrCodeOutboxPublishFailed).
			With("seq", row.Seq).
			With("subject", row.Subject).
			Wrap(err)
	}
	if _, err := p.pool.Exec(ctx,
		`UPDATE outbox.events SET published_at = now() WHERE seq = $1`, row.Seq,
	); err != nil {
		// The message is in NATS but we failed to mark it. Log and continue;
		// the duplicate will be deduplicated by the consumer (at-least-once).
		p.logger.Error().
			Err(oops.In("publisher").Code(ErrCodeOutboxMarkFailed).Wrap(err)).
			Int64("seq", row.Seq).
			Msg("failed to mark outbox row published; possible duplicate delivery")
	}
	return nil
}
