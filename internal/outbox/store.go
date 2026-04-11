// Package outbox implements the transactional outbox pattern. The write
// side (Publish + Store) is in this package; a separate publisher binary
// will consume rows via LISTEN/NOTIFY + poll fallback and emit them to
// NATS JetStream (out of scope for this milestone).
package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/tx"
)

// Record is the infrastructure form of one outbox row.
//
// AggregateType and AggregateID are captured from the EventSource at
// Publish time and become envelope.aggregate_type / envelope.aggregate_id
// on the wire. Both are transport-level identifiers and are stored as
// strings here (the domain keeps its own uuid.UUID types; the boundary
// between domain and outbox is the string conversion).
//
// EventType is the stable identifier used to look up EventInfo for
// deserialisation.
//
// Payload is json.Marshal of the domain event struct (no proto
// involvement at write time).
type Record struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	CreatedAt     time.Time
}

// Store persists outbox rows. The transaction is an explicit parameter
// (not picked up from ctx like regular repositories) because an outbox
// write MUST share a transaction with the aggregate write. Making it
// explicit forces the invariant at compile time: a caller cannot pass
// nil or forget the tx.
type Store interface {
	Append(ctx context.Context, t tx.Tx, record Record) error
}
