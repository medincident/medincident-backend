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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// UpdateClinicPhysicalAddressPayload carries the new physical address
// for an existing clinic.
type UpdateClinicPhysicalAddressPayload struct {
	ID      string       `validate:"required,uuid"`
	Address AddressInput `validate:"required"`
}

// UpdateClinicPhysicalAddressCommand = caller + payload.
type UpdateClinicPhysicalAddressCommand struct {
	Caller  authz.Caller
	Payload UpdateClinicPhysicalAddressPayload
}

// UpdatePhysicalAddress replaces the clinic's physical address.
//
// See: docs/services/OrgStructure.md
func (s *ClinicService) UpdatePhysicalAddress(
	ctx context.Context,
	cmd UpdateClinicPhysicalAddressCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var clinic model.Clinic
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&clinic, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", id).
					Errorf("clinic not found")
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		newAddress := model.Address{Text: strings.TrimSpace(cmd.Payload.Address.Text)}
		if cmd.Payload.Address.Point != nil {
			newAddress.Point = null.ValueFrom(model.Point{
				Longitude: cmd.Payload.Address.Point.Longitude,
				Latitude:  cmd.Payload.Address.Point.Latitude,
			})
		}
		if clinic.PhysicalAddress.Equal(newAddress) {
			return nil
		}
		clinic.PhysicalAddress = newAddress
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		return projector.ClinicPhysicalAddressChanged(tx, &clinic)
	})
}
