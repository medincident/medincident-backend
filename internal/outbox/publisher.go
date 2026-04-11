package outbox

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/medincident/medincident-command-service/internal/shared/clock"
	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

// Error codes emitted by Publisher.Publish.
const (
	ErrCodeNoMapper     = "no_mapper"
	ErrCodeProtoConvert = "proto_convert_failed"
	ErrCodeAnyWrap      = "any_wrap_failed"
	ErrCodeProtoMarshal = "proto_marshal_failed"
	ErrCodeStoreAppend  = "store_append_failed"
)

// natsMsgIDHeader is the JetStream deduplication header. The publisher
// service uses jetstream.WithMsgID(headers[natsMsgIDHeader]) when
// emitting the row.
const natsMsgIDHeader = "Nats-Msg-Id"

// EventSource is the contract Publish expects from each aggregate
// source. AggregateType / AggregateID are stamped on every row produced
// from the source; PullEvents drains the aggregate's event buffer.
type EventSource interface {
	AggregateType() string
	AggregateID() string
	PullEvents() []any
}

// Publisher writes outbox rows for one or more aggregate event sources
// inside a caller-supplied transaction.
type Publisher interface {
	Publish(ctx context.Context, t tx.Tx, sources ...EventSource) error
}

type publisher struct {
	store    Store
	registry Registry
	clock    clock.Clock
}

var _ Publisher = (*publisher)(nil)

// NewPublisher wires a Publisher implementation. Returns the concrete
// type by design — the DI factory adapts it to the Publisher interface.
func NewPublisher(store Store, registry Registry, clk clock.Clock) *publisher {
	return &publisher{store: store, registry: registry, clock: clk}
}

// Publish drains every source's event buffer and writes one row per
// event. The transaction MUST be the same one used for the aggregate
// write; that's the whole point of the outbox pattern.
func (p *publisher) Publish(ctx context.Context, t tx.Tx, sources ...EventSource) error {
	now := p.clock.Now().UTC()

	for _, src := range sources {
		aggType := src.AggregateType()
		aggID := src.AggregateID()

		for _, ev := range src.PullEvents() {
			info, ok := p.registry.Lookup(reflect.TypeOf(ev))
			if !ok {
				return oops.In("outbox").
					Code(ErrCodeNoMapper).
					With("go_type", fmt.Sprintf("%T", ev)).
					With("aggregate_type", aggType).
					With("aggregate_id", aggID).
					Hint("register an EventInfo in feature/infra/RegisterOutboxMappers").
					Errorf("no EventInfo for event")
			}

			msg, err := info.ToProto(ev)
			if err != nil {
				return oops.In("outbox").
					Code(ErrCodeProtoConvert).
					With("subject", info.Subject).
					With("aggregate_type", aggType).
					With("aggregate_id", aggID).
					Wrap(err)
			}

			anyMsg, err := anypb.New(msg)
			if err != nil {
				return oops.In("outbox").
					Code(ErrCodeAnyWrap).
					With("subject", info.Subject).
					Wrap(err)
			}

			payload, err := proto.Marshal(anyMsg)
			if err != nil {
				return oops.In("outbox").
					Code(ErrCodeProtoMarshal).
					With("subject", info.Subject).
					Wrap(err)
			}

			eventID := uuid.New()
			rec := &Record{
				EventID:       eventID,
				OccurredAt:    now,
				AggregateType: aggType,
				AggregateID:   aggID,
				Subject:       info.Subject,
				Headers:       map[string]string{natsMsgIDHeader: eventID.String()},
				Payload:       payload,
				CreatedAt:     now,
			}
			if err := p.store.Append(ctx, t, rec); err != nil {
				return oops.In("outbox").
					Code(ErrCodeStoreAppend).
					With("subject", info.Subject).
					With("aggregate_type", aggType).
					With("aggregate_id", aggID).
					Wrap(err)
			}
		}
	}
	return nil
}
