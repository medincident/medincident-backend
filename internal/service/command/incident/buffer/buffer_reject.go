package buffer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

type RejectPayload struct {
	BufferID string `validate:"required,uuid"`
}

type RejectCommand struct {
	Caller  authz.Caller
	Payload RejectPayload
}

// Reject closes a buffer entry without creating an incident. Allowed
// for OrgDispatcher / OrgAdmin / SystemAdmin of the buffer's org while
// the entry is pending.
func (s *BufferService) Reject(ctx context.Context, cmd *RejectCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.BufferID)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Load buffer to read organization_id, then authorize BEFORE the
		// state check so unauthorized callers cannot distinguish
		// pending vs non-pending vs not-found from the error response —
		// they all uniformly return permission_denied.
		b, err := s.loadBuffer(tx, id)
		if err != nil {
			return err
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			authz.AnyOf(
				authz.SystemAdmin,
				authz.OrgAdminOf.Organization(b.OrganizationID),
				authz.OrgDispatcherOf.Organization(b.OrganizationID),
			),
		); err != nil {
			return err
		}
		if b.Status != model.BufferStatusPending {
			return oops.In(scope).Code(ErrCodeBufferNotPending).
				Public("Buffer entry is not pending.").
				With("status", b.Status).Errorf("not pending")
		}
		b.Status = model.BufferStatusRejected
		b.UpdatedAt = now
		if err := tx.Save(b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		return projector.BufferUpdated(tx, b)
	})
}
