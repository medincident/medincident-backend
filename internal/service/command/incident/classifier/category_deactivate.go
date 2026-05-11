package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeactivateIncidentCategoryPayload identifies the incident category
// to deactivate (cascades to descendants and types).
type DeactivateIncidentCategoryPayload struct {
	CategoryID string `validate:"required,uuid"`
}

// DeactivateIncidentCategoryCommand = caller + payload.
type DeactivateIncidentCategoryCommand struct {
	Caller  authz.Caller
	Payload DeactivateIncidentCategoryPayload
}

// DeactivateIncidentCategoryResult is empty — events carry the real
// outcome.
type DeactivateIncidentCategoryResult struct{}

// Deactivate performs a cascading deactivation: the target category,
// every descendant category, and every type whose category_id is in
// the resulting subtree. One outbox event is written per row actually
// toggled, all inside a single transaction.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentCategoryService) Deactivate(
	ctx context.Context,
	cmd DeactivateIncidentCategoryCommand,
) (DeactivateIncidentCategoryResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeactivateIncidentCategoryResult{}, err
	}
	categoryID := uuid.MustParse(cmd.Payload.CategoryID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Category(categoryID)); err != nil {
		return DeactivateIncidentCategoryResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var root model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&root, "id = ?", categoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNotFound).
					Public("Incident category not found.").
					With("incident_category_id", categoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", categoryID).
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
	msg := &classifierv1.IncidentCategoryDeactivated{CategoryId: categoryID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.incident.classifier.category").Code(ErrCodeIncidentCategorySaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{OccurredAt: timestamppb.New(now), AggregateType: "incident_category", AggregateId: categoryID.String(), Payload: payload}
	return outbox.Append(tx, "medincident.event.incident_category.v1.deactivated", env)
}

func appendTypeDeactivatedEvent(tx *gorm.DB, typeID uuid.UUID, now time.Time) error {
	msg := &classifierv1.IncidentTypeDeactivated{TypeId: typeID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.incident.classifier.type").Code(ErrCodeIncidentTypeSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{OccurredAt: timestamppb.New(now), AggregateType: "incident_type", AggregateId: typeID.String(), Payload: payload}
	return outbox.Append(tx, "medincident.event.incident_type.v1.deactivated", env)
}
