package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
)

type DeleteIncidentTypeCommand struct {
	TypeID uuid.UUID
}

type DeleteIncidentTypeResult struct{}

func (s *IncidentTypeService) Delete(
	ctx context.Context,
	cmd DeleteIncidentTypeCommand,
) (DeleteIncidentTypeResult, error) {
	if err := requireTypeID(cmd.TypeID); err != nil {
		return DeleteIncidentTypeResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.IncidentType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", cmd.TypeID).Error; err != nil {
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

		if err := lockClassifierOrg(tx, row.OrganizationID); err != nil {
			return err
		}

		if err := projector.TypeDeleted(tx, row.ID); err != nil {
			return err
		}
		if err := tx.Delete(&model.IncidentType{}, "id = ?", row.ID).Error; err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		return nil
	})
	return DeleteIncidentTypeResult{}, err
}
