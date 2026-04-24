package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

const (
	ErrCodeIncidentTypeMoveOrganizationMismatch = "incident_type_move_organization_mismatch"
)

// MoveIncidentTypePayload carries the identifiers needed to move an
// incident type into a different category.
type MoveIncidentTypePayload struct {
	TypeID        string `validate:"required,uuid"`
	NewCategoryID string `validate:"required,uuid"`
}

// MoveIncidentTypeCommand = caller + payload.
type MoveIncidentTypeCommand struct {
	Caller  authz.Caller
	Payload MoveIncidentTypePayload
}

// MoveIncidentTypeResult is empty.
type MoveIncidentTypeResult struct{}

func (s *IncidentTypeService) Move(
	ctx context.Context,
	cmd MoveIncidentTypeCommand,
) (MoveIncidentTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return MoveIncidentTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	newCategoryID := uuid.MustParse(cmd.Payload.NewCategoryID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return MoveIncidentTypeResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var moving model.IncidentType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&moving, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNotFound).
					Public("Incident type not found.").
					With("incident_type_id", typeID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", typeID).
				Wrap(err)
		}

		if err := lockClassifierOrg(tx, moving.OrganizationID); err != nil {
			return err
		}

		var newCategory model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			First(&newCategory, "id = ?", newCategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeCategoryNotFound).
					Public("New incident category not found.").
					With("incident_category_id", newCategoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_category_id", newCategoryID).
				Wrap(err)
		}

		if newCategory.OrganizationID != moving.OrganizationID {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeMoveOrganizationMismatch).
				Public("New incident category belongs to a different organization.").
				With("incident_type_id", moving.ID).
				With("incident_category_id", newCategory.ID).
				Errorf("organization mismatch on type move")
		}

		var updatedAt time.Time
		if err := tx.Raw(`
			UPDATE domain.incident_types
			SET category_id = ?, updated_at = now()
			WHERE id = ?
			RETURNING updated_at`, newCategory.ID, moving.ID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", moving.ID).
				Wrap(err)
		}

		return projector.TypeMove(tx, moving.ID, newCategory.ID, updatedAt)
	})
	return MoveIncidentTypeResult{}, err
}
