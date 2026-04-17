package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/outbox"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	categoryeventv1 "github.com/medincident/medincident-command-service/pkg/event/incident/category/v1"
	typeeventv1 "github.com/medincident/medincident-command-service/pkg/event/incident/type/v1"
)

type DeactivateIncidentCategoryCommand struct {
	CategoryID uuid.UUID
}

type DeactivateIncidentCategoryResult struct{}

// Deactivate performs a cascading deactivation: the target category,
// every descendant category, and every type whose category_id is in
// the resulting subtree. One outbox event is written per row actually
// toggled, all inside a single transaction.
func (s *IncidentCategoryService) Deactivate(
	ctx context.Context,
	cmd DeactivateIncidentCategoryCommand,
) (DeactivateIncidentCategoryResult, error) {
	if err := requireCategoryID(cmd.CategoryID); err != nil {
		return DeactivateIncidentCategoryResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var root model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&root, "id = ?", cmd.CategoryID).Error; err != nil {
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

		if err := lockClassifierOrg(tx, root.OrganizationID); err != nil {
			return err
		}

		deactivatedCategoryIDs, err := deactivateCategorySubtree(tx, root.ID)
		if err != nil {
			return err
		}
		deactivatedTypeIDs, err := deactivateTypesInSubtree(tx, root.ID)
		if err != nil {
			return err
		}

		for _, id := range deactivatedCategoryIDs {
			if err := appendCategoryDeactivatedEvent(tx, id, now); err != nil {
				return err
			}
		}
		for _, id := range deactivatedTypeIDs {
			if err := appendTypeDeactivatedEvent(tx, id, now); err != nil {
				return err
			}
		}
		return nil
	})
	return DeactivateIncidentCategoryResult{}, err
}

func deactivateCategorySubtree(tx *gorm.DB, root uuid.UUID) ([]uuid.UUID, error) {
	const q = `
		WITH RECURSIVE subtree AS (
			SELECT id FROM domain.incident_categories WHERE id = ?
			UNION ALL
			SELECT c.id
			FROM domain.incident_categories c
			JOIN subtree s ON c.parent_category_id = s.id
		)
		UPDATE domain.incident_categories
		SET is_active = FALSE, updated_at = now()
		WHERE id IN (SELECT id FROM subtree) AND is_active
		RETURNING id;
	`
	var ids []uuid.UUID
	if err := tx.Raw(q, root).Scan(&ids).Error; err != nil {
		return nil, oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategorySaveFailed).
			With("incident_category_id", root).
			Wrap(err)
	}
	return ids, nil
}

func deactivateTypesInSubtree(tx *gorm.DB, root uuid.UUID) ([]uuid.UUID, error) {
	const q = `
		WITH RECURSIVE subtree AS (
			SELECT id FROM domain.incident_categories WHERE id = ?
			UNION ALL
			SELECT c.id
			FROM domain.incident_categories c
			JOIN subtree s ON c.parent_category_id = s.id
		)
		UPDATE domain.incident_types
		SET is_active = FALSE, updated_at = now()
		WHERE category_id IN (SELECT id FROM subtree) AND is_active
		RETURNING id;
	`
	var ids []uuid.UUID
	if err := tx.Raw(q, root).Scan(&ids).Error; err != nil {
		return nil, oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeSaveFailed).
			With("root_category_id", root).
			Wrap(err)
	}
	return ids, nil
}

func appendCategoryDeactivatedEvent(tx *gorm.DB, categoryID uuid.UUID, now time.Time) error {
	if err := projector.CategoryDeactivate(tx, categoryID, now); err != nil {
		return err
	}
	return outbox.Publish(tx, SubjectIncidentCategoryDeactivated, AggregateTypeIncidentCategory, categoryID.String(), now, &categoryeventv1.IncidentCategoryDeactivated{})
}

func appendTypeDeactivatedEvent(tx *gorm.DB, typeID uuid.UUID, now time.Time) error {
	if err := projector.TypeDeactivate(tx, typeID, now); err != nil {
		return err
	}
	return outbox.Publish(tx, SubjectIncidentTypeDeactivated, AggregateTypeIncidentType, typeID.String(), now, &typeeventv1.IncidentTypeDeactivated{})
}
