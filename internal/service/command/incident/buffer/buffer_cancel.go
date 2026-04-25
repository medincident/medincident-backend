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

type CancelPayload struct {
	BufferID string `validate:"required,uuid"`
}

type CancelCommand struct {
	Caller  authz.Caller
	Payload CancelPayload
}

// Cancel marks the buffer entry cancelled. Only the patient owner,
// only while pending.
func (s *BufferService) Cancel(ctx context.Context, cmd *CancelCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.Authenticated); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.BufferID)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		b, err := s.loadBuffer(tx, id)
		if err != nil {
			return err
		}
		if b.PatientZitadelUserID != cmd.Caller.ZitadelUserID {
			return oops.In(scope).Code(ErrCodeBufferNotPatient).
				Public("Only the submitting patient may cancel this entry.").
				Errorf("not owner")
		}
		if b.Status != model.BufferStatusPending {
			return oops.In(scope).Code(ErrCodeBufferNotPending).
				Public("Buffer entry is no longer cancellable.").
				With("status", b.Status).Errorf("not pending")
		}
		b.Status = model.BufferStatusCancelled
		b.UpdatedAt = now
		if err := tx.Save(b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		return projector.BufferUpdated(tx, b)
	})
}
