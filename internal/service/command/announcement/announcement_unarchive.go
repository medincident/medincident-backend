package announcement

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UnarchiveAnnouncementPayload is the validated payload.
type UnarchiveAnnouncementPayload struct {
	ID string `validate:"required,uuid"`
}

// UnarchiveAnnouncementCommand is the full command DTO.
type UnarchiveAnnouncementCommand struct {
	Caller  authz.Caller
	Payload UnarchiveAnnouncementPayload
}

// Unarchive sets is_archived=false. Idempotent.
//
// See: docs/services/Announcements.md
func (s *AnnouncementService) Unarchive(
	ctx context.Context, cmd *UnarchiveAnnouncementCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}

	id := uuid.MustParse(cmd.Payload.ID)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		a, err := s.loadAnnouncement(tx, id)
		if err != nil {
			return err
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			writePolicy(a.OrganizationID, a.ClinicID, a.DepartmentID)); err != nil {
			return err
		}
		if !a.IsArchived {
			return nil // idempotent
		}
		if err := tx.Model(a).Updates(map[string]any{
			"is_archived": false,
			"updated_at":  time.Now(),
		}).Error; err != nil {
			return oops.In(scope).
				Code(ErrCodeAnnouncementSaveFailed).
				With("announcement_id", id).
				Wrap(err)
		}
		return nil
	})
}
