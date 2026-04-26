package buffer

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

type UpdatePayload struct {
	BufferID    string  `validate:"required,uuid"`
	CategoryID  *string `validate:"omitnil,uuid"`
	TypeID      *string `validate:"omitnil,uuid"`
	Description *string `validate:"omitnil,no_extra_ws,min=1,max=10000"`
	OccurredAt  *string
}

type UpdateCommand struct {
	Caller  authz.Caller
	Payload UpdatePayload
}

// Update lets the patient owner edit a pending buffer entry.
// All four input fields are independently mutable; passing nil keeps
// the current value (no clear-to-null semantics).
//
// See: docs/services/incident/Buffer.md
func (s *BufferService) Update(ctx context.Context, cmd UpdateCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.Authenticated); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.BufferID)
	now := time.Now()
	var occurredAt *time.Time
	if cmd.Payload.OccurredAt != nil {
		t, err := parseOccurredAt(*cmd.Payload.OccurredAt)
		if err != nil {
			return err
		}
		if err := validateOccurredAt(t, now); err != nil {
			return err
		}
		occurredAt = &t
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		b, err := s.loadBuffer(tx, id)
		if err != nil {
			return err
		}
		if b.PatientZitadelUserID != cmd.Caller.ZitadelUserID {
			return oops.In(scope).Code(ErrCodeBufferNotPatient).
				Public("Only the submitting patient may edit this entry.").
				Errorf("not owner")
		}
		if b.Status != model.BufferStatusPending {
			return oops.In(scope).Code(ErrCodeBufferNotPending).
				Public("Buffer entry is no longer editable.").
				With("status", b.Status).Errorf("not pending")
		}

		if cmd.Payload.CategoryID != nil {
			b.CategoryID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.CategoryID), Valid: true}
		}
		if cmd.Payload.TypeID != nil {
			b.TypeID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.TypeID), Valid: true}
		}
		if err := validatePatientCategoryType(tx, b.OrganizationID, b.CategoryID, b.TypeID); err != nil {
			return err
		}
		if cmd.Payload.Description != nil {
			b.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if occurredAt != nil {
			b.OccurredAt = null.TimeFrom(*occurredAt)
		}
		b.UpdatedAt = now

		if err := tx.Save(b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		return projector.BufferUpdated(tx, b)
	})
}
