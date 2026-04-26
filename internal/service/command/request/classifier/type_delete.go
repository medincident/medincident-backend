package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// DeleteRequestTypePayload identifies the request type to delete.
type DeleteRequestTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// DeleteRequestTypeCommand = caller + payload.
type DeleteRequestTypeCommand struct {
	Caller  authz.Caller
	Payload DeleteRequestTypePayload
}

// DeleteRequestTypeResult is empty.
type DeleteRequestTypeResult struct{}

// Delete permanently removes a request type.
// See: https://github.com/medincident/medincident-backend/wiki/Service-Requests#deleterequesttype
func (s *RequestTypeService) Delete(
	ctx context.Context,
	cmd DeleteRequestTypeCommand,
) (DeleteRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeleteRequestTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.RequestType(typeID)); err != nil {
		return DeleteRequestTypeResult{}, err
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
		if err := projector.RequestTypeDeleted(tx, row.ID); err != nil {
			return err
		}
		if err := tx.Delete(&model.RequestType{}, "id = ?", row.ID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", row.ID).Wrap(err)
		}
		return nil
	})
	return DeleteRequestTypeResult{}, err
}
