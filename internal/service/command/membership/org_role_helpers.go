package membership

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
)

// Shared sub-helpers used by the five role-assignment commands
// (AssignClinicHead, AssignDepartmentResponsible,
// AssignOrganizationAdmin, AssignOrganizationHead,
// AssignOrganizationDispatcher). Each helper folds one piece of the
// assign-flow boilerplate that used to live inline in every command
// file, so the command bodies read as "validate → exist → load emp →
// scope → insert" instead of 100 lines of copy-paste.

// requireOrganizationExists fails with ErrCodeOrganizationNotFound if
// the row is missing, or wraps other driver errors with
// ErrCodeOrganizationLookupFailed.
func requireOrganizationExists(tx *gorm.DB, scope string, orgID uuid.UUID) error {
	var exists int64
	if err := tx.Raw(`SELECT count(*) FROM domain.organizations WHERE id = ?`, orgID).
		Scan(&exists).Error; err != nil {
		return oops.In(scope).Code(ErrCodeOrganizationLookupFailed).Wrap(err)
	}
	if exists == 0 {
		return oops.In(scope).
			Code(ErrCodeOrganizationNotFound).
			Public("Organization not found.").
			With("organization_id", orgID).
			Errorf("organization not found")
	}
	return nil
}

// requireClinicExists mirrors requireOrganizationExists for clinics.
// Used by clinic_head_assign.
func requireClinicExists(tx *gorm.DB, scope string, clinicID uuid.UUID) error {
	var exists int64
	if err := tx.Raw(`SELECT count(*) FROM domain.clinics WHERE id = ?`, clinicID).
		Scan(&exists).Error; err != nil {
		return oops.In(scope).Code(ErrCodeClinicLookupFailed).Wrap(err)
	}
	if exists == 0 {
		return oops.In(scope).
			Code(ErrCodeClinicNotFound).
			Public("Clinic not found.").
			With("clinic_id", clinicID).
			Errorf("clinic not found")
	}
	return nil
}

// requireDepartmentExists mirrors requireOrganizationExists for
// departments. Used by department_responsible_assign.
func requireDepartmentExists(tx *gorm.DB, scope string, deptID uuid.UUID) error {
	var exists int64
	if err := tx.Raw(`SELECT count(*) FROM domain.departments WHERE id = ?`, deptID).
		Scan(&exists).Error; err != nil {
		return oops.In(scope).Code(ErrCodeDepartmentLookupFailed).Wrap(err)
	}
	if exists == 0 {
		return oops.In(scope).
			Code(ErrCodeDepartmentNotFound).
			Public("Department not found.").
			With("department_id", deptID).
			Errorf("department not found")
	}
	return nil
}

// loadEmployeeForRoleAssignment loads an employee row with FOR UPDATE
// and wraps NotFound/other errors with the right oops code under the
// caller's scope.
func loadEmployeeForRoleAssignment(tx *gorm.DB, scope string, empID uuid.UUID) (*model.Employee, error) {
	var emp model.Employee
	err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
		Where("id = ?", empID).
		First(&emp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, oops.In(scope).
				Code(ErrCodeEmployeeNotFound).
				Public("Employee not found.").
				With("employee_id", empID).
				Errorf("employee not found")
		}
		return nil, oops.In(scope).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
	}
	return &emp, nil
}

// requireEmployeeInOrganization enforces the "employee belongs to the
// target organization" invariant via the denormalized
// employees.organization_id column.
func requireEmployeeInOrganization(scope string, emp *model.Employee, orgID uuid.UUID) error {
	if emp.OrganizationID == orgID {
		return nil
	}
	return oops.In(scope).
		Code(ErrCodeEmployeeNotInOrganization).
		Public("Employee does not belong to this organization.").
		With("employee_id", emp.ID).
		With("employee_organization_id", emp.OrganizationID).
		With("target_organization_id", orgID).
		Errorf("employee not in target organization")
}

// requireEmployeeInDepartment enforces the same-department invariant
// via the denormalized employees.department_id column.
func requireEmployeeInDepartment(scope string, emp *model.Employee, deptID uuid.UUID) error {
	if emp.DepartmentID == deptID {
		return nil
	}
	return oops.In(scope).
		Code(ErrCodeEmployeeNotInDepartment).
		Public("Employee does not belong to this department.").
		With("employee_id", emp.ID).
		With("employee_department_id", emp.DepartmentID).
		With("target_department_id", deptID).
		Errorf("employee not in target department")
}

// requireEmployeeInClinic enforces the same-clinic invariant by
// looking up the department and checking its clinic_id, since
// clinic_id is not denormalized on employees.
func requireEmployeeInClinic(tx *gorm.DB, scope string, emp *model.Employee, clinicID uuid.UUID) error {
	var inClinic int64
	if err := tx.Raw(
		`SELECT count(*) FROM domain.departments WHERE id = ? AND clinic_id = ?`,
		emp.DepartmentID, clinicID,
	).Scan(&inClinic).Error; err != nil {
		return oops.In(scope).Code(ErrCodeDepartmentLookupFailed).Wrap(err)
	}
	if inClinic == 0 {
		return oops.In(scope).
			Code(ErrCodeEmployeeNotInClinic).
			Public("Employee does not belong to this clinic.").
			With("employee_id", emp.ID).
			With("employee_department_id", emp.DepartmentID).
			With("target_clinic_id", clinicID).
			Errorf("employee not in target clinic")
	}
	return nil
}
