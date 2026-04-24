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

// UpdateOrganizationLegalAddressPayload carries the new legal address.
type UpdateOrganizationLegalAddressPayload struct {
	ID      string `validate:"required,uuid"`
	Address AddressInput
}

// UpdateOrganizationLegalAddressCommand = caller + payload.
type UpdateOrganizationLegalAddressCommand struct {
	Caller  authz.Caller
	Payload UpdateOrganizationLegalAddressPayload
}

// UpdateLegalAddress replaces the organization's legal address. Returns
// nil without writing when the new address is equal to the stored one.
func (s *OrganizationService) UpdateLegalAddress(
	ctx context.Context,
	cmd UpdateOrganizationLegalAddressCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", id).
					Errorf("organization not found")
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", id).
				Wrap(err)
		}

		newAddress := model.Address{
			Text: strings.TrimSpace(cmd.Payload.Address.Text),
		}
		if cmd.Payload.Address.Point != nil {
			newAddress.Point = null.ValueFrom(model.Point{
				Longitude: cmd.Payload.Address.Point.Longitude,
				Latitude:  cmd.Payload.Address.Point.Latitude,
			})
		}
		if org.LegalAddress.Equal(newAddress) {
			return nil
		}
		org.LegalAddress = newAddress
		if err := tx.Save(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", id).
				Wrap(err)
		}

		return projector.OrganizationLegalAddressChanged(tx, &org)
	})
}
