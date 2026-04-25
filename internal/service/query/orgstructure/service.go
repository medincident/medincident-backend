// Package orgstructure exposes read methods over the orgstructure
// projections (projections.organizations, projections.clinics,
// projections.departments). Queries are plain SQL (tx.Raw(...)) — no
// gorm models, no AutoMigrate, no First/Find/Count — to keep the
// read side decoupled from the command-side gorm tag layout.
//
// Transport (gRPC/HTTP) is plugged in by a higher layer; readers take
// plain Go inputs and return per-query view structs directly.
//
// Authorization model: organizations are the public catalog of the
// platform — any authenticated caller can list and fetch them, so
// OrganizationReader does not consult authz. Everything below the
// organization (clinics, departments) is scoped via authz.ReaderOf.X:
// system admins, organization admins, and members of the containing
// organization may read. A caller outside that set receives
// permission_denied with no distinction between "scope does not
// exist" and "you are not authorized" — see the package authz
// preamble for the rationale.
package orgstructure

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// OrganizationReader exposes read methods for the Organization
// projection. Organizations are public to any authenticated caller, so
// this reader takes no *authz.Authz — authentication is enforced by
// the gRPC interceptor chain upstream.
type OrganizationReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationReader returns an OrganizationReader bound to the
// given gorm DB.
//
// Organizations are a public catalog for authenticated callers.
// OrganizationReader performs no authorization/scope check; authentication
// is enforced upstream by the gRPC interceptor chain (see AGENTS.md § Authorization model).
func NewOrganizationReader(db *gorm.DB, logger *zerolog.Logger) *OrganizationReader {
	return &OrganizationReader{db: db, logger: logger}
}

// ClinicReader exposes read methods for the Clinic projection. Reads
// are gated by authz.ReaderOf.{Clinic,Organization}, so the caller's
// identity must be supplied on every method.
type ClinicReader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewClinicReader returns a ClinicReader bound to the given gorm DB
// and authorization service.
func NewClinicReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *ClinicReader {
	return &ClinicReader{db: db, authz: az, logger: logger}
}

// DepartmentReader exposes read methods for the Department projection.
// Reads are gated by authz.ReaderOf.{Department,Clinic}.
type DepartmentReader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewDepartmentReader returns a DepartmentReader bound to the given
// gorm DB and authorization service.
func NewDepartmentReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *DepartmentReader {
	return &DepartmentReader{db: db, authz: az, logger: logger}
}
