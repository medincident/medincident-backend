package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// DeactivateRequestTypePayload identifies the request type to deactivate.
type DeactivateRequestTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// DeactivateRequestTypeCommand = caller + payload.
type DeactivateRequestTypeCommand struct {
	Caller  authz.Caller
	Payload DeactivateRequestTypePayload
}

// DeactivateRequestTypeResult is empty.
type DeactivateRequestTypeResult struct{}

// Deactivate marks a request type as inactive. Idempotent.
// See: docs/services/request/Classifier.md
func (s *RequestTypeService) Deactivate(
	ctx context.Context,
	cmd DeactivateRequestTypeCommand,
) (DeactivateRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeactivateRequestTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.RequestType(typeID)); err != nil {
		return DeactivateRequestTypeResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.RequestType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNotFound).
					Public("Request type not found.").
					With("request_type_id", typeID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeRequestTypeLoadFailed).
				With("request_type_id", typeID).Wrap(err)
		}
		if !row.IsActive {
			return nil
		}
		now := time.Now()
		row.IsActive = false
		row.UpdatedAt = now
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", row.ID).Wrap(err)
		}
		return projector.RequestTypeDeactivated(tx, row.ID, now)
	})
	return DeactivateRequestTypeResult{}, err
}
