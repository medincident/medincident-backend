package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// Role projectors mirror the domain.<role> tables into projections.<role>.
// Each function applies one domain mutation inside the caller's
// transaction — holders of the five org-structure roles have a common
// (scope_id, employee_id) composite key plus an optional deputy slot;
// system_admin is a simple (zitadel_user_id) key.

// ClinicHeadAssigned inserts the projection row for a freshly created
// clinic head. Called after the domain INSERT.
func ClinicHeadAssigned(tx *gorm.DB, r *model.ClinicHead) error {
	if err := tx.Exec(`
		INSERT INTO projections.clinic_heads
		    (clinic_id, employee_id, deputy_employee_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		r.ClinicID, r.EmployeeID, r.DeputyEmployeeID,
		r.CreatedAt, r.UpdatedAt,
	).Error; err != nil {
		return wrapRole(err, "clinic_head", r.EmployeeID)
	}
	return nil
}

// ClinicHeadDeputyAssigned sets the deputy slot on an existing clinic
// head row.
func ClinicHeadDeputyAssigned(tx *gorm.DB, r *model.ClinicHead) error {
	if err := tx.Exec(`
		UPDATE projections.clinic_heads
		   SET deputy_employee_id = ?, updated_at = ?
		 WHERE clinic_id = ? AND employee_id = ?`,
		r.DeputyEmployeeID, r.UpdatedAt, r.ClinicID, r.EmployeeID,
	).Error; err != nil {
		return wrapRole(err, "clinic_head", r.EmployeeID)
	}
	return nil
}

// ClinicHeadDeputyRemoved clears the deputy slot on an existing clinic
// head row.
func ClinicHeadDeputyRemoved(tx *gorm.DB, clinicID, employeeID uuid.UUID, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.clinic_heads
		   SET deputy_employee_id = NULL, updated_at = ?
		 WHERE clinic_id = ? AND employee_id = ?`,
		now, clinicID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "clinic_head", employeeID)
	}
	return nil
}

// ClinicHeadRevoked deletes the projection row.
func ClinicHeadRevoked(tx *gorm.DB, clinicID, employeeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		clinicID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "clinic_head", employeeID)
	}
	return nil
}

// DepartmentResponsibleAssigned inserts the projection row for a
// freshly created department responsible.
func DepartmentResponsibleAssigned(tx *gorm.DB, r *model.DepartmentResponsible) error {
	if err := tx.Exec(`
		INSERT INTO projections.department_responsibles
		    (department_id, employee_id, deputy_employee_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		r.DepartmentID, r.EmployeeID, r.DeputyEmployeeID,
		r.CreatedAt, r.UpdatedAt,
	).Error; err != nil {
		return wrapRole(err, "department_responsible", r.EmployeeID)
	}
	return nil
}

// DepartmentResponsibleDeputyAssigned sets the deputy slot.
func DepartmentResponsibleDeputyAssigned(tx *gorm.DB, r *model.DepartmentResponsible) error {
	if err := tx.Exec(`
		UPDATE projections.department_responsibles
		   SET deputy_employee_id = ?, updated_at = ?
		 WHERE department_id = ? AND employee_id = ?`,
		r.DeputyEmployeeID, r.UpdatedAt, r.DepartmentID, r.EmployeeID,
	).Error; err != nil {
		return wrapRole(err, "department_responsible", r.EmployeeID)
	}
	return nil
}

// DepartmentResponsibleDeputyRemoved clears the deputy slot.
func DepartmentResponsibleDeputyRemoved(tx *gorm.DB, departmentID, employeeID uuid.UUID, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.department_responsibles
		   SET deputy_employee_id = NULL, updated_at = ?
		 WHERE department_id = ? AND employee_id = ?`,
		now, departmentID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "department_responsible", employeeID)
	}
	return nil
}

// DepartmentResponsibleRevoked deletes the projection row.
func DepartmentResponsibleRevoked(tx *gorm.DB, departmentID, employeeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		departmentID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "department_responsible", employeeID)
	}
	return nil
}

// OrgAdminAssigned inserts the projection row for a freshly created
// org admin role.
func OrgAdminAssigned(tx *gorm.DB, r *model.OrgAdmin) error {
	if err := tx.Exec(`
		INSERT INTO projections.org_admins
		    (organization_id, employee_id, deputy_employee_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		r.OrganizationID, r.EmployeeID, r.DeputyEmployeeID,
		r.CreatedAt, r.UpdatedAt,
	).Error; err != nil {
		return wrapRole(err, "org_admin", r.EmployeeID)
	}
	return nil
}

// OrgAdminDeputyAssigned sets the deputy slot.
func OrgAdminDeputyAssigned(tx *gorm.DB, r *model.OrgAdmin) error {
	if err := tx.Exec(`
		UPDATE projections.org_admins
		   SET deputy_employee_id = ?, updated_at = ?
		 WHERE organization_id = ? AND employee_id = ?`,
		r.DeputyEmployeeID, r.UpdatedAt, r.OrganizationID, r.EmployeeID,
	).Error; err != nil {
		return wrapRole(err, "org_admin", r.EmployeeID)
	}
	return nil
}

// OrgAdminDeputyRemoved clears the deputy slot.
func OrgAdminDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.org_admins
		   SET deputy_employee_id = NULL, updated_at = ?
		 WHERE organization_id = ? AND employee_id = ?`,
		now, organizationID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "org_admin", employeeID)
	}
	return nil
}

