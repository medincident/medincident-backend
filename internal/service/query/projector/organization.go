package projector

import (
	"time"

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
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.organization", ErrCodeOrganizationProjectionFailed)
	if err != nil {
		return err
	}
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
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.organization", ErrCodeOrganizationProjectionFailed)
	if err != nil {
		return err
	}
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
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.organization", ErrCodeOrganizationProjectionFailed)
	if err != nil {
		return err
	}
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

// OrganizationDeactivated sets is_active=false in projections.organizations.
//
// See: docs/services/OrgStructure.md
func OrganizationDeactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *orgv1.OrganizationDeactivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.organization", ErrCodeOrganizationProjectionFailed)
	if err != nil {
		return err
	}
	if err := tx.Exec(`
		UPDATE projections.organizations
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		ev.GetUpdatedAt().AsTime(), id,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// OrganizationActivated sets is_active=true in projections.organizations.
//
// See: docs/services/OrgStructure.md
func OrganizationActivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *orgv1.OrganizationActivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.organization", ErrCodeOrganizationProjectionFailed)
	if err != nil {
		return err
	}
	if err := tx.Exec(`
		UPDATE projections.organizations
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		ev.GetUpdatedAt().AsTime(), id,
	).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// OrganizationDeleted removes the row from projections.organizations.
// Cascade FK rules in the projection schema clean up dependent rows.
//
// See: docs/services/OrgStructure.md
func OrganizationDeleted(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	_ *orgv1.OrganizationDeleted,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.organization", ErrCodeOrganizationProjectionFailed)
	if err != nil {
		return err
	}
	if err := tx.Exec(`DELETE FROM projections.organizations WHERE id = ?`, id).Error; err != nil {
		return oops.In("projector.organization").
			Code(ErrCodeOrganizationProjectionFailed).
			With("organization_id", aggregateID).
			Wrap(err)
	}
	return nil
}

func (p *Projectors) OrganizationCreated(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrganizationCreated) error {
	return OrganizationCreated(tx, id, t, ev)
}

func (p *Projectors) OrganizationDetailsChanged(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrganizationDetailsChanged) error {
	return OrganizationDetailsChanged(tx, id, t, ev)
}

func (p *Projectors) OrganizationLegalAddressChanged(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrganizationLegalAddressChanged) error {
	return OrganizationLegalAddressChanged(tx, id, t, ev)
}

func (p *Projectors) OrganizationDeactivated(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrganizationDeactivated) error {
	return OrganizationDeactivated(tx, id, t, ev)
}

func (p *Projectors) OrganizationActivated(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrganizationActivated) error {
	return OrganizationActivated(tx, id, t, ev)
}

func (p *Projectors) OrganizationDeleted(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrganizationDeleted) error {
	return OrganizationDeleted(tx, id, t, ev)
}
