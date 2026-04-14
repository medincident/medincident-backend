package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
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

// Error codes emitted by CreateIncidentCategory and shared with other
// category service methods.
const (
	ErrCodeIncidentCategoryIDGenerationFailed         = "incident_category_id_generation_failed"
	ErrCodeIncidentCategorySaveFailed                 = "incident_category_save_failed"
	ErrCodeIncidentCategoryLoadFailed                 = "incident_category_load_failed"
	ErrCodeIncidentCategoryNotFound                   = "incident_category_not_found"
	ErrCodeIncidentCategoryParentNotFound             = "incident_category_parent_not_found"
	ErrCodeIncidentCategoryParentOrganizationMismatch = "incident_category_parent_organization_mismatch"
	ErrCodeIncidentCategoryNameConflict               = "incident_category_name_conflict"
	ErrCodeIncidentCategoryEventBuildFailed           = "incident_category_event_build_failed"
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

func buildIncidentCategoryCreatedEvent(c *model.IncidentCategory) *categoryeventv1.IncidentCategoryCreated {
	ev := &categoryeventv1.IncidentCategoryCreated{
		OrganizationId: c.OrganizationID.String(),
		Name:           c.Name,
	}
	if c.ParentCategoryID.Valid {
		s := c.ParentCategoryID.UUID.String()
		ev.ParentCategoryId = &s
	}
	if c.Description.Valid {
		d := c.Description.String
		ev.Description = &d
	}
	return ev
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
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case pgErrCodeUniqueViolation:
					return oops.In("services.incident.classifier.category").
						Code(ErrCodeIncidentCategoryNameConflict).
						Public("An active incident category with this name already exists.").
						With("organization_id", cmd.OrganizationID).
						With("name", cat.Name).
						Wrap(err)
				case pgErrCodeForeignKeyViolation:
					return oops.In("services.incident.classifier.category").
						Code(ErrCodeIncidentCategoryParentNotFound).
						Public("Parent incident category or organization not found.").
						Wrap(err)
				}
			}
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategorySaveFailed).
				With("incident_category_id", id).
				Wrap(err)
		}

		event := buildIncidentCategoryCreatedEvent(&cat)
		payload, err := anypb.New(event)
		if err != nil {
			return oops.In("services.incident.classifier.category").
				Code(ErrCodeIncidentCategoryEventBuildFailed).
				With("incident_category_id", id).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(cat.UpdatedAt),
			AggregateType: AggregateTypeIncidentCategory,
			AggregateId:   cat.ID.String(),
			Payload:       payload,
		}
		if err := outbox.AppendEvent(tx, SubjectIncidentCategoryCreated, envelope, nil); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
