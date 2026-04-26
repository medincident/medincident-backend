// Package announcement is the write-side service for announcements.
// Announcements are one-way informational broadcasts from administration
// to employees, scoped to an organization, clinic, or department.
//
// See: https://github.com/medincident/medincident-backend/wiki/Service-Announcements
package announcement

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

const (
	ErrCodeAnnouncementNotFound       = "announcement_not_found"
	ErrCodeAnnouncementLoadFailed     = "announcement_load_failed"
	ErrCodeAnnouncementSaveFailed     = "announcement_save_failed"
	ErrCodeAnnouncementIDGenFailed    = "announcement_id_generation_failed"
	ErrCodeAnnouncementOrgNotFound    = "announcement_organization_not_found"
	ErrCodeAnnouncementClinicNotFound = "announcement_clinic_not_found"
	ErrCodeAnnouncementDeptNotFound   = "announcement_department_not_found"
	ErrCodeAnnouncementArchived       = "announcement_archived"
	ErrCodeAnnouncementInvalidScope   = "announcement_invalid_scope"
	ErrCodeAnnouncementTimeRange      = "announcement_invalid_time_range"
)

const scope = "services.command.announcement"

// AnnouncementService is the concrete write-side service for announcements.
type AnnouncementService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewAnnouncementService wires the service.
func NewAnnouncementService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *AnnouncementService {
	return &AnnouncementService{db: db, authz: az, logger: logger}
}

// writePolicy returns the authorization policy for mutating an
// announcement at the given scope. Callers with broader privileges
// (OrgAdmin) may always manage clinic- and dept-level announcements.
func writePolicy(orgID uuid.UUID, clinicID, deptID uuid.NullUUID) authz.Policy {
	base := authz.AnyOf(
		authz.SystemAdmin,
		authz.OrgAdminOf.Organization(orgID),
	)
	if !clinicID.Valid {
		return base
	}
	withClinic := authz.AnyOf(base, authz.ClinicHeadOf.Clinic(clinicID.UUID))
	if !deptID.Valid {
		return withClinic
	}
	return authz.AnyOf(withClinic, authz.DeptResponsibleOf.Department(deptID.UUID))
}

// loadAnnouncement reads an announcement by id; not-found is mapped to a
// public error.
func (s *AnnouncementService) loadAnnouncement(tx *gorm.DB, id uuid.UUID) (*model.Announcement, error) {
	var a model.Announcement
	if err := tx.First(&a, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, oops.In(scope).
				Code(ErrCodeAnnouncementNotFound).
				Public("Announcement not found.").
				With("announcement_id", id).
				Wrap(err)
		}
		return nil, oops.In(scope).
			Code(ErrCodeAnnouncementLoadFailed).
			With("announcement_id", id).
			Wrap(err)
	}
	return &a, nil
}
