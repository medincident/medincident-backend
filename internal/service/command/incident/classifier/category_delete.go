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

// DeleteIncidentCategoryPayload identifies the incident category to
// delete along with its entire subtree.
type DeleteIncidentCategoryPayload struct {
	CategoryID string `validate:"required,uuid"`
}

// DeleteIncidentCategoryCommand = caller + payload.
type DeleteIncidentCategoryCommand struct {
	Caller  authz.Caller
	Payload DeleteIncidentCategoryPayload
}

// DeleteIncidentCategoryResult is empty — events carry the real
// outcome.
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

// Delete permanently removes a category and its entire subtree.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentCategoryService) Delete(
	ctx context.Context,
	cmd DeleteIncidentCategoryCommand,
) (DeleteIncidentCategoryResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeleteIncidentCategoryResult{}, err
	}
	categoryID := uuid.MustParse(cmd.Payload.CategoryID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Category(categoryID)); err != nil {
		return DeleteIncidentCategoryResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		// Serialise classifier mutations for this org so a concurrent
		// Move can't reparent a node out from under our subtree walk.
		if err := lockClassifierOrg(tx, root.OrganizationID); err != nil {
			return err
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

		now := time.Now().UTC()

		// Append outbox events children-first (types first, then
		// categories in depth-DESC order) before the actual SQL DELETEs
		// so downstream projections never see a parent removed before
		// its descendants.
		for _, id := range typeIDs {
			msg := &classifierv1.IncidentTypeDeleted{TypeId: id.String(), DeletedAt: timestamppb.New(now)}
			payload, err := anypb.New(msg)
			if err != nil {
				return oops.In("services.incident.classifier.type").Code(ErrCodeIncidentTypeSaveFailed).Wrap(err)
			}
			if err := outbox.Append(tx, "medincident.event.incident_type.v1.deleted",
				&eventv1.Envelope{OccurredAt: timestamppb.New(now), AggregateType: "incident_type", AggregateId: id.String(), Payload: payload}); err != nil {
				return err
			}
		}
		for _, id := range categoryIDs {
			msg := &classifierv1.IncidentCategoryDeleted{CategoryId: id.String(), DeletedAt: timestamppb.New(now)}
			payload, err := anypb.New(msg)
			if err != nil {
				return oops.In("services.incident.classifier.category").Code(ErrCodeIncidentCategorySaveFailed).Wrap(err)
			}
			if err := outbox.Append(tx, "medincident.event.incident_category.v1.deleted",
				&eventv1.Envelope{OccurredAt: timestamppb.New(now), AggregateType: "incident_category", AggregateId: id.String(), Payload: payload}); err != nil {
				return err
			}
		}

		// Delete types first — they FK-reference the categories and
		// the parent_category_id FK is now ON DELETE RESTRICT.
		if len(typeIDs) > 0 {
			if err := tx.Delete(&model.IncidentType{}, "id IN ?", typeIDs).Error; err != nil {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeSaveFailed).
					With("root_category_id", root.ID).
					Wrap(err)
			}
		}

		// Delete categories children-first (depth DESC) so no row is
		// removed while its descendants still exist.
		for _, id := range categoryIDs {
			if err := tx.Delete(&model.IncidentCategory{}, "id = ?", id).Error; err != nil {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategorySaveFailed).
					With("incident_category_id", id).
					Wrap(err)
			}
		}
		return nil
	})
	return DeleteIncidentCategoryResult{}, err
}
