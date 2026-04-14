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

	categoryeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/category/v1"
	typeeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/type/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

type DeleteIncidentCategoryCommand struct {
	CategoryID uuid.UUID
}

type DeleteIncidentCategoryResult struct{}

// lockCategorySubtreeIDs returns the full list of category ids in the
// subtree rooted at `root`, inclusive, ordered children-first
// (deepest depth first, then by id as a stable tiebreaker), and takes
// a row-level FOR UPDATE lock on every returned row. Children-first
// ordering matches the order in which Delete events must be emitted
// so downstream projections never see a parent removed before its
// descendants.
func lockCategorySubtreeIDs(tx *gorm.DB, root uuid.UUID) ([]uuid.UUID, error) {
	const q = `
		WITH RECURSIVE subtree(id, depth) AS (
			SELECT id, 0
			FROM domain.incident_categories
			WHERE id = ?
			UNION ALL
			SELECT c.id, s.depth + 1
			FROM domain.incident_categories c
			JOIN subtree s ON c.parent_category_id = s.id
		)
		SELECT t.id
		FROM subtree s
		JOIN domain.incident_categories t ON t.id = s.id
		ORDER BY s.depth DESC, t.id
		FOR UPDATE OF t;
	`
	var ids []uuid.UUID
	if err := tx.Raw(q, root).Scan(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// lockTypeIDsInCategories returns all type ids whose category_id is
// in the supplied slice, ordered by id, and takes a row-level FOR
// UPDATE lock on every returned row.
func lockTypeIDsInCategories(tx *gorm.DB, categoryIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}
	var ids []uuid.UUID
	if err := tx.Model(&model.IncidentType{}).
		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
		Where("category_id IN ?", categoryIDs).
		Order("id").
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *IncidentCategoryService) Delete(
	ctx context.Context,
	cmd DeleteIncidentCategoryCommand,
) (DeleteIncidentCategoryResult, error) {
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

		categoryIDs, err := lockCategorySubtreeIDs(tx, root.ID)
		if err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", root.ID).
				Wrap(err)
		}
		typeIDs, err := lockTypeIDsInCategories(tx, categoryIDs)
		if err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("root_category_id", root.ID).
				Wrap(err)
		}

		for _, id := range typeIDs {
			payload, err := anypb.New(&typeeventv1.IncidentTypeDeleted{})
			if err != nil {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeEventBuildFailed).
					With("incident_type_id", id).
					Wrap(err)
			}
			envelope := &envelopev1.Envelope{
				OccurredAt:    timestamppb.New(now),
				AggregateType: AggregateTypeIncidentType,
				AggregateId:   id.String(),
				Payload:       payload,
			}
			if err := outbox.AppendEvent(tx, SubjectIncidentTypeDeleted, envelope, nil); err != nil {
				return err
			}
		}
		for _, id := range categoryIDs {
			payload, err := anypb.New(&categoryeventv1.IncidentCategoryDeleted{})
			if err != nil {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryEventBuildFailed).
					With("incident_category_id", id).
					Wrap(err)
			}
			envelope := &envelopev1.Envelope{
				OccurredAt:    timestamppb.New(now),
				AggregateType: AggregateTypeIncidentCategory,
				AggregateId:   id.String(),
				Payload:       payload,
			}
			if err := outbox.AppendEvent(tx, SubjectIncidentCategoryDeleted, envelope, nil); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.IncidentCategory{}, "id = ?", root.ID).Error; err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", root.ID).
				Wrap(err)
		}
		return nil
	})
	return DeleteIncidentCategoryResult{}, err
}
