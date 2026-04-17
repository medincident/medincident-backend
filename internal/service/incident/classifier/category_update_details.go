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

	categoryeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/category/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
)

type UpdateIncidentCategoryDetailsCommand struct {
	CategoryID  uuid.UUID
	Name        string
	Description *string
}

type UpdateIncidentCategoryDetailsResult struct{}

func buildIncidentCategoryDetailsChangedEvent(c *model.IncidentCategory) *categoryeventv1.IncidentCategoryDetailsChanged {
	return &categoryeventv1.IncidentCategoryDetailsChanged{
		Name:        c.Name,
		Description: c.Description.Ptr(),
	}
}

func (s *IncidentCategoryService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentCategoryDetailsCommand,
) (UpdateIncidentCategoryDetailsResult, error) {
	var errs []error
	if err := requireCategoryID(cmd.CategoryID); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentCategoryName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentCategoryDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return UpdateIncidentCategoryDetailsResult{}, errors.Join(errs...)
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
		newDescription := null.StringFromPtr(trimmedStringPtr(cmd.Description))
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

		event := buildIncidentCategoryDetailsChangedEvent(&cat)
		return outbox.Publish(tx, SubjectIncidentCategoryDetailsChanged, AggregateTypeIncidentCategory, cat.ID.String(), cat.UpdatedAt, event)
	})
	return UpdateIncidentCategoryDetailsResult{}, err
}
