package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// UpdateIncidentTypeDetailsCommand carries the new name and (optional)
// description for an existing incident type.
type UpdateIncidentTypeDetailsCommand struct {
	TypeID      uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=2,max=256"`
	Description *string   `validate:"omitnil,min=8,max=2048"`
}

// UpdateIncidentTypeDetailsResult is empty.
type UpdateIncidentTypeDetailsResult struct{}

func (s *IncidentTypeService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentTypeDetailsCommand,
) (UpdateIncidentTypeDetailsResult, error) {
	if err := validation.Struct(cmd); err != nil {
		return UpdateIncidentTypeDetailsResult{}, err
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

		newName := strings.TrimSpace(cmd.Name)
		var newDescription null.String
		if cmd.Description != nil {
			newDescription = null.StringFrom(strings.TrimSpace(*cmd.Description))
		}
		if row.Name == newName && row.Description == newDescription {
			return nil
		}
		row.Name = newName
		row.Description = newDescription
		if err := tx.Save(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNameConflict).
					Public("An active incident type with this name already exists.").
					With("incident_type_id", row.ID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}

		return projector.TypeUpdateDetails(tx, &row)
	})
	return UpdateIncidentTypeDetailsResult{}, err
}
