package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	categoryeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/category/v1"
	typeeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/type/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
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
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		deactivatedCategoryIDs, err := deactivateCategorySubtree(tx, root.ID)
		if err != nil {
			return err
		}
		deactivatedTypeIDs, err := deactivateTypesInSubtree(tx, root.ID)
		if err != nil {
			return err
		}

		for _, id := range deactivatedCategoryIDs {
			if err := appendCategoryDeactivatedEvent(tx, id); err != nil {
				return err
			}
		}
		for _, id := range deactivatedTypeIDs {
			if err := appendTypeDeactivatedEvent(tx, id); err != nil {
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

func appendCategoryDeactivatedEvent(tx *gorm.DB, categoryID uuid.UUID) error {
	payload, err := anypb.New(&categoryeventv1.IncidentCategoryDeactivated{})
	if err != nil {
		return oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryEventBuildFailed).
			With("incident_category_id", categoryID).
			Wrap(err)
	}
	envelope := &envelopev1.Envelope{
		OccurredAt:    timestamppb.Now(),
		AggregateType: AggregateTypeIncidentCategory,
		AggregateId:   categoryID.String(),
		Payload:       payload,
	}
	return outbox.AppendEvent(tx, SubjectIncidentCategoryDeactivated, envelope)
}

func appendTypeDeactivatedEvent(tx *gorm.DB, typeID uuid.UUID) error {
	payload, err := anypb.New(&typeeventv1.IncidentTypeDeactivated{})
	if err != nil {
		return oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeEventBuildFailed).
			With("incident_type_id", typeID).
			Wrap(err)
	}
	envelope := &envelopev1.Envelope{
		OccurredAt:    timestamppb.Now(),
		AggregateType: AggregateTypeIncidentType,
		AggregateId:   typeID.String(),
		Payload:       payload,
	}
	return outbox.AppendEvent(tx, SubjectIncidentTypeDeactivated, envelope)
}
