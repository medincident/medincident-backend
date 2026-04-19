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
)

type UpdateClinicPhysicalAddressCommand struct {
	ID      uuid.UUID
	Address AddressInput
}

func (s *ClinicService) UpdatePhysicalAddress(
	ctx context.Context,
	cmd UpdateClinicPhysicalAddressCommand,
) error {
	var errs []error
	if err := validateAddressInput(cmd.Address); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
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

		newAddress := model.Address{Text: strings.TrimSpace(cmd.Address.Text)}
		if cmd.Address.Point != nil {
			newAddress.Point = null.ValueFrom(model.Point{
				Longitude: cmd.Address.Point.Longitude,
				Latitude:  cmd.Address.Point.Latitude,
			})
		}
		if clinic.PhysicalAddress.Equal(newAddress) {
			return nil
		}
		clinic.PhysicalAddress = newAddress
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}

		return projector.ClinicPhysicalAddressChanged(tx, &clinic)
	})
}
