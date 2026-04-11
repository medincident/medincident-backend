package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/tx"
	sqlcgen "github.com/medincident/medincident-command-service/internal/storage/postgres/sqlc/gen"
)

// OutboxStore persists outbox rows into outbox.events via the
// sqlc-generated AppendOutboxEvent query.
type OutboxStore struct{}

// NewOutboxStore returns a new OutboxStore.
func NewOutboxStore() *OutboxStore { return &OutboxStore{} }

// Append inserts one outbox row inside the given transaction.
func (s *OutboxStore) Append(ctx context.Context, t tx.Tx, rec *outbox.Record) error {
	if rec == nil {
		return oops.In("storage.postgres").
			Code(outbox.ErrCodeNilRecord).
			Errorf("nil outbox record")
	}
	pg, err := unwrap(t)
	if err != nil {
		return err
	}

	headersJSON, err := json.Marshal(rec.Headers)
	if err != nil {
		return oops.In("storage.postgres").
			Code(outbox.ErrCodeAppendFailed).
			With("subject", rec.Subject).
			Wrap(err)
	}

	q := sqlcgen.New(pg.raw)
	if err := q.AppendOutboxEvent(ctx, sqlcgen.AppendOutboxEventParams{
		EventID:       uuidToPg(rec.EventID),
		OccurredAt:    pgtype.Timestamptz{Time: rec.OccurredAt, Valid: true},
		AggregateType: rec.AggregateType,
		AggregateID:   rec.AggregateID,
		CorrelationID: rec.CorrelationID,
		Subject:       rec.Subject,
		Headers:       headersJSON,
		Payload:       rec.Payload,
		CreatedAt:     pgtype.Timestamptz{Time: rec.CreatedAt, Valid: true},
	}); err != nil {
		return oops.In("storage.postgres").
			Code(outbox.ErrCodeAppendFailed).
			With("event_id", rec.EventID).
			With("subject", rec.Subject).
			With("aggregate_type", rec.AggregateType).
			With("aggregate_id", rec.AggregateID).
			Wrap(err)
	}
	return nil
}

var _ outbox.Store = (*OutboxStore)(nil)
