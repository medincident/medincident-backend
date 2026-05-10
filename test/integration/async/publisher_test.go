//go:build integration

package async_integration_test

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
)

// runTestPublisher is a simplified publisher for integration tests.
// It polls outbox.events on the command DB every 50ms and publishes to NATS.
// No advisory lock — single process in tests.
func runTestPublisher(ctx context.Context, dsn string, js jetstream.JetStream) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return
	}
	defer pool.Close()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rows, err := pool.Query(ctx,
				`SELECT seq, subject, payload FROM outbox.events WHERE published_at IS NULL ORDER BY seq LIMIT 100`)
			if err != nil {
				continue
			}
			type row struct {
				seq     int64
				subject string
				payload []byte
			}
			var batch []row
			for rows.Next() {
				var r row
				_ = rows.Scan(&r.seq, &r.subject, &r.payload)
				batch = append(batch, r)
			}
			rows.Close()
			for _, r := range batch {
				if _, err := js.Publish(ctx, r.subject, r.payload); err != nil {
					break
				}
				_, _ = pool.Exec(ctx, `UPDATE outbox.events SET published_at = now() WHERE seq = $1`, r.seq)
			}
		}
	}
}
