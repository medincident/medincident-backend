package organizationapp

import "github.com/google/uuid"

// CreateCommand is the input of Service.Create.
type CreateCommand struct {
	// Name is the organization's display name. Trimmed and validated
	// inside organization.New.
	Name string
	// Description is optional. Empty string means "not set".
	Description string
	// LegalAddress is optional. nil means "not set".
	LegalAddress *AddressInput
}

// CreateResult is the output of Service.Create.
type CreateResult struct {
	ID uuid.UUID
}

// RenameCommand is the input of Service.Rename.
type RenameCommand struct {
	ID   uuid.UUID
	Name string
}

// RenameResult is the output of Service.Rename.
type RenameResult struct{}

// UpdateDescriptionCommand is the input of Service.UpdateDescription.
// An empty Description clears the description.
type UpdateDescriptionCommand struct {
	ID          uuid.UUID
	Description string
}

// UpdateDescriptionResult is the output of Service.UpdateDescription.
type UpdateDescriptionResult struct{}

// RelocateLegalAddressCommand is the input of Service.RelocateLegalAddress.
// Address may be nil, which removes the legal address from the aggregate.
type RelocateLegalAddressCommand struct {
	ID      uuid.UUID
	Address *AddressInput
}

// RelocateLegalAddressResult is the output of Service.RelocateLegalAddress.
type RelocateLegalAddressResult struct{}

// AddressInput is the transport-agnostic carrier for a legal address
// input coming from a command. The application service builds a
// domain geo.Address from it, which is the boundary where the input
// is validated (via geo.NewAddress and geo.NewPoint).
type AddressInput struct {
	Text  string
	Point *PointInput
}

// PointInput carries optional geographic coordinates.
type PointInput struct {
	Longitude float64
	Latitude  float64
}
