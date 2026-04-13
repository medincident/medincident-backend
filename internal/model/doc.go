// Package model holds the gorm-backed persistence types for the
// command service. Only this package is allowed to import
// github.com/guregu/null/v6 — callers convert at the service boundary
// via null.StringFromPtr / null.FloatFromPtr / etc.
//
// Every field carries an explicit gorm field-level permission tag so
// that immutability at the Go layer is enforced at the persistence
// boundary:
//
//   - primary keys and creation-time metadata use <-:create
//   - parent foreign keys (Clinic.OrganizationID, Department.ClinicID)
//     use <-:create because parenthood is immutable per the spec
//   - mutable body fields use the default <- written explicitly for clarity
//   - OutboxEvent.PublishedAt is marked `-` — the publisher service
//     owns that column; command-service never touches it
package model
