package projector

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
)

// DepartmentCreated writes the projections.departments row, its matching
// projections.department_counters row (initialised to zero), bumps
// projections.clinic_counters.departments_total for the parent clinic,
// and bumps projections.organization_counters.departments_total for
// the grandparent organization. Called by DepartmentService.Create
// inside the same transaction as the domain.departments INSERT. The
// organization id is looked up via projections.clinics so the caller
// does not need to pre-load it.
func DepartmentCreated(tx *gorm.DB, d *model.Department) error {
	if err := tx.Exec(`
		INSERT INTO projections.departments
		    (id, clinic_id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		d.ID, d.ClinicID, d.Name, d.Description,
		d.CreatedAt, d.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", d.ID).
			Wrap(err)
	}

	// Parent clinic is always already projected (sync projector in the
	// same tx chain): look up organization_id to populate the counters.
	var orgID uuid.UUID
	if err := tx.Raw(
		`SELECT organization_id FROM projections.clinics WHERE id = ?`,
		d.ClinicID,
	).Row().Scan(&orgID); err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", d.ID).
			With("clinic_id", d.ClinicID).
			Wrap(err)
	}

	if err := tx.Exec(`
		INSERT INTO projections.department_counters
		    (department_id, clinic_id, organization_id, employees_total, updated_at)
		VALUES (?, ?, ?, 0, ?)`,
		d.ID, d.ClinicID, orgID, d.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", d.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.clinic_counters
		   SET departments_total = departments_total + 1, updated_at = ?
		 WHERE clinic_id = ?`,
		d.UpdatedAt, d.ClinicID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("clinic_id", d.ClinicID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.organization_counters
		   SET departments_total = departments_total + 1, updated_at = ?
		 WHERE organization_id = ?`,
		d.UpdatedAt, orgID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	return nil
}

// DepartmentDetailsChanged updates name + description + updated_at on
// projections.departments and mirrors the new department_name into
// projections.employee_cards for every card under this department.
func DepartmentDetailsChanged(tx *gorm.DB, d *model.Department) error {
	if err := tx.Exec(`
		UPDATE projections.departments
		   SET name = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		d.Name, d.Description, d.UpdatedAt, d.ID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", d.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET department_name = ?, updated_at = ?
		 WHERE department_id = ?`,
		d.Name, d.UpdatedAt, d.ID,
	).Error; err != nil {
		return oops.In("projector.department").
			Code(ErrCodeDepartmentProjectionFailed).
			With("department_id", d.ID).
			Wrap(err)
	}
	return nil
}
