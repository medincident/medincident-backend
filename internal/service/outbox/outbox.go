// Package outbox provides the shared outbox-append helper used by
// every command service. It replaces the per-package duplicates that
// previously lived in services/orgstructure and services/membership.
//
// Deduplication keys for NATS are owned by the publisher service and
// derived from outbox.events.id (stable because the column is
// GENERATED ALWAYS AS IDENTITY). The command service does not write
// any dedup metadata — if you find yourself wanting to, add it in the
// publisher instead.
package outbox

import (
	"time"

	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

// Error codes emitted by Publish.
const (
	ErrCodeOutboxMarshalFailed = "outbox_marshal_failed"
	ErrCodeOutboxAppendFailed  = "outbox_append_failed"
)

// Publish builds an envelope around `event`, marshals it, and writes
// one row into outbox.events inside the given transaction. It is the
// single entry point for every command service — AGENTS.md §10 says
// the event proto is still constructed inline at the call site, but
// the envelope assembly, proto marshal, and INSERT all live here.
func Publish(
	tx *gorm.DB,
	subject, aggregateType, aggregateID string,
	occurredAt time.Time,
	event proto.Message,
) error {
	payload, err := anypb.New(event)
	if err != nil {
		return oops.In("services.outbox").
			Code(ErrCodeOutboxMarshalFailed).
			With("subject", subject).
			With("aggregate_type", aggregateType).
			With("aggregate_id", aggregateID).
			Wrap(err)
	}
	envelope := &envelopev1.Envelope{
		OccurredAt:    timestamppb.New(occurredAt),
		AggregateType: aggregateType,
		AggregateId:   aggregateID,
		Payload:       payload,
	}
	envelopeBytes, err := proto.Marshal(envelope)
	if err != nil {
		return oops.In("services.outbox").
			Code(ErrCodeOutboxMarshalFailed).
			With("subject", subject).
			Wrap(err)
	}
	row := model.OutboxEvent{
		Subject: subject,
		Payload: envelopeBytes,
	}
	if err := tx.Create(&row).Error; err != nil {
		return oops.In("services.outbox").
			Code(ErrCodeOutboxAppendFailed).
			With("subject", subject).
			Wrap(err)
	}
	return nil
}
