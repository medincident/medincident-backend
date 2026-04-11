package organization

import (
	"time"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

// Domain events are plain Go structs with NO struct tags. The outbox
// layer marshals them via encoding/json using exported field names
// directly; the resulting JSONB keys are CamelCase (e.g. "NewName",
// "LegalAddress"). Because the Go type is the sole source of truth for
// both write (outbox.Publish) and read (publisher) sides, there is no
// cross-party contract to negotiate.
//
// Caveat: renaming a Go field on a domain event is a BREAKING CHANGE for
// already-persisted outbox rows. Any such rename must be accompanied by
// a drain of the outbox (all rows published before the new code ships)
// or a data migration that rewrites existing payloads.

// Events describe the NEW state after a mutation. They never carry the
// "before" value — a consumer that needs a diff has its own previous
// state to compare against. This keeps payloads small and avoids
// duplicating information that is already implicit in the consumer's
// local view of the aggregate.

// Created is raised when an Organization is first constructed via New.
type Created struct {
	ID           uuid.UUID
	Name         string
	Description  string
	LegalAddress *geo.Address
	At           time.Time
}

// Renamed is raised when Organization.Rename successfully changes the name.
type Renamed struct {
	ID   uuid.UUID
	Name string
	At   time.Time
}

// DescriptionUpdated is raised when Organization.UpdateDescription
// successfully replaces the description (including clearing it).
type DescriptionUpdated struct {
	ID          uuid.UUID
	Description string
	At          time.Time
}

// LegalAddressRelocated is raised when Organization.RelocateLegalAddress
// successfully replaces the legal address (including removing it via nil).
type LegalAddressRelocated struct {
	ID           uuid.UUID
	LegalAddress *geo.Address
	At           time.Time
}
