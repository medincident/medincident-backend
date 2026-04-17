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

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
)

// UpdateOrganizationLegalAddressCommand carries the new legal address.
type UpdateOrganizationLegalAddressCommand struct {
	ID      uuid.UUID
	Address AddressInput
}

// buildOrganizationLegalAddressChangedEvent assembles the event from
// the updated model's LegalAddress. Text is always present; point is
// conditionally attached when both coordinates are Valid.
func buildOrganizationLegalAddressChangedEvent(org *model.Organization) *organizationv1.OrganizationLegalAddressChanged {
	ev := &organizationv1.OrganizationLegalAddressChanged{
		LegalAddress: &organizationv1.Address{Text: org.LegalAddress.Text},
	}
	if org.LegalAddress.Point.Longitude.Valid && org.LegalAddress.Point.Latitude.Valid {
		ev.LegalAddress.Point = &organizationv1.Point{
			Longitude: org.LegalAddress.Point.Longitude.Float64,
			Latitude:  org.LegalAddress.Point.Latitude.Float64,
		}
	}
	return ev
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
	if err := validatePointInput(cmd.Address.Point); err != nil {
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
			newAddress.Point = model.Point{
				Longitude: null.FloatFrom(cmd.Address.Point.Longitude),
				Latitude:  null.FloatFrom(cmd.Address.Point.Latitude),
			}
		}
		if org.LegalAddress == newAddress {
			return nil
		}
		org.LegalAddress = newAddress
		if err := tx.Save(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", cmd.ID).
				Wrap(err)
		}

		event := buildOrganizationLegalAddressChangedEvent(&org)
		return outbox.Publish(tx, SubjectOrganizationLegalAddressChanged, AggregateTypeOrganization, org.ID.String(), org.UpdatedAt, event)
	})
}
