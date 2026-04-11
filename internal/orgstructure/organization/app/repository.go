package organizationapp

import (
	"context"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
)

// Repository is the application-layer port for persistence of
// Organization aggregates. Implementations live in the storage adapter
// (internal/storage/postgres) and pick up a transaction from ctx via
// tx.FromContext when one is present.
//
// The interface is deliberately context-only — it does not accept a
// tx.Tx parameter. Transactions flow through ctx so that the same
// repository method works both inside a Within block (transactional
// writes) and outside (one-shot reads).
type Repository interface {
	// GetByID loads an organization by primary key. Returns a
	// repository-level not-found error when the row does not exist.
	GetByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error)

	// Save performs an upsert of the aggregate's current state. Callers
	// are expected to pass the same aggregate value they mutated inside
	// a Within block; Save always writes the full row.
	Save(ctx context.Context, o *organization.Organization) error

	// List returns organizations matching the filter in created_at DESC
	// order. The filter's zero value means "no filter".
	List(ctx context.Context, filter ListFilter) ([]*organization.Organization, error)
}

// ListFilter carries optional listing criteria. All fields are optional;
// the zero value returns (up to Limit) organizations with no filter.
type ListFilter struct {
	// NameLike is a case-insensitive substring match on Name. Empty
	// means "no name filter".
	NameLike string
	// Limit caps the number of rows returned. 0 means "use default".
	Limit int32
	// Offset is the row offset. 0 means "from the start".
	Offset int32
}
