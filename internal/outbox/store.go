// Package outbox implements the transactional outbox pattern. The write
// side (Publish + Store) is in this package; a separate publisher binary
// will consume rows via LISTEN/NOTIFY + poll fallback and emit them to
// NATS JetStream (out of scope for this milestone).
package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

// Record is the infrastructure form of one outbox row. AggregateType
// and AggregateID cross from the domain as strings; EventType is the
// stable identifier used to look up EventInfo at read time; Payload is
// json.Marshal of the domain event struct (no proto at write time).
type Record struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   string
	EventType     string
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
