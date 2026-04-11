// Package outbox implements the transactional outbox pattern. Records
// are written in the same transaction as the aggregate state. A separate
// publisher service drains the table and ships rows to NATS JetStream.
package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

// Error codes emitted by Store implementations.
const (
	ErrCodeAppendFailed = "outbox_append_failed"
	ErrCodeNilRecord    = "outbox_nil_record"
)

// Record is the infrastructure form of one outbox row. It mirrors a
// medincident.events.v1.Envelope: every envelope field is a column;
// the inner proto message is marshaled into a google.protobuf.Any whose
// bytes go into Payload. Headers carry transport-level metadata such as
// the JetStream Nats-Msg-Id used for deduplication.
type Record struct {
	EventID       uuid.UUID
	OccurredAt    time.Time
	AggregateType string
	AggregateID   string
	CorrelationID string
	Subject       string
	Headers       map[string]string
	Payload       []byte
	CreatedAt     time.Time
}

// Store persists outbox rows. The transaction is an explicit parameter
// (not pulled from ctx) because an outbox write MUST share a
// transaction with the aggregate write; making it explicit enforces
// that invariant at the call site.
type Store interface {
	Append(ctx context.Context, t tx.Tx, record *Record) error
}
