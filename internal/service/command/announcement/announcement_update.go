package announcement

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UpdateAnnouncementPayload is the validated payload for updating mutable fields.
type UpdateAnnouncementPayload struct {
	ID       string     `validate:"required,uuid"`
	Title    string     `validate:"required,no_extra_ws,min=3,max=200"`
	Content  string     `validate:"required,no_extra_ws,min=10,max=5000"`
	StartsAt *time.Time // nil = clear
	EndsAt   *time.Time // nil = clear
}

// UpdateAnnouncementCommand is the full command DTO.
type UpdateAnnouncementCommand struct {
	Caller  authz.Caller
	Payload UpdateAnnouncementPayload
}

// Update mutates title, content, starts_at, ends_at.
//
// See: docs/services/Announcements.md
func (s *AnnouncementService) Update(
	ctx context.Context, cmd *UpdateAnnouncementCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}

	if cmd.Payload.StartsAt != nil && cmd.Payload.EndsAt != nil &&
		!cmd.Payload.EndsAt.After(*cmd.Payload.StartsAt) {
		return oops.In(scope).
			Code(ErrCodeAnnouncementTimeRange).
			Public("ends_at must be after starts_at.").
			Errorf("invalid time range")
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

		a.Title = strings.TrimSpace(cmd.Payload.Title)
		a.Content = strings.TrimSpace(cmd.Payload.Content)
		a.StartsAt = null.Time{}
		a.EndsAt = null.Time{}
		if cmd.Payload.StartsAt != nil {
			a.StartsAt = null.TimeFrom(*cmd.Payload.StartsAt)
		}
		if cmd.Payload.EndsAt != nil {
			a.EndsAt = null.TimeFrom(*cmd.Payload.EndsAt)
		}
		a.UpdatedAt = time.Now()

		if err := tx.Model(a).
			Select("title", "content", "starts_at", "ends_at", "updated_at").
			Updates(a).Error; err != nil {
			return oops.In(scope).
				Code(ErrCodeAnnouncementSaveFailed).
				With("announcement_id", id).
				Wrap(err)
		}
		return nil
	})
}
