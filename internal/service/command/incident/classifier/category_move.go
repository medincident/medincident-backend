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
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

const (
	ErrCodeIncidentCategoryMoveWouldCreateCycle     = "incident_category_move_would_create_cycle"
	ErrCodeIncidentCategoryMoveWouldExceedDepth     = "incident_category_move_would_exceed_depth"
	ErrCodeIncidentCategoryMoveOrganizationMismatch = "incident_category_move_organization_mismatch"
)

type MoveIncidentCategoryCommand struct {
	CategoryID          uuid.UUID
	NewParentCategoryID *uuid.UUID
}

type MoveIncidentCategoryResult struct{}

// categorySubtreeDepth returns the internal depth of the subtree rooted
// at `root` (1 = only the root itself).
func categorySubtreeDepth(tx *gorm.DB, root uuid.UUID) (int, error) {
	var depth int
	const q = `
		WITH RECURSIVE subtree AS (
			SELECT id, 1 AS depth
			FROM domain.incident_categories
			WHERE id = ?
			UNION ALL
			SELECT c.id, s.depth + 1
			FROM domain.incident_categories c
			JOIN subtree s ON c.parent_category_id = s.id
		)
		SELECT COALESCE(MAX(depth), 0) FROM subtree;
	`
	if err := tx.Raw(q, root).Scan(&depth).Error; err != nil {
		return 0, err
	}
	return depth, nil
}

// categoryIsAncestorOf returns true iff `ancestor` is `descendant` or
// one of its transitive ancestors in the classifier tree.
func categoryIsAncestorOf(tx *gorm.DB, ancestor, descendant uuid.UUID) (bool, error) {
	var found bool
	const q = `
		WITH RECURSIVE chain AS (
			SELECT id, parent_category_id
			FROM domain.incident_categories
			WHERE id = ?
			UNION ALL
			SELECT c.id, c.parent_category_id
			FROM domain.incident_categories c
			JOIN chain ch ON ch.parent_category_id = c.id
		)
		SELECT EXISTS(SELECT 1 FROM chain WHERE id = ?);
	`
	if err := tx.Raw(q, descendant, ancestor).Scan(&found).Error; err != nil {
		return false, err
	}
	return found, nil
}

func (s *IncidentCategoryService) Move(
	ctx context.Context,
	cmd MoveIncidentCategoryCommand,
) (MoveIncidentCategoryResult, error) {
	if err := requireCategoryID(cmd.CategoryID); err != nil {
		return MoveIncidentCategoryResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var moving model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&moving, "id = ?", cmd.CategoryID).Error; err != nil {
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

		// Serialise classifier mutations for this org: depth / cycle
		// checks on concurrent moves of unrelated categories in the
		// same organisation would otherwise race.
		if err := lockClassifierOrg(tx, moving.OrganizationID); err != nil {
			return err
		}

		var newParent uuid.NullUUID
		if cmd.NewParentCategoryID != nil {
			var parent model.IncidentCategory
			if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
				First(&parent, "id = ?", *cmd.NewParentCategoryID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return oops.In("services.incident.classifier.category").
						Code(ErrCodeIncidentCategoryParentNotFound).
						Public("New parent incident category not found.").
						With("new_parent_category_id", *cmd.NewParentCategoryID).
						Wrap(err)
				}
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryLoadFailed).
					With("new_parent_category_id", *cmd.NewParentCategoryID).
					Wrap(err)
			}
			if parent.OrganizationID != moving.OrganizationID {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryMoveOrganizationMismatch).
					Public("New parent incident category belongs to a different organization.").
					With("incident_category_id", moving.ID).
					With("new_parent_category_id", parent.ID).
					Errorf("organization mismatch on move")
			}
			cycle, err := categoryIsAncestorOf(tx, moving.ID, parent.ID)
			if err != nil {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryLoadFailed).
					With("incident_category_id", moving.ID).
					Wrap(err)
			}
			if cycle {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryMoveWouldCreateCycle).
					Public("Moving this category under the target would create a cycle.").
					With("incident_category_id", moving.ID).
					With("new_parent_category_id", parent.ID).
					Errorf("move would create cycle")
			}
			parentDepth, err := categoryDepth(tx, parent.ID)
			if err != nil {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryLoadFailed).
					With("new_parent_category_id", parent.ID).
					Wrap(err)
			}
			subtreeDepth, err := categorySubtreeDepth(tx, moving.ID)
			if err != nil {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryLoadFailed).
					With("incident_category_id", moving.ID).
					Wrap(err)
			}
			if parentDepth+subtreeDepth > incidentClassifierMaxDepth {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryMoveWouldExceedDepth).
					Public("Moving this category would exceed the maximum classifier depth.").
					With("max_depth", incidentClassifierMaxDepth).
					With("parent_depth", parentDepth).
					With("subtree_depth", subtreeDepth).
					Errorf("move would exceed max depth")
			}
			newParent = uuid.NullUUID{UUID: parent.ID, Valid: true}
		}

		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_categories
			SET parent_category_id = ?, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, newParent, moving.ID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", moving.ID).
				Wrap(err)
		}

		return projector.CategoryMove(tx, moving.ID, newParent, updatedAt)
	})
	return MoveIncidentCategoryResult{}, err
}
