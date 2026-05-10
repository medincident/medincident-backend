package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
)

// ClinicCreated writes projections.clinics + projections.clinic_counters
// and bumps projections.organization_counters.clinics_total.
//
// See: docs/services/OrgStructure.md
func ClinicCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *clinicv1.ClinicCreated,
) error {
	id := uuid.MustParse(aggregateID)
	orgID := uuid.MustParse(ev.GetOrganizationId())
	var lon, lat *float64
	if p := ev.GetPhysicalAddress().GetPoint(); p != nil {
		lo, la := p.GetLongitude(), p.GetLatitude()
		lon, lat = &lo, &la
	}
	createdAt := ev.GetCreatedAt().AsTime()

	if err := tx.Exec(`
		INSERT INTO projections.clinics
		    (id, organization_id, name, description,
		     physical_address_text, physical_address_longitude, physical_address_latitude,
		     created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING`,
		id, orgID, ev.GetName(), ev.GetDescription(),
		ev.GetPhysicalAddress().GetText(), lon, lat,
		createdAt, createdAt,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		INSERT INTO projections.clinic_counters
		    (clinic_id, organization_id, employees_total, departments_total, updated_at)
		VALUES (?, ?, 0, 0, ?)
		ON CONFLICT (clinic_id) DO NOTHING`,
		id, orgID, createdAt,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.organization_counters
		   SET clinics_total = clinics_total + 1, updated_at = ?
		 WHERE organization_id = ?`,
		createdAt, orgID,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// ClinicDetailsChanged updates name + description + mirrors name into employee_cards.
//
// See: docs/services/OrgStructure.md
func ClinicDetailsChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *clinicv1.ClinicDetailsChanged,
) error {
	id := uuid.MustParse(aggregateID)
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.clinics
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetName(), ev.GetDescription(), updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET clinic_name = ?, updated_at = ?
		 WHERE clinic_id = ?`,
		ev.GetName(), updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// ClinicPhysicalAddressChanged updates address columns.
//
// See: docs/services/OrgStructure.md
func ClinicPhysicalAddressChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *clinicv1.ClinicPhysicalAddressChanged,
) error {
	id := uuid.MustParse(aggregateID)
	var lon, lat *float64
	if p := ev.GetPhysicalAddress().GetPoint(); p != nil {
		lo, la := p.GetLongitude(), p.GetLatitude()
		lon, lat = &lo, &la
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.clinics
		   SET physical_address_text      = ?,
		       physical_address_longitude = ?,
		       physical_address_latitude  = ?,
		       updated_at                 = ?
		 WHERE id = ?`,
		ev.GetPhysicalAddress().GetText(), lon, lat, updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", aggregateID).
			Wrap(err)
	}
	return nil
}
