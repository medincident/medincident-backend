package classifier

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	categoryeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/category/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

const (
	ErrCodeIncidentCategoryReactivateInactiveAncestor = "incident_category_reactivate_inactive_ancestor"
	ErrCodeIncidentCategoryReactivateNameConflict     = "incident_category_reactivate_name_conflict"
)

type ReactivateIncidentCategoryCommand struct {
	CategoryID uuid.UUID
}

type ReactivateIncidentCategoryResult struct{}

// firstInactiveAncestor returns the id of the nearest inactive strict
// ancestor of the given category, or uuid.Nil if all strict ancestors
// are active (or the category is a root).
func firstInactiveAncestor(tx *gorm.DB, categoryID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	const q = `
		WITH RECURSIVE chain AS (
			SELECT id, parent_category_id, is_active, 0 AS depth
			FROM domain.incident_categories
			WHERE id = ?
			UNION ALL
			SELECT c.id, c.parent_category_id, c.is_active, ch.depth + 1
			FROM domain.incident_categories c
			JOIN chain ch ON ch.parent_category_id = c.id
		)
		SELECT id FROM chain
		WHERE depth > 0 AND is_active = FALSE
		ORDER BY depth ASC
		LIMIT 1;
	`
	row := tx.Raw(q, categoryID).Row()
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (s *IncidentCategoryService) Reactivate(
	ctx context.Context,
	cmd ReactivateIncidentCategoryCommand,
) (ReactivateIncidentCategoryResult, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cat model.IncidentCategory
		if err := tx.First(&cat, "id = ?", cmd.CategoryID).Error; err != nil {
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
		if cat.IsActive {
			return nil
		}

		ancestor, err := firstInactiveAncestor(tx, cat.ID)
		if err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryLoadFailed).
				With("incident_category_id", cat.ID).
				Wrap(err)
		}
		if ancestor != uuid.Nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryReactivateInactiveAncestor).
				Public("Cannot reactivate: an ancestor incident category is inactive.").
				With("incident_category_id", cat.ID).
				With("inactive_ancestor_id", ancestor).
				Errorf("inactive ancestor blocks reactivation")
		}

		if err := tx.Model(&model.IncidentCategory{}).
			Where("id = ?", cat.ID).
			Updates(map[string]any{"is_active": true}).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryReactivateNameConflict).
					Public("Cannot reactivate: another active incident category uses this name.").
					With("incident_category_id", cat.ID).
					With("name", cat.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", cat.ID).
				Wrap(err)
		}

		payload, err := anypb.New(&categoryeventv1.IncidentCategoryReactivated{})
		if err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryEventBuildFailed).
				With("incident_category_id", cat.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.Now(),
			AggregateType: AggregateTypeIncidentCategory,
			AggregateId:   cat.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectIncidentCategoryReactivated, envelope, nil)
	})
	return ReactivateIncidentCategoryResult{}, err
}
