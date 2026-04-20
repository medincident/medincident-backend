package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// UpdateIncidentCategoryDetailsCommand carries the new name and
// (optional) description for an existing incident category.
type UpdateIncidentCategoryDetailsCommand struct {
	CategoryID  uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=2,max=256"`
	Description *string   `validate:"omitnil,min=8,max=2048"`
}

// UpdateIncidentCategoryDetailsResult is empty — the event is the
// meaningful result.
type UpdateIncidentCategoryDetailsResult struct{}

func (s *IncidentCategoryService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentCategoryDetailsCommand,
) (UpdateIncidentCategoryDetailsResult, error) {
	if err := validation.Struct(cmd); err != nil {
		return UpdateIncidentCategoryDetailsResult{}, err
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cat model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&cat, "id = ?", cmd.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNotFound).
					Public("Incident category not found.").
					With("incident_category_id", cmd.CategoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", cmd.CategoryID).
				Wrap(err)
		}

		if err := lockClassifierOrg(tx, cat.OrganizationID); err != nil {
			return err
		}

		newName := strings.TrimSpace(cmd.Name)
		var newDescription null.String
		if cmd.Description != nil {
			newDescription = null.StringFrom(strings.TrimSpace(*cmd.Description))
		}
		if cat.Name == newName && cat.Description == newDescription {
			return nil
		}
		cat.Name = newName
		cat.Description = newDescription

		if err := tx.Save(&cat).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNameConflict).
					Public("An active incident category with this name already exists.").
					With("incident_category_id", cat.ID).
					With("name", cat.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", cat.ID).
				Wrap(err)
		}

		return projector.CategoryUpdateDetails(tx, &cat)
	})
	return UpdateIncidentCategoryDetailsResult{}, err
}
