package projector

import (
	"time"

	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
)

// DepartmentCreated writes projections.departments + projections.department_counters
// and bumps clinic/org counters.
//
// See: docs/services/OrgStructure.md
func DepartmentCreated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *deptv1.DepartmentCreated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.department", ErrCodeDepartmentProjectionFailed)
	if err != nil {
		return err
	}
	clinicID, err := parseUUID(ev.GetClinicId(), "clinic_id", "projector.department", ErrCodeDepartmentProjectionFailed)
	if err != nil {
		return err
	}
	createdAt := ev.GetCreatedAt().AsTime()

	orgID, err := lookupOrgIDForClinic(tx, clinicID)
	if err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			With("clinic_id", ev.GetClinicId()).
			Wrap(err)
	}

	var desc null.String
	if ev.GetDescription() != "" {
		desc = null.StringFrom(ev.GetDescription())
	}

	if err := tx.Exec(`
		INSERT INTO projections.departments
		    (id, clinic_id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING`,
		id, clinicID, ev.GetName(), desc,
		createdAt, createdAt,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		INSERT INTO projections.department_counters
		    (department_id, clinic_id, organization_id, employees_total, updated_at)
		VALUES (?, ?, ?, 0, ?)
		ON CONFLICT (department_id) DO NOTHING`,
		id, clinicID, orgID, createdAt,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.clinic_counters
		   SET departments_total = (SELECT count(*) FROM projections.departments WHERE clinic_id = ?),
		       updated_at = ?
		 WHERE clinic_id = ?`,
		clinicID, createdAt, clinicID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("clinic_id", ev.GetClinicId()).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.organization_counters
		   SET departments_total = (SELECT count(*) FROM projections.department_counters WHERE organization_id = ?),
		       updated_at = ?
		 WHERE organization_id = ?`,
		orgID, createdAt, orgID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	return nil
}

// DepartmentDetailsChanged updates name + description + mirrors name into employee_cards.
//
// See: docs/services/OrgStructure.md
func DepartmentDetailsChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *deptv1.DepartmentDetailsChanged,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.department", ErrCodeDepartmentProjectionFailed)
	if err != nil {
		return err
	}
	updatedAt := ev.GetUpdatedAt().AsTime()
	var desc null.String
	if ev.GetDescription() != "" {
		desc = null.StringFrom(ev.GetDescription())
	}

	if err := tx.Exec(`
		UPDATE projections.departments
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		ev.GetName(), desc, updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET department_name = ?, updated_at = ?
		 WHERE department_id = ?`,
		ev.GetName(), updatedAt, id,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// DepartmentDeactivated sets is_active=false in projections.departments.
//
// See: docs/services/OrgStructure.md
func DepartmentDeactivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *deptv1.DepartmentDeactivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.department", ErrCodeDepartmentProjectionFailed)
	if err != nil {
		return err
	}
	if err := tx.Exec(`
		UPDATE projections.departments
		   SET is_active = FALSE, updated_at = ?
		 WHERE id = ?`,
		ev.GetUpdatedAt().AsTime(), id,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// DepartmentActivated sets is_active=true in projections.departments.
//
// See: docs/services/OrgStructure.md
func DepartmentActivated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *deptv1.DepartmentActivated,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.department", ErrCodeDepartmentProjectionFailed)
	if err != nil {
		return err
	}
	if err := tx.Exec(`
		UPDATE projections.departments
		   SET is_active = TRUE, updated_at = ?
		 WHERE id = ?`,
		ev.GetUpdatedAt().AsTime(), id,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	return nil
}

// DepartmentDeleted removes the row from projections.departments.
//
// See: docs/services/OrgStructure.md
func DepartmentDeleted(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	_ *deptv1.DepartmentDeleted,
) error {
	id, err := parseUUID(aggregateID, "aggregate_id", "projector.department", ErrCodeDepartmentProjectionFailed)
	if err != nil {
		return err
	}
	if err := tx.Exec(`DELETE FROM projections.departments WHERE id = ?`, id).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", aggregateID).
			Wrap(err)
	}
	return nil
}

func (p *Projectors) DepartmentCreated(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentCreated) error {
	return DepartmentCreated(tx, id, t, ev)
}

func (p *Projectors) DepartmentDetailsChanged(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentDetailsChanged) error {
	return DepartmentDetailsChanged(tx, id, t, ev)
}

func (p *Projectors) DepartmentDeactivated(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentDeactivated) error {
	return DepartmentDeactivated(tx, id, t, ev)
}

func (p *Projectors) DepartmentActivated(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentActivated) error {
	return DepartmentActivated(tx, id, t, ev)
}

func (p *Projectors) DepartmentDeleted(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentDeleted) error {
	return DepartmentDeleted(tx, id, t, ev)
}
