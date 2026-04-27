package buffer

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
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

type SubmitPayload struct {
	OrganizationID string  `validate:"required,uuid"`
	CategoryID     *string `validate:"omitnil,uuid"`
	TypeID         *string `validate:"omitnil,uuid"`
	Description    *string `validate:"omitnil,no_extra_ws,min=1,max=10000"`
	OccurredAt     *string
}

type SubmitCommand struct {
	Caller  authz.Caller
	Payload SubmitPayload
}

type SubmitResult struct {
	ID uuid.UUID
}

// Submit creates a new buffer entry for the authenticated patient.
// Authorization: any authenticated caller (Authenticated policy).
//
// See: docs/services/incident/Buffer.md
func (s *BufferService) Submit(ctx context.Context, cmd SubmitCommand) (SubmitResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return SubmitResult{}, err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.Authenticated); err != nil {
		return SubmitResult{}, err
	}
	orgID := uuid.MustParse(cmd.Payload.OrganizationID)
	now := time.Now()
	var occurredTime time.Time
	var hasOccurred bool
	if cmd.Payload.OccurredAt != nil {
		t, err := time.Parse(time.RFC3339Nano, *cmd.Payload.OccurredAt)
		if err != nil {
			return SubmitResult{}, oops.In(scope).
				Code(ErrCodeBufferOccurredAtInvalid).
				Public("occurred_at is not a valid RFC3339 timestamp.").
				With("occurred_at", *cmd.Payload.OccurredAt).Wrap(err)
		}
		if err := validateOccurredAt(t, now); err != nil {
			return SubmitResult{}, err
		}
		occurredTime = t
		hasOccurred = true
	}

	var categoryID, typeID uuid.NullUUID
	if cmd.Payload.CategoryID != nil {
		categoryID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.CategoryID), Valid: true}
	}
	if cmd.Payload.TypeID != nil {
		typeID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.TypeID), Valid: true}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return SubmitResult{}, oops.In(scope).Code(ErrCodeBufferIDGenerationFailed).Wrap(err)
	}

	var result SubmitResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.First(&org, "id = ?", orgID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).Code(ErrCodeBufferOrgNotFound).
					Public("Organization not found.").
					With("organization_id", orgID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeBufferLoadFailed).Wrap(err)
		}
		if err := validatePatientCategoryType(tx, orgID, categoryID, typeID); err != nil {
			return err
		}

		b := model.PatientIncidentBuffer{
			ID:                   id,
			OrganizationID:       orgID,
			PatientZitadelUserID: cmd.Caller.ZitadelUserID,
			CategoryID:           categoryID,
			TypeID:               typeID,
			Status:               model.BufferStatusPending,
			CreatedAt:            now,
			UpdatedAt:            now,
		}
		if cmd.Payload.Description != nil {
			b.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if hasOccurred {
			b.OccurredAt = null.TimeFrom(occurredTime)
		}
		if err := tx.Create(&b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		if err := projector.BufferCreated(tx, &b); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