// OrgAdminRevoked deletes the projection row.
func OrgAdminRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.org_admins WHERE organization_id = ? AND employee_id = ?`,
		organizationID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "org_admin", employeeID)
	}
	return nil
}

// OrgDispatcherAssigned inserts the projection row.
func OrgDispatcherAssigned(tx *gorm.DB, r *model.OrgDispatcher) error {
	if err := tx.Exec(`
		INSERT INTO projections.org_dispatchers
		    (organization_id, employee_id, deputy_employee_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		r.OrganizationID, r.EmployeeID, r.DeputyEmployeeID,
		r.CreatedAt, r.UpdatedAt,
	).Error; err != nil {
		return wrapRole(err, "org_dispatcher", r.EmployeeID)
	}
	return nil
}

// OrgDispatcherDeputyAssigned sets the deputy slot.
func OrgDispatcherDeputyAssigned(tx *gorm.DB, r *model.OrgDispatcher) error {
	if err := tx.Exec(`
		UPDATE projections.org_dispatchers
		   SET deputy_employee_id = ?, updated_at = ?
		 WHERE organization_id = ? AND employee_id = ?`,
		r.DeputyEmployeeID, r.UpdatedAt, r.OrganizationID, r.EmployeeID,
	).Error; err != nil {
		return wrapRole(err, "org_dispatcher", r.EmployeeID)
	}
	return nil
}

// OrgDispatcherDeputyRemoved clears the deputy slot.
func OrgDispatcherDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.org_dispatchers
		   SET deputy_employee_id = NULL, updated_at = ?
		 WHERE organization_id = ? AND employee_id = ?`,
		now, organizationID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "org_dispatcher", employeeID)
	}
	return nil
}

// OrgDispatcherRevoked deletes the projection row.
func OrgDispatcherRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.org_dispatchers WHERE organization_id = ? AND employee_id = ?`,
		organizationID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "org_dispatcher", employeeID)
	}
	return nil
}

// OrgHeadAssigned inserts the projection row.
func OrgHeadAssigned(tx *gorm.DB, r *model.OrgHead) error {
	if err := tx.Exec(`
		INSERT INTO projections.org_heads
		    (organization_id, employee_id, deputy_employee_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		r.OrganizationID, r.EmployeeID, r.DeputyEmployeeID,
		r.CreatedAt, r.UpdatedAt,
	).Error; err != nil {
		return wrapRole(err, "org_head", r.EmployeeID)
	}
	return nil
}

// OrgHeadDeputyAssigned sets the deputy slot.
func OrgHeadDeputyAssigned(tx *gorm.DB, r *model.OrgHead) error {
	if err := tx.Exec(`
		UPDATE projections.org_heads
		   SET deputy_employee_id = ?, updated_at = ?
		 WHERE organization_id = ? AND employee_id = ?`,
		r.DeputyEmployeeID, r.UpdatedAt, r.OrganizationID, r.EmployeeID,
	).Error; err != nil {
		return wrapRole(err, "org_head", r.EmployeeID)
	}
	return nil
}

// OrgHeadDeputyRemoved clears the deputy slot.
func OrgHeadDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.org_heads
		   SET deputy_employee_id = NULL, updated_at = ?
		 WHERE organization_id = ? AND employee_id = ?`,
		now, organizationID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "org_head", employeeID)
	}
	return nil
}

// OrgHeadRevoked deletes the projection row.
func OrgHeadRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID) error {
	if err := tx.Exec(
		`DELETE FROM projections.org_heads WHERE organization_id = ? AND employee_id = ?`,
		organizationID, employeeID,
	).Error; err != nil {
		return wrapRole(err, "org_head", employeeID)
	}
	return nil
}

// SystemAdminGranted inserts a row in projections.system_admins.
func SystemAdminGranted(tx *gorm.DB, zitadelUserID string, createdAt time.Time) error {
	if err := tx.Exec(`
		INSERT INTO projections.system_admins (zitadel_user_id, created_at)
		VALUES (?, ?)`,
		zitadelUserID, createdAt,
	).Error; err != nil {
		return oops.In("projector.role").
			Code(ErrCodeRoleProjectionFailed).
			With("zitadel_user_id", zitadelUserID).
			Wrap(err)
	}
	return nil
}

// SystemAdminRevoked deletes a row from projections.system_admins.
func SystemAdminRevoked(tx *gorm.DB, zitadelUserID string) error {
	if err := tx.Exec(
		`DELETE FROM projections.system_admins WHERE zitadel_user_id = ?`,
		zitadelUserID,
	).Error; err != nil {
		return oops.In("projector.role").
			Code(ErrCodeRoleProjectionFailed).
			With("zitadel_user_id", zitadelUserID).
			Wrap(err)
	}
	return nil
}

// wrapRole wraps a gorm error with the per-role oops scope.
func wrapRole(err error, role string, employeeID uuid.UUID) error {
	return oops.In("projector.role."+role).
		Code(ErrCodeRoleProjectionFailed).
		With("employee_id", employeeID).
		Wrap(err)
}
