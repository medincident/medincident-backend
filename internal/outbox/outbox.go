// Package outbox provides the transactional outbox write helper used by
// every command service. Append serialises an Envelope proto and inserts
// a row into outbox.events inside the caller-supplied gorm transaction.
// The publisher-server reads unpublished rows and forwards them to NATS
// JetStream.
package outbox

import (
	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

const ErrCodeOutboxAppendFailed = "outbox_append_failed"

// Append serialises env and inserts a row into outbox.events inside tx.
// Must be called from within an active gorm transaction on the command DB.
// subject follows the scheme medincident.event.<aggregate_type>.v1.<action>.
func Append(tx *gorm.DB, subject string, env *eventv1.Envelope) error {
	payload, err := proto.Marshal(env)
	if err != nil {
		return oops.In("outbox").
			Code(ErrCodeOutboxAppendFailed).
			With("subject", subject).
			Wrap(err)
	}
	if err := tx.Exec(
		`INSERT INTO outbox.events (subject, payload) VALUES (?, ?)`,
		subject, payload,
	).Error; err != nil {
		return oops.In("outbox").
			Code(ErrCodeOutboxAppendFailed).
			With("subject", subject).
			Wrap(err)
	}
	return nil
}
