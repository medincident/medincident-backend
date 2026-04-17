package projector

import (
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// OrganizationCreated writes the projections.organizations row and the
// matching projections.organization_counters row (initialised to zero)
// for a freshly persisted Organization. Called by command service
// OrganizationService.Create inside the same transaction as the
// domain.organizations INSERT.
func OrganizationCreated(tx *gorm.DB, org *model.Organization) error {
	if err := tx.Exec(`
		INSERT INTO projections.organizations
		    (id, name, description,
		     legal_address_text, legal_address_longitude, legal_address_latitude,
		     created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		org.ID, org.Name, org.Description,
		org.LegalAddress.Text,
		nullableLongitude(org.LegalAddress.Point),
		nullableLatitude(org.LegalAddress.Point),
		org.CreatedAt, org.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", org.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		INSERT INTO projections.organization_counters
		    (organization_id, employees_total, clinics_total, departments_total, updated_at)
		VALUES (?, 0, 0, 0, ?)`,
		org.ID, org.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", org.ID).
			Wrap(err)
	}
	return nil
}

// OrganizationDetailsChanged updates name + description + updated_at on
// projections.organizations and mirrors the new name into
// projections.employee_cards.organization_name for every card under
// this organization.
func OrganizationDetailsChanged(tx *gorm.DB, org *model.Organization) error {
	if err := tx.Exec(`
		UPDATE projections.organizations
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		org.Name, org.Description, org.UpdatedAt, org.ID,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", org.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET organization_name = ?, updated_at = ?
		 WHERE organization_id = ?`,
		org.Name, org.UpdatedAt, org.ID,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", org.ID).
			Wrap(err)
	}
	return nil
}

// OrganizationLegalAddressChanged updates the legal address columns
// on projections.organizations.
func OrganizationLegalAddressChanged(tx *gorm.DB, org *model.Organization) error {
	if err := tx.Exec(`
		UPDATE projections.organizations
		   SET legal_address_text = ?, legal_address_longitude = ?, legal_address_latitude = ?, updated_at = ?
		 WHERE id = ?`,
		org.LegalAddress.Text,
		nullableLongitude(org.LegalAddress.Point),
		nullableLatitude(org.LegalAddress.Point),
		org.UpdatedAt, org.ID,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", org.ID).
			Wrap(err)
	}
	return nil
}

// nullableLongitude returns the float pointer or nil for storage in
// the optional legal_address_longitude / physical_address_longitude
// columns.
func nullableLongitude(p *model.Point) *float64 {
	if p == nil {
		return nil
	}
	v := p.Longitude
	return &v
}

// nullableLatitude is the latitude analogue of nullableLongitude.
func nullableLatitude(p *model.Point) *float64 {
	if p == nil {
		return nil
	}
	v := p.Latitude
	return &v
}
