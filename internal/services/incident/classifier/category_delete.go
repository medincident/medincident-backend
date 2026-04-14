package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

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

// collectCategorySubtreeIDs returns the full list of category ids in
// the subtree rooted at `root`, inclusive.
func collectCategorySubtreeIDs(tx *gorm.DB, root uuid.UUID) ([]uuid.UUID, error) {
	const q = `
		WITH RECURSIVE subtree AS (
			SELECT id FROM domain.incident_categories WHERE id = ?
			UNION ALL
			SELECT c.id
			FROM domain.incident_categories c
			JOIN subtree s ON c.parent_category_id = s.id
		)
		SELECT id FROM subtree;
	`
	var ids []uuid.UUID
	if err := tx.Raw(q, root).Scan(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// collectTypeIDsInCategories returns all type ids whose category_id is
// in the supplied slice.
func collectTypeIDsInCategories(tx *gorm.DB, categoryIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}
	var ids []uuid.UUID
	if err := tx.Model(&model.IncidentType{}).
		Where("category_id IN ?", categoryIDs).
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
		var root model.IncidentCategory
		if err := tx.First(&root, "id = ?", cmd.CategoryID).Error; err != nil {
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

		categoryIDs, err := collectCategorySubtreeIDs(tx, root.ID)
		if err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", root.ID).
				Wrap(err)
		}
		typeIDs, err := collectTypeIDsInCategories(tx, categoryIDs)
		if err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("root_category_id", root.ID).
				Wrap(err)
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
				OccurredAt:    timestamppb.Now(),
				AggregateType: AggregateTypeIncidentCategory,
				AggregateId:   id.String(),
				Payload:       payload,
			}
			if err := outbox.AppendEvent(tx, SubjectIncidentCategoryDeleted, envelope, nil); err != nil {
				return err
			}
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
				OccurredAt:    timestamppb.Now(),
				AggregateType: AggregateTypeIncidentType,
				AggregateId:   id.String(),
				Payload:       payload,
			}
			if err := outbox.AppendEvent(tx, SubjectIncidentTypeDeleted, envelope, nil); err != nil {
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
