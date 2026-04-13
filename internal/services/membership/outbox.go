package membership

import (
	"encoding/json"

	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

const (
	ErrCodeOutboxMarshalFailed = "membership_outbox_marshal_failed"
	ErrCodeOutboxAppendFailed  = "membership_outbox_append_failed"
)

// AppendMembershipOutboxEvent serialises the envelope and writes one
// row into outbox.events inside the given transaction. Mirrors the
// helper in services/orgstructure/outbox.go.
func AppendMembershipOutboxEvent(
	tx *gorm.DB,
	subject string,
	envelope *envelopev1.Envelope,
	headers map[string]string,
) error {
	payload, err := proto.Marshal(envelope)
	if err != nil {
		return oops.In("services.membership.outbox").
			Code(ErrCodeOutboxMarshalFailed).
			With("subject", subject).
			Wrap(err)
	}
	if headers == nil {
		headers = map[string]string{}
	}
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return oops.In("services.membership.outbox").
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
		return oops.In("services.membership.outbox").
			Code(ErrCodeOutboxAppendFailed).
			With("subject", subject).
			Wrap(err)
	}
	return nil
}
