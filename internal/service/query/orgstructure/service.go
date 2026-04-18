// Package orgstructure exposes read methods over the orgstructure
// projections (projections.organizations, projections.clinics,
// projections.departments). Queries are plain SQL (tx.Raw(...)) — no
// gorm models, no AutoMigrate, no First/Find/Count — to keep the
// read side decoupled from the command-side gorm tag layout.
//
// Transport (gRPC/HTTP) is plugged in by a higher layer; readers take
// plain Go inputs and return per-query view structs directly.
package orgstructure

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Bounds applied to List pagination. Readers reject inputs outside
// these ranges with a domain error.
const (
	listMinLimit     = 1
	listMaxLimit     = 500
	listDefaultLimit = 50
)

// OrganizationReader exposes read methods for the Organization
// projection.
type OrganizationReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationReader returns an OrganizationReader bound to the
// given gorm DB.
func NewOrganizationReader(db *gorm.DB, logger *zerolog.Logger) *OrganizationReader {
	return &OrganizationReader{db: db, logger: logger}
}

// ClinicReader exposes read methods for the Clinic projection.
type ClinicReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewClinicReader returns a ClinicReader bound to the given gorm DB.
func NewClinicReader(db *gorm.DB, logger *zerolog.Logger) *ClinicReader {
	return &ClinicReader{db: db, logger: logger}
}

// DepartmentReader exposes read methods for the Department projection.
type DepartmentReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewDepartmentReader returns a DepartmentReader bound to the given
// gorm DB.
func NewDepartmentReader(db *gorm.DB, logger *zerolog.Logger) *DepartmentReader {
	return &DepartmentReader{db: db, logger: logger}
}
