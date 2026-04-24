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

const (
	ErrCodeIncidentTypeIDGenerationFailed = "incident_type_id_generation_failed"
	ErrCodeIncidentTypeSaveFailed         = "incident_type_save_failed"
	ErrCodeIncidentTypeLoadFailed         = "incident_type_load_failed"
	ErrCodeIncidentTypeNotFound           = "incident_type_not_found"
	ErrCodeIncidentTypeCategoryNotFound   = "incident_type_category_not_found"
	ErrCodeIncidentTypeCategoryInactive   = "incident_type_category_inactive"
	ErrCodeIncidentTypeNameConflict       = "incident_type_name_conflict"
)

// CreateIncidentTypePayload is the validated client-facing payload.
type CreateIncidentTypePayload struct {
	CategoryID  string  `validate:"required,uuid"`
	Name        string  `validate:"required,min=2,max=256"`
	Description *string `validate:"omitnil,min=8,max=2048"`
}

// CreateIncidentTypeCommand = caller + payload.
type CreateIncidentTypeCommand struct {
	Caller  authz.Caller
	Payload CreateIncidentTypePayload
}

// CreateIncidentTypeResult is the output of IncidentTypeService.Create.
type CreateIncidentTypeResult struct {
	ID uuid.UUID
}

func (s *IncidentTypeService) Create(
	ctx context.Context,
	cmd CreateIncidentTypeCommand,
) (CreateIncidentTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateIncidentTypeResult{}, err
	}
	categoryID := uuid.MustParse(cmd.Payload.CategoryID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Category(categoryID)); err != nil {
		return CreateIncidentTypeResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateIncidentTypeResult{}, oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeIDGenerationFailed).
			Public("Failed to create incident type.").
			Wrap(err)
	}

	var result CreateIncidentTypeResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cat model.IncidentCategory
		if err := tx.First(&cat, "id = ?", categoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeCategoryNotFound).
					Public("Incident category not found.").
					With("incident_category_id", categoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_category_id", categoryID).
				Wrap(err)
		}

		// Serialise classifier mutations for this organisation.
		if err := lockClassifierOrg(tx, cat.OrganizationID); err != nil {
			return err
		}

		// Forbid creating a type under an inactive category. The type
		// would be born-inactive-by-ancestor and Reactivate would refuse
		// to promote it until the category chain is reactivated.
		if !cat.IsActive {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeCategoryInactive).
				Public("Owning incident category is inactive.").
				With("incident_category_id", cat.ID).
				Errorf("category inactive")
		}

		row := model.IncidentType{
			ID:                   id,
			OrganizationID:       cat.OrganizationID,
			CategoryID:           cat.ID,
			Name:                 strings.TrimSpace(cmd.Payload.Name),
			IsActive:             true,
			IsAllowedForPatients: false,
		}
		if cmd.Payload.Description != nil {
			row.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}

		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNameConflict).
					Public("An active incident type with this name already exists.").
					With("organization_id", cat.OrganizationID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", id).
				Wrap(err)
		}

		if err := projector.TypeCreated(tx, &row); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
