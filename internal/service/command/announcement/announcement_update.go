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
// See: https://github.com/medincident/medincident-backend/wiki/Service-Announcements#updateannouncement
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

		var startsAt null.Time
		if cmd.Payload.StartsAt != nil {
			startsAt = null.TimeFrom(*cmd.Payload.StartsAt)
		}
		var endsAt null.Time
		if cmd.Payload.EndsAt != nil {
			endsAt = null.TimeFrom(*cmd.Payload.EndsAt)
		}

		if err := tx.Model(a).Updates(map[string]any{
			"title":      strings.TrimSpace(cmd.Payload.Title),
			"content":    strings.TrimSpace(cmd.Payload.Content),
			"starts_at":  startsAt,
			"ends_at":    endsAt,
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
