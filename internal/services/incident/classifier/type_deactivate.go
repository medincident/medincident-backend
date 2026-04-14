package classifier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

type DeactivateIncidentTypeCommand struct {
	TypeID uuid.UUID
}

type DeactivateIncidentTypeResult struct{}

func (s *IncidentTypeService) Deactivate(
	ctx context.Context,
	cmd DeactivateIncidentTypeCommand,
) (DeactivateIncidentTypeResult, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.IncidentType
		if err := tx.First(&row, "id = ?", cmd.TypeID).Error; err != nil {
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
		if !row.IsActive {
			return nil
		}
		if err := tx.Model(&model.IncidentType{}).
			Where("id = ?", row.ID).
			Updates(map[string]any{"is_active": false}).Error; err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		return appendTypeDeactivatedEvent(tx, row.ID)
	})
	return DeactivateIncidentTypeResult{}, err
}
