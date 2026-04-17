package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

const (
	ErrCodeIncidentTypeIDGenerationFailed = "incident_type_id_generation_failed"
	ErrCodeIncidentTypeIDEmpty            = "incident_type_id_empty"
	ErrCodeIncidentTypeSaveFailed         = "incident_type_save_failed"
	ErrCodeIncidentTypeLoadFailed         = "incident_type_load_failed"
	ErrCodeIncidentTypeNotFound           = "incident_type_not_found"
	ErrCodeIncidentTypeCategoryNotFound   = "incident_type_category_not_found"
	ErrCodeIncidentTypeCategoryInactive   = "incident_type_category_inactive"
	ErrCodeIncidentTypeNameConflict       = "incident_type_name_conflict"
)

type CreateIncidentTypeCommand struct {
	CategoryID  uuid.UUID
	Name        string
	Description *string
}

type CreateIncidentTypeResult struct {
	ID uuid.UUID
}

func (s *IncidentTypeService) Create(
	ctx context.Context,
	cmd CreateIncidentTypeCommand,
) (CreateIncidentTypeResult, error) {
	var errs []error
	if err := validateIncidentTypeName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentTypeDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return CreateIncidentTypeResult{}, errors.Join(errs...)
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
		if err := tx.First(&cat, "id = ?", cmd.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeCategoryNotFound).
					Public("Incident category not found.").
					With("incident_category_id", cmd.CategoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_category_id", cmd.CategoryID).
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
			ID:             id,
			OrganizationID: cat.OrganizationID,
			CategoryID:     cat.ID,
			Name:           strings.TrimSpace(cmd.Name),
			Description:    null.StringFromPtr(trimmedStringPtr(cmd.Description)),
			IsActive:       true,
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
