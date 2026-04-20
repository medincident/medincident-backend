package orgstructure

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

// UpdateClinicDetailsCommand carries the new name and (optional)
// description for an existing clinic.
type UpdateClinicDetailsCommand struct {
	ID          uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=4,max=256"`
	Description *string   `validate:"omitnil,min=8,max=2048"`
}

// UpdateDetails changes a clinic's name and description.
func (s *ClinicService) UpdateDetails(
	ctx context.Context,
	cmd UpdateClinicDetailsCommand,
) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var clinic model.Clinic
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&clinic, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", cmd.ID).
					Errorf("clinic not found")
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Name)
		var newDesc null.String
		if cmd.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Description))
		}
		if clinic.Name == newName && clinic.Description == newDesc {
			return nil
		}
		clinic.Name = newName
		clinic.Description = newDesc
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}

		return projector.ClinicDetailsChanged(tx, &clinic)
	})
}
