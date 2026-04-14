// Package outbox provides the shared outbox-append helper used by
// every command service. It replaces the per-package duplicates that
// previously lived in services/orgstructure and services/membership.
package outbox

import (
	"encoding/json"

	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

// Error codes emitted by AppendEvent.
const (
	ErrCodeOutboxMarshalFailed = "outbox_marshal_failed"
	ErrCodeOutboxAppendFailed  = "outbox_append_failed"
)

// AppendEvent serialises the envelope and writes one row into
// outbox.events inside the given transaction. Callers own the
// envelope construction (occurred_at, aggregate_type, aggregate_id,
// payload); this helper only handles the marshal + insert.
func AppendEvent(
	tx *gorm.DB,
	subject string,
	envelope *envelopev1.Envelope,
	headers map[string]string,
) error {
	payload, err := proto.Marshal(envelope)
	if err != nil {
		return oops.In("services.outbox").
			Code(ErrCodeOutboxMarshalFailed).
			With("subject", subject).
			Wrap(err)
	}
	if headers == nil {
		headers = map[string]string{}
	}
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return oops.In("services.outbox").
			Code(ErrCodeOutboxMarshalFailed).
			With("subject", subject).
			Wrap(err)
	}
	row := model.OutboxEvent{
		Subject: subject,
		Payload: payload,
		Headers: datatypes.JSON(headersJSON),
	}
	if err := tx.Create(&row).Error; err != nil {
		return oops.In("services.outbox").
			Code(ErrCodeOutboxAppendFailed).
			With("subject", subject).
			Wrap(err)
	}
	return nil
}
