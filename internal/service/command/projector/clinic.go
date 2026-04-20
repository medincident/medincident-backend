package projector

import (
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
)

// ClinicCreated writes the projections.clinics row and the matching
// projections.clinic_counters row (initialised to zero) for a freshly
// persisted Clinic, and bumps projections.organization_counters
// .clinics_total for the parent organization. Called by command service
// ClinicService.Create inside the same transaction as the
// domain.clinics INSERT.
func ClinicCreated(tx *gorm.DB, c *model.Clinic) error {
	if err := tx.Exec(`
		INSERT INTO projections.clinics
		    (id, organization_id, name, description,
		     physical_address_text, physical_address_longitude, physical_address_latitude,
		     created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.OrganizationID, c.Name, c.Description,
		c.PhysicalAddress.Text,
		nullableLongitude(c.PhysicalAddress.Point),
		nullableLatitude(c.PhysicalAddress.Point),
		c.CreatedAt, c.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", c.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		INSERT INTO projections.clinic_counters
		    (clinic_id, organization_id, employees_total, departments_total, updated_at)
		VALUES (?, ?, 0, 0, ?)`,
		c.ID, c.OrganizationID, c.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", c.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.organization_counters
		   SET clinics_total = clinics_total + 1, updated_at = ?
		 WHERE organization_id = ?`,
		c.UpdatedAt, c.OrganizationID,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", c.ID).
			Wrap(err)
	}
	return nil
}

// ClinicDetailsChanged updates name + description + updated_at on
// projections.clinics and mirrors the new clinic_name into
// projections.employee_cards for every card under this clinic.
func ClinicDetailsChanged(tx *gorm.DB, c *model.Clinic) error {
	if err := tx.Exec(`
		UPDATE projections.clinics
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		c.Name, c.Description, c.UpdatedAt, c.ID,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", c.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET clinic_name = ?, updated_at = ?
		 WHERE clinic_id = ?`,
		c.Name, c.UpdatedAt, c.ID,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", c.ID).
			Wrap(err)
	}
	return nil
}

// ClinicPhysicalAddressChanged updates the physical address columns on
// projections.clinics.
func ClinicPhysicalAddressChanged(tx *gorm.DB, c *model.Clinic) error {
	if err := tx.Exec(`
		UPDATE projections.clinics
		   SET physical_address_text = ?, physical_address_longitude = ?, physical_address_latitude = ?, updated_at = ?
		 WHERE id = ?`,
		c.PhysicalAddress.Text,
		nullableLongitude(c.PhysicalAddress.Point),
		nullableLatitude(c.PhysicalAddress.Point),
		c.UpdatedAt, c.ID,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", c.ID).
			Wrap(err)
	}
	return nil
}
