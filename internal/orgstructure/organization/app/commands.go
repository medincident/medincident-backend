// Package organizationapp is the application layer for the Organization
// aggregate. It exposes commands (XCommand structs), results (XResult
// structs), the Repository port, and the Service that orchestrates
// aggregate loading, mutation, persistence, and outbox publishing.
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
	ID      uuid.UUID
	NewName string
}

// RenameResult is the output of Service.Rename. Empty for now; returned
// as a typed struct anyway so the signature is stable across future
// additions (e.g., an updated_at timestamp).
type RenameResult struct{}

// UpdateDescriptionCommand is the input of Service.UpdateDescription.
// An empty NewDescription clears the description.
type UpdateDescriptionCommand struct {
	ID             uuid.UUID
	NewDescription string
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
	// Text is the free-form postal text. Required when AddressInput is
	// present; validated by geo.NewAddress.
	Text string
	// Point is optional. nil means "no coordinates".
	Point *PointInput
}

// PointInput carries optional geographic coordinates. Validated by
// geo.NewPoint when present.
type PointInput struct {
	Longitude float64
	Latitude  float64
}
