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
	"github.com/medincident/medincident-backend/internal/service/command/projector"
)

const (
	ErrCodeIncidentTypeMoveOrganizationMismatch = "incident_type_move_organization_mismatch"
)

type MoveIncidentTypeCommand struct {
	TypeID        uuid.UUID
	NewCategoryID uuid.UUID
}

type MoveIncidentTypeResult struct{}

func (s *IncidentTypeService) Move(
	ctx context.Context,
	cmd MoveIncidentTypeCommand,
) (MoveIncidentTypeResult, error) {
	var errs []error
	if err := requireTypeID(cmd.TypeID); err != nil {
		errs = append(errs, err)
	}
	if cmd.NewCategoryID == uuid.Nil {
		errs = append(errs, oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentCategoryIDEmpty).
			Public("New incident category ID is required.").
			With("field", "new_category_id").
			Errorf("new category id is empty"))
	}
	if len(errs) > 0 {
		return MoveIncidentTypeResult{}, errors.Join(errs...)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var moving model.IncidentType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&moving, "id = ?", cmd.TypeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNotFound).
					Public("Incident type not found.").
					With("incident_type_id", cmd.TypeID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", cmd.TypeID).
				Wrap(err)
		}

		if err := lockClassifierOrg(tx, moving.OrganizationID); err != nil {
			return err
		}

		var newCategory model.IncidentCategory
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			First(&newCategory, "id = ?", cmd.NewCategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeCategoryNotFound).
					Public("New incident category not found.").
					With("incident_category_id", cmd.NewCategoryID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_category_id", cmd.NewCategoryID).
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
