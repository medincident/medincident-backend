// Package classifier owns the incident classifier command-side logic:
// two aggregates (IncidentCategory, IncidentType) that together model
// a per-organisation, hierarchical classifier with max depth 5 and
// active-name uniqueness scoped per organisation. Two service structs
// share the package: category mutations and type mutations. They hold
// nothing beyond a *gorm.DB and a *zerolog.Logger.
package classifier

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/service/authz"
)

// IncidentCategoryService handles mutations of domain.incident_categories,
// including all cascade and hierarchy invariants.
type IncidentCategoryService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewIncidentCategoryService returns a service bound to the given gorm
// DB and authorization service.
func NewIncidentCategoryService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *IncidentCategoryService {
	return &IncidentCategoryService{db: db, authz: az, logger: logger}
}

// IncidentTypeService handles mutations of domain.incident_types.
type IncidentTypeService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewIncidentTypeService returns a service bound to the given gorm DB
// and authorization service.
func NewIncidentTypeService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *IncidentTypeService {
	return &IncidentTypeService{db: db, authz: az, logger: logger}
}
