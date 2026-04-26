package announcement

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// CreateAnnouncementPayload is the validated client-facing payload.
type CreateAnnouncementPayload struct {
	OrganizationID string  `validate:"required,uuid"`
	ClinicID       *string `validate:"omitnil,uuid"`
	DepartmentID   *string `validate:"omitnil,uuid"`
	Title          string  `validate:"required,no_extra_ws,min=3,max=200"`
	Content        string  `validate:"required,no_extra_ws,min=10,max=5000"`
	Priority       string  `validate:"required,oneof=normal high"`
	StartsAt       *time.Time
	EndsAt         *time.Time
}

// CreateAnnouncementCommand is the full command DTO.
type CreateAnnouncementCommand struct {
	Caller  authz.Caller
	Payload CreateAnnouncementPayload
}

// CreateAnnouncementResult carries the new announcement ID.
type CreateAnnouncementResult struct {
	ID uuid.UUID
}

// Create persists a new announcement.
//
// See: docs/services/Announcements.md
func (s *AnnouncementService) Create(
	ctx context.Context, cmd *CreateAnnouncementCommand,
) (CreateAnnouncementResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateAnnouncementResult{}, err
	}

	orgID := uuid.MustParse(cmd.Payload.OrganizationID)

	var clinicID uuid.NullUUID
	var deptID uuid.NullUUID
	if cmd.Payload.ClinicID != nil {
		clinicID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.ClinicID), Valid: true}
	}
	if cmd.Payload.DepartmentID != nil {
		if !clinicID.Valid {
			return CreateAnnouncementResult{}, oops.In(scope).
				Code(ErrCodeAnnouncementInvalidScope).
				Public("department_id requires clinic_id.").
				Errorf("dept without clinic")
		}
		deptID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.DepartmentID), Valid: true}
	}

	if cmd.Payload.StartsAt != nil && cmd.Payload.EndsAt != nil &&
		!cmd.Payload.EndsAt.After(*cmd.Payload.StartsAt) {
		return CreateAnnouncementResult{}, oops.In(scope).
			Code(ErrCodeAnnouncementTimeRange).
			Public("ends_at must be after starts_at.").
			Errorf("invalid time range")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateAnnouncementResult{}, oops.In(scope).
			Code(ErrCodeAnnouncementIDGenFailed).
			Public("Failed to create announcement.").
			Wrap(err)
	}

	var result CreateAnnouncementResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate organization exists.
		var orgCount int64
		if err := tx.Raw(`SELECT 1 FROM domain.organizations WHERE id = ? LIMIT 1`, orgID).Scan(&orgCount).Error; err != nil {
			return oops.In(scope).Code(ErrCodeAnnouncementLoadFailed).Wrap(err)
		}
		if orgCount == 0 {
			return oops.In(scope).
				Code(ErrCodeAnnouncementOrgNotFound).
				Public("Organization not found.").
				With("organization_id", orgID).
				Errorf("organization not found")
		}

		// Validate clinic belongs to org if given.
		if clinicID.Valid {
			var clinicOrgID uuid.UUID
			if err := tx.Raw(`SELECT organization_id FROM domain.clinics WHERE id = ? LIMIT 1`, clinicID.UUID).
				Row().Scan(&clinicOrgID); err != nil || clinicOrgID == uuid.Nil {
				return oops.In(scope).
					Code(ErrCodeAnnouncementClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", clinicID.UUID).
					Wrap(err)
			}
			if clinicOrgID != orgID {
				return oops.In(scope).
					Code(ErrCodeAnnouncementClinicNotFound).
					Public("Clinic does not belong to the organization.").
					With("clinic_id", clinicID.UUID).
					Errorf("clinic org mismatch")
			}
		}

		// Validate department belongs to clinic if given.
		if deptID.Valid {
			var deptClinicID uuid.UUID
			if err := tx.Raw(`SELECT clinic_id FROM domain.departments WHERE id = ? LIMIT 1`, deptID.UUID).
				Row().Scan(&deptClinicID); err != nil || deptClinicID == uuid.Nil {
				return oops.In(scope).
					Code(ErrCodeAnnouncementDeptNotFound).
					Public("Department not found.").
					With("department_id", deptID.UUID).
					Wrap(err)
			}
			if deptClinicID != clinicID.UUID {
				return oops.In(scope).
					Code(ErrCodeAnnouncementDeptNotFound).
					Public("Department does not belong to the clinic.").
					With("department_id", deptID.UUID).
					Errorf("dept clinic mismatch")
			}
		}

		// Authorization.
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, writePolicy(orgID, clinicID, deptID)); err != nil {
			return err
		}

		now := time.Now()
		priority := model.AnnouncementPriority(cmd.Payload.Priority)

		var startsAt null.Time
		if cmd.Payload.StartsAt != nil {
			startsAt = null.TimeFrom(*cmd.Payload.StartsAt)
		}
		var endsAt null.Time
		if cmd.Payload.EndsAt != nil {
			endsAt = null.TimeFrom(*cmd.Payload.EndsAt)
		}

		a := model.Announcement{
			ID:             id,
			OrganizationID: orgID,
			ClinicID:       clinicID,
			DepartmentID:   deptID,
			AuthorID:       cmd.Caller.ZitadelUserID,
			Title:          strings.TrimSpace(cmd.Payload.Title),
			Content:        strings.TrimSpace(cmd.Payload.Content),
			Priority:       priority,
			IsArchived:     false,
			StartsAt:       startsAt,
			EndsAt:         endsAt,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(&a).Error; err != nil {
			return oops.In(scope).
				Code(ErrCodeAnnouncementSaveFailed).
				With("announcement_id", id).
				Wrap(err)
		}

		// Seed the view counter row.
		if err := tx.Exec(
			`INSERT INTO projections.announcement_views (announcement_id, view_count) VALUES (?, 0)`, id,
		).Error; err != nil {
			return oops.In(scope).
				Code(ErrCodeAnnouncementSaveFailed).
				With("announcement_id", id).
				Wrap(errors.New("failed to seed view counter: " + err.Error()))
		}

		result.ID = id
		return nil
	})
	return result, err
}
