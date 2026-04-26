package announcement

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UpdateAnnouncementPriorityPayload is the validated payload.
type UpdateAnnouncementPriorityPayload struct {
	ID       string `validate:"required,uuid"`
	Priority string `validate:"required,oneof=normal high"`
}

// UpdateAnnouncementPriorityCommand is the full command DTO.
type UpdateAnnouncementPriorityCommand struct {
	Caller  authz.Caller
	Payload UpdateAnnouncementPriorityPayload
}

// UpdatePriority changes the priority of a non-archived announcement.
//
// See: docs/services/Announcements.md
func (s *AnnouncementService) UpdatePriority(
	ctx context.Context, cmd *UpdateAnnouncementPriorityCommand,
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
		if a.IsArchived {
			return oops.In(scope).
				Code(ErrCodeAnnouncementArchived).
				Public("Cannot change priority of an archived announcement.").
				With("announcement_id", id).
				Errorf("archived")
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			writePolicy(a.OrganizationID, a.ClinicID, a.DepartmentID)); err != nil {
			return err
		}
		if err := tx.Model(a).Updates(map[string]any{
			"priority":   model.AnnouncementPriority(cmd.Payload.Priority),
			"updated_at": time.Now(),
		}).Error; err != nil {
			return oops.In(scope).
				Code(ErrCodeAnnouncementSaveFailed).
				With("announcement_id", id).
				Wrap(err)
		}
		return nil
	})
}
