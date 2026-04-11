package organization

import (
	"time"

	"github.com/google/uuid"

	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

// Domain events are plain Go structs with no struct tags: the outbox
// layer marshals them via encoding/json using exported field names
// directly. Renaming an exported field is a BREAKING CHANGE for any
// outbox rows already persisted under the old shape and must be paired
// with an outbox drain or a payload-rewrite migration.
//
// Events carry the NEW state only. Consumers compute diffs from their
// own prior projection.

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
