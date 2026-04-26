package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// CreateRequestTypePayload is the validated client-facing payload.
type CreateRequestTypePayload struct {
	OrganizationID string  `validate:"required,uuid"`
	Name           string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description    *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// CreateRequestTypeCommand = caller + payload.
type CreateRequestTypeCommand struct {
	Caller  authz.Caller
	Payload CreateRequestTypePayload
}

// CreateRequestTypeResult is the output of RequestTypeService.Create.
type CreateRequestTypeResult struct {
	ID uuid.UUID
}

// Create creates a new request type for an organization.
// See: https://github.com/medincident/medincident-backend/wiki/Service-Requests#createrequesttype
func (s *RequestTypeService) Create(
	ctx context.Context,
	cmd CreateRequestTypeCommand,
) (CreateRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateRequestTypeResult{}, err
	}
	orgID := uuid.MustParse(cmd.Payload.OrganizationID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return CreateRequestTypeResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateRequestTypeResult{}, oops.In(scope).
			Code(ErrCodeRequestTypeIDGenerationFailed).
			Public("Failed to create request type.").
			Wrap(err)
	}

	var result CreateRequestTypeResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := model.RequestType{
			ID:             id,
			OrganizationID: orgID,
			Name:           strings.TrimSpace(cmd.Payload.Name),
			IsActive:       true,
			CreatedAt:      tx.NowFunc(),
			UpdatedAt:      tx.NowFunc(),
		}
		if cmd.Payload.Description != nil {
			row.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNameConflict).
					Public("An active request type with this name already exists.").
					With("organization_id", orgID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In(scope).
				Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", id).
				Wrap(err)
		}
		if err := projector.RequestTypeCreated(tx, &row); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
