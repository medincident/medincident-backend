package postgres

import (
	"context"

	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/tx"
)

// Error codes emitted by OutboxStore.Append.
const (
	CodeOutboxInsertFailed = "postgres_outbox_insert_failed"
	CodeOutboxNilRecord    = "postgres_outbox_nil_record"
)

// OutboxStore persists outbox rows into outbox.events. Stateless — the
// transaction passed to Append already carries the connection.
type OutboxStore struct{}

// NewOutboxStore returns a new OutboxStore.
func NewOutboxStore() *OutboxStore { return &OutboxStore{} }

const insertOutboxSQL = `
INSERT INTO outbox.events (
    id, aggregate_type, aggregate_id, event_type, payload, created_at
)
VALUES ($1, $2, $3, $4, $5, $6)
`

// Append inserts one outbox row in the given transaction.
func (s *OutboxStore) Append(ctx context.Context, t tx.Tx, rec *outbox.Record) error {
	if rec == nil {
		return oops.In("storage.postgres").
			Code(CodeOutboxNilRecord).
			Errorf("nil outbox record")
	}
	pg, err := unwrap(t)
	if err != nil {
		return err
	}
	if _, err := pg.raw.Exec(ctx, insertOutboxSQL,
		rec.ID,
		rec.AggregateType,
		rec.AggregateID,
		rec.EventType,
		rec.Payload,
		rec.CreatedAt,
	); err != nil {
		return oops.In("storage.postgres").
			Code(CodeOutboxInsertFailed).
			With("id", rec.ID).
			With("aggregate_type", rec.AggregateType).
			With("aggregate_id", rec.AggregateID).
			With("event_type", rec.EventType).
			Wrap(err)
	}
	return nil
}

var _ outbox.Store = (*OutboxStore)(nil)
