package orgstructure

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// OrganizationService handles mutations of domain.organizations.
type OrganizationService struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrganizationService returns a service bound to the given gorm DB.
func NewOrganizationService(db *gorm.DB, logger *zerolog.Logger) *OrganizationService {
	return &OrganizationService{db: db, logger: logger}
}

// ClinicService handles mutations of domain.clinics.
type ClinicService struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewClinicService returns a service bound to the given gorm DB.
func NewClinicService(db *gorm.DB, logger *zerolog.Logger) *ClinicService {
	return &ClinicService{db: db, logger: logger}
}

// DepartmentService handles mutations of domain.departments.
type DepartmentService struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewDepartmentService returns a service bound to the given gorm DB.
func NewDepartmentService(db *gorm.DB, logger *zerolog.Logger) *DepartmentService {
	return &DepartmentService{db: db, logger: logger}
}
