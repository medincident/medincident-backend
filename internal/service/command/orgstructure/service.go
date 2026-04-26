// Package orgstructure is the write-side service for organizational structure.
//
// See: docs/services/OrgStructure.md
package orgstructure

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// OrganizationService handles mutations of domain.organizations.
type OrganizationService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewOrganizationService returns a service bound to the given gorm DB
// and authorization service.
func NewOrganizationService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *OrganizationService {
	return &OrganizationService{db: db, authz: az, logger: logger}
}

// ClinicService handles mutations of domain.clinics.
type ClinicService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewClinicService returns a service bound to the given gorm DB and
// authorization service.
func NewClinicService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *ClinicService {
	return &ClinicService{db: db, authz: az, logger: logger}
}

// DepartmentService handles mutations of domain.departments.
type DepartmentService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewDepartmentService returns a service bound to the given gorm DB
// and authorization service.
func NewDepartmentService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *DepartmentService {
	return &DepartmentService{db: db, authz: az, logger: logger}
}
