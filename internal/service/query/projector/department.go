package projector

import (
	"time"

	"github.com/google/uuid"
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
	id := uuid.MustParse(aggregateID)
	clinicID := uuid.MustParse(ev.GetClinicId())
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
		   SET departments_total = departments_total + 1, updated_at = ?
		 WHERE clinic_id = ?`,
		createdAt, clinicID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("clinic_id", ev.GetClinicId()).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.organization_counters
		   SET departments_total = departments_total + 1, updated_at = ?
		 WHERE organization_id = ?`,
		createdAt, orgID,
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
	id := uuid.MustParse(aggregateID)
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

func (p *Projectors) DepartmentCreated(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentCreated) error {
	return DepartmentCreated(tx, id, t, ev)
}

func (p *Projectors) DepartmentDetailsChanged(tx *gorm.DB, id string, t time.Time, ev *deptv1.DepartmentDetailsChanged) error {
	return DepartmentDetailsChanged(tx, id, t, ev)
}
