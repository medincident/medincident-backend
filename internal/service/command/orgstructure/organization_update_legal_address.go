package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

// UpdateOrganizationLegalAddressCommand carries the new legal address.
type UpdateOrganizationLegalAddressCommand struct {
	ID      uuid.UUID
	Address AddressInput
}

// UpdateLegalAddress replaces the organization's legal address. Returns
// nil without writing when the new address is equal to the stored one.
func (s *OrganizationService) UpdateLegalAddress(
	ctx context.Context,
	cmd UpdateOrganizationLegalAddressCommand,
) error {
	var errs []error
	if err := validateAddressInput(cmd.Address); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", cmd.ID).
					Errorf("organization not found")
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", cmd.ID).
				Wrap(err)
		}

		newAddress := model.Address{
			Text: strings.TrimSpace(cmd.Address.Text),
		}
		if cmd.Address.Point != nil {
			newAddress.Point = &model.Point{
				Longitude: cmd.Address.Point.Longitude,
				Latitude:  cmd.Address.Point.Latitude,
			}
		}
		if org.LegalAddress.Equal(newAddress) {
			return nil
		}
		org.LegalAddress = newAddress
		if err := tx.Save(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", cmd.ID).
				Wrap(err)
		}

		return projector.OrganizationLegalAddressChanged(tx, &org)
	})
}
