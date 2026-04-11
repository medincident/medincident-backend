package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/tx"
)

// Error codes emitted by Publish.
const (
	CodeNoMapper           = "outbox_no_mapper"
	CodeJSONMarshalFailed  = "outbox_json_marshal_failed"
	CodeIDGenerationFailed = "outbox_id_generation_failed"
	CodeStoreAppendFailed  = "outbox_store_append_failed"
)

// EventSource is the contract Publish expects from each aggregate source.
//
// AggregateType is a short BC-level name (e.g. "organization") that
// becomes envelope.aggregate_type on the wire. AggregateID is the
// aggregate's own id rendered as a string — the domain keeps its uuid.UUID
// type internally but crosses into outbox-land as a string because
// envelope.aggregate_id, the outbox column, and every downstream consumer
// treat it as opaque text. Both values are captured ONCE per source
// inside Publish and applied to every outbox row produced for that source.
// PullEvents drains and clears the event buffer (inherited from embedded
// aggregate.Root).
type EventSource interface {
	AggregateType() string
	AggregateID() string
	PullEvents() []any
}

// Publish collects all domain events from the given sources and writes
// them as JSONB rows into the outbox store. The transaction is REQUIRED
// and passed explicitly — outbox writes must share a transaction with
// aggregate writes, and the function signature enforces that.
func Publish(
	ctx context.Context,
	t tx.Tx,
	store Store,
	registry Registry,
	sources ...EventSource,
) error {
	now := time.Now().UTC()

	for _, src := range sources {
		aggType := src.AggregateType()
		aggID := src.AggregateID()

		for _, ev := range src.PullEvents() {
			info, ok := registry.ByGoType(reflect.TypeOf(ev))
			if !ok {
				return oops.In("outbox").
					Code(CodeNoMapper).
					With("go_type", fmt.Sprintf("%T", ev)).
					With("aggregate_type", aggType).
					With("aggregate_id", aggID).
					Hint("register an EventInfo in feature/infra/RegisterOutboxMappers").
					Errorf("no EventInfo for event")
			}

			payload, err := json.Marshal(ev)
			if err != nil {
				return oops.In("outbox").
					Code(CodeJSONMarshalFailed).
					With("type_name", info.TypeName).
					With("aggregate_type", aggType).
					With("aggregate_id", aggID).
					Wrap(err)
			}

			id, err := uuid.NewV7()
			if err != nil {
				return oops.In("outbox").
					Code(CodeIDGenerationFailed).
					Wrap(err)
			}

			rec := &Record{
				ID:            id,
				AggregateType: aggType,
				AggregateID:   aggID,
				EventType:     info.TypeName,
				Payload:       payload,
				CreatedAt:     now,
			}
			if err := store.Append(ctx, t, rec); err != nil {
				return oops.In("outbox").
					Code(CodeStoreAppendFailed).
					With("type_name", info.TypeName).
					With("aggregate_type", aggType).
					With("aggregate_id", aggID).
					Wrap(err)
			}
		}
	}
	return nil
}
