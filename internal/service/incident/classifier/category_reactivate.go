package classifier

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	categoryeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/category/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/pgerr"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
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
	if err := requireCategoryID(cmd.CategoryID); err != nil {
		return ReactivateIncidentCategoryResult{}, err
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

		// Use RETURNING to get the updated_at in the same roundtrip so
		// the outbox envelope's OccurredAt matches the DB clock.
		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_categories
			SET is_active = TRUE, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, cat.ID,
		).Row().Scan(&updatedAt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerr.CodeUniqueViolation {
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

		return outbox.Publish(tx, SubjectIncidentCategoryReactivated, AggregateTypeIncidentCategory, cat.ID.String(), updatedAt, &categoryeventv1.IncidentCategoryReactivated{})
	})
	return ReactivateIncidentCategoryResult{}, err
}
