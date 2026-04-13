package orgstructure

import (
	"encoding/json"

	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

// Error codes emitted by AppendOutboxEvent.
const (
	ErrCodeOutboxMarshalFailed = "outbox_marshal_failed"
	ErrCodeOutboxAppendFailed  = "outbox_append_failed"
)

// AppendOutboxEvent serialises the envelope and writes one row into
// outbox.events inside the given transaction. The caller is responsible
// for building the Envelope (occurred_at, aggregate_type, aggregate_id,
// payload). headers may be nil; it is marshaled to `{}` either way.
//
// This helper exists purely so the three-line proto.Marshal +
// json.Marshal + tx.Create sequence is written once, not nine times.
func AppendOutboxEvent(
	tx *gorm.DB,
	subject string,
	envelope *envelopev1.Envelope,
	headers map[string]string,
) error {
	payload, err := proto.Marshal(envelope)
	if err != nil {
		return oops.In("services.orgstructure.outbox").
			Code(ErrCodeOutboxMarshalFailed).
			With("subject", subject).
			Wrap(err)
	}
	if headers == nil {
		headers = map[string]string{}
	}
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return oops.In("services.orgstructure.outbox").
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
		return oops.In("services.orgstructure.outbox").
			Code(ErrCodeOutboxAppendFailed).
			With("subject", subject).
			Wrap(err)
	}
	return nil
}
