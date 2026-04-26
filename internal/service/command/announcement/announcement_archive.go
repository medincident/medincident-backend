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

// ArchiveAnnouncementPayload is the validated payload.
type ArchiveAnnouncementPayload struct {
	ID string `validate:"required,uuid"`
}

// ArchiveAnnouncementCommand is the full command DTO.
type ArchiveAnnouncementCommand struct {
	Caller  authz.Caller
	Payload ArchiveAnnouncementPayload
}

// Archive sets is_archived=true. Idempotent.
//
// See: https://github.com/medincident/medincident-backend/wiki/Service-Announcements#archiveannouncement
func (s *AnnouncementService) Archive(
	ctx context.Context, cmd *ArchiveAnnouncementCommand,
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
		if a.IsArchived {
			return nil // idempotent
		}
		if err := tx.Model(a).Updates(map[string]any{
			"is_archived": true,
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
