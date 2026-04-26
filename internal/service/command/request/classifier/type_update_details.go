package classifier

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UpdateRequestTypeDetailsPayload carries the new name and (optional)
// description for an existing request type.
type UpdateRequestTypeDetailsPayload struct {
	TypeID      string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateRequestTypeDetailsCommand = caller + payload.
type UpdateRequestTypeDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateRequestTypeDetailsPayload
}

// UpdateRequestTypeDetailsResult is empty.
type UpdateRequestTypeDetailsResult struct{}

// UpdateDetails updates the name and description of a request type.
// See: https://github.com/medincident/medincident-backend/wiki/Service-Requests#updaterequesttypedetails
func (s *RequestTypeService) UpdateDetails(
	ctx context.Context,
	cmd UpdateRequestTypeDetailsCommand,
) (UpdateRequestTypeDetailsResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return UpdateRequestTypeDetailsResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.RequestType(typeID)); err != nil {
		return UpdateRequestTypeDetailsResult{}, err
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.RequestType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNotFound).
					Public("Request type not found.").
					With("request_type_id", typeID).
					Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeRequestTypeLoadFailed).
				With("request_type_id", typeID).Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDescription null.String
		if cmd.Payload.Description != nil {
			newDescription = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if row.Name == newName && row.Description == newDescription {
			return nil
		}

		row.Name = newName
		row.Description = newDescription
		row.UpdatedAt = time.Now()

		if err := tx.Save(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNameConflict).
					Public("An active request type with this name already exists.").
					With("organization_id", row.OrganizationID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", row.ID).Wrap(err)
		}
		return projector.RequestTypeDetailsUpdated(tx, &row)
	})
	return UpdateRequestTypeDetailsResult{}, err
}
