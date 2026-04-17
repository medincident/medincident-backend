package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

// Error codes emitted by CreateIncidentCategory and shared with other
// category service methods.
const (
	ErrCodeIncidentCategoryIDGenerationFailed         = "incident_category_id_generation_failed"
	ErrCodeIncidentCategoryIDEmpty                    = "incident_category_id_empty"
	ErrCodeIncidentCategorySaveFailed                 = "incident_category_save_failed"
	ErrCodeIncidentCategoryLoadFailed                 = "incident_category_load_failed"
	ErrCodeIncidentCategoryNotFound                   = "incident_category_not_found"
	ErrCodeIncidentCategoryParentNotFound             = "incident_category_parent_not_found"
	ErrCodeIncidentCategoryParentOrganizationMismatch = "incident_category_parent_organization_mismatch"
	ErrCodeIncidentCategoryParentInactive             = "incident_category_parent_inactive"
	ErrCodeIncidentCategoryNameConflict               = "incident_category_name_conflict"
	ErrCodeIncidentCategoryMaxDepthExceeded           = "incident_category_max_depth_exceeded"
)

// CreateIncidentCategoryCommand is the input of IncidentCategoryService.Create.
type CreateIncidentCategoryCommand struct {
	OrganizationID   uuid.UUID
	ParentCategoryID *uuid.UUID
	Name             string
	Description      *string
}

// CreateIncidentCategoryResult is the output of IncidentCategoryService.Create.
type CreateIncidentCategoryResult struct {
	ID uuid.UUID
}

// categoryDepth computes the 1-based depth of a category (root = 1)
// via a recursive CTE.
func categoryDepth(tx *gorm.DB, categoryID uuid.UUID) (int, error) {
	var depth int
	const q = `
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_category_id, 1 AS depth
			FROM domain.incident_categories
			WHERE id = ?
			UNION ALL
			SELECT c.id, c.parent_category_id, a.depth + 1
			FROM domain.incident_categories c
			JOIN ancestors a ON a.parent_category_id = c.id
		)
		SELECT COALESCE(MAX(depth), 0) FROM ancestors;
	`
	if err := tx.Raw(q, categoryID).Scan(&depth).Error; err != nil {
		return 0, err
	}
	return depth, nil
}

// Create persists a new IncidentCategory and appends an
// IncidentCategoryCreated event in the same transaction.
func (s *IncidentCategoryService) Create(
	ctx context.Context,
	cmd CreateIncidentCategoryCommand,
) (CreateIncidentCategoryResult, error) {
	var errs []error
	if err := validateIncidentCategoryName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentCategoryDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return CreateIncidentCategoryResult{}, errors.Join(errs...)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateIncidentCategoryResult{}, oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryIDGenerationFailed).
			Public("Failed to create incident category.").
			Wrap(err)
	}

	cat := model.IncidentCategory{
		ID:             id,
		OrganizationID: cmd.OrganizationID,
		Name:           strings.TrimSpace(cmd.Name),
		Description:    null.StringFromPtr(trimmedStringPtr(cmd.Description)),
		IsActive:       true,
	}
	if cmd.ParentCategoryID != nil {
		cat.ParentCategoryID = uuid.NullUUID{UUID: *cmd.ParentCategoryID, Valid: true}
	}

	var result CreateIncidentCategoryResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Serialise all classifier tree mutations for this organisation so
		// concurrent Create/Move/Delete in the same org can't race the
		// depth / cycle / uniqueness checks.
		if err := lockClassifierOrg(tx, cmd.OrganizationID); err != nil {
			return err
		}

		if cmd.ParentCategoryID != nil {
			var parent model.IncidentCategory
			if err := tx.First(&parent, "id = ?", *cmd.ParentCategoryID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return oops.In("services.incident.classifier.category").
						Code(ErrCodeIncidentCategoryParentNotFound).
						Public("Parent incident category not found.").
						With("parent_category_id", *cmd.ParentCategoryID).
						Wrap(err)
				}
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryLoadFailed).
					With("parent_category_id", *cmd.ParentCategoryID).
					Wrap(err)
			}
			if parent.OrganizationID != cmd.OrganizationID {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryParentOrganizationMismatch).
					Public("Parent incident category belongs to a different organization.").
					With("parent_category_id", parent.ID).
					With("parent_organization_id", parent.OrganizationID).
					With("organization_id", cmd.OrganizationID).
					Errorf("organization mismatch")
			}
			// Creating under an inactive parent is forbidden: the child
			// would be born-inactive-by-ancestor and Reactivate would
			// refuse to promote it (see category_reactivate.go's
			// firstInactiveAncestor check), leaving the row permanently
			// stuck until the whole ancestor chain is reactivated.
			if !parent.IsActive {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryParentInactive).
					Public("Parent incident category is inactive.").
					With("parent_category_id", parent.ID).
					Errorf("parent inactive")
			}
			depth, err := categoryDepth(tx, parent.ID)
			if err != nil {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryLoadFailed).
					With("parent_category_id", parent.ID).
					Wrap(err)
			}
			if depth+1 > incidentClassifierMaxDepth {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryMaxDepthExceeded).
					Public("Incident category tree would exceed maximum depth.").
					With("max_depth", incidentClassifierMaxDepth).
					With("parent_depth", depth).
					Errorf("max depth exceeded")
			}
		}

		if err := tx.Create(&cat).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryNameConflict).
					Public("An active incident category with this name already exists.").
					With("organization_id", cmd.OrganizationID).
					With("name", cat.Name).
					Wrap(err)
			} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return oops.In("services.incident.classifier.category").
					Code(ErrCodeIncidentCategoryParentNotFound).
					Public("Parent incident category or organization not found.").
					Wrap(err)
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", id).
				Wrap(err)
		}

		if err := projector.CategoryCreated(tx, &cat); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
