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
// AggregateType and AggregateID are captured once per source and stamped
// on every row produced for it; PullEvents drains the event buffer
// (inherited from embedded aggregate.Root).
type EventSource interface {
	AggregateType() string
	AggregateID() string
	PullEvents() []any
}

// Publish drains every source's event buffer and writes one outbox row
// per event inside t. The transaction must be the same one used for the
// aggregate write.
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
