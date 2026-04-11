package organizationapp

import (
	"context"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
)

// Error codes emitted by Repository implementations. Defined here, next
// to the contract, so that adapters and consumers reference one symbol.
// The string values are unprefixed because oops.In("storage.postgres")
// already provides the namespace.
const (
	ErrCodeNotFound   = "organization_not_found"
	ErrCodeFindFailed = "organization_find_failed"
	ErrCodeSaveFailed = "organization_save_failed"
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
	// Find loads an organization by primary key. Returns a repository-
	// level not-found error when the row does not exist.
	Find(ctx context.Context, id uuid.UUID) (*organization.Organization, error)

	// Save performs an upsert of the aggregate's current state.
	Save(ctx context.Context, o *organization.Organization) error
}
