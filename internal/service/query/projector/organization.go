package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
)

// OrganizationCreated writes projections.organizations + projections.organization_counters.
//
// See: docs/services/OrgStructure.md
func OrganizationCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *orgv1.OrganizationCreated,
) error {
	id := uuid.MustParse(aggregateID)
	var lon, lat *float64
	if p := ev.GetLegalAddress().GetPoint(); p != nil {
		lo, la := p.GetLongitude(), p.GetLatitude()
		lon, lat = &lo, &la
	}
	var desc null.String
	if ev.GetDescription() != "" {
		desc = null.StringFrom(ev.GetDescription())
	}
	createdAt := ev.GetCreatedAt().AsTime()

	if err := tx.Exec(`
		INSERT INTO projections.organizations
		    (id, name, description,
		     legal_address_text, legal_address_longitude, legal_address_latitude,
		     created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING`,
		id, ev.GetName(), desc,
		ev.GetLegalAddress().GetText(), lon, lat,
		createdAt, createdAt,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		INSERT INTO projections.organization_counters
		    (organization_id, employees_total, clinics_total, departments_total, updated_at)
		VALUES (?, 0, 0, 0, ?)
		ON CONFLICT (organization_id) DO NOTHING`,
		id, createdAt,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// OrganizationDetailsChanged updates name + description + mirrors name into employee_cards.
//
// See: docs/services/OrgStructure.md
func OrganizationDetailsChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *orgv1.OrganizationDetailsChanged,
) error {
	id := uuid.MustParse(aggregateID)
	var desc null.String
	if ev.GetDescription() != "" {
		desc = null.StringFrom(ev.GetDescription())
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.organizations
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetName(), desc, updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET organization_name = ?, updated_at = ?
		 WHERE organization_id = ?`,
		ev.GetName(), updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// OrganizationLegalAddressChanged updates address columns.
//
// See: docs/services/OrgStructure.md
func OrganizationLegalAddressChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *orgv1.OrganizationLegalAddressChanged,
) error {
	id := uuid.MustParse(aggregateID)
	var lon, lat *float64
	if p := ev.GetLegalAddress().GetPoint(); p != nil {
		lo, la := p.GetLongitude(), p.GetLatitude()
		lon, lat = &lo, &la
	}
	updatedAt := ev.GetUpdatedAt().AsTime()

	if err := tx.Exec(`
		UPDATE projections.organizations
		   SET legal_address_text      = ?,
		       legal_address_longitude = ?,
		       legal_address_latitude  = ?,
		       updated_at              = ?
		 WHERE id = ?`,
		ev.GetLegalAddress().GetText(), lon, lat, updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	return nil
}
