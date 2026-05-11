package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
)

// EmployeeHired writes projections.employees + projections.employee_cards
// and bumps employee counters on org/clinic/dept.
//
// See: docs/services/Membership.md
func EmployeeHired(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *empv1.EmployeeHired,
) error {
	empID := uuid.MustParse(aggregateID)
	orgID := uuid.MustParse(ev.GetOrganizationId())
	deptID := uuid.MustParse(ev.GetDepartmentId())
	hiredAt := ev.GetHiredAt().AsTime()

	clinicID, err := lookupClinicID(tx, deptID)
	if err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", aggregateID).
			Wrap(err)
	}

	var pos null.String
	if ev.GetPosition() != "" {
		pos = null.StringFrom(ev.GetPosition())
	}

	if err := tx.Exec(`
		INSERT INTO projections.employees
		    (id, zitadel_user_id, organization_id, clinic_id, department_id,
		     position, hired_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING`,
		empID, ev.GetZitadelUserId(), orgID, clinicID, deptID,
		pos, hiredAt, hiredAt, hiredAt,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", aggregateID).
			Wrap(err)
	}

	first, last, display, email := lookupUserName(tx, ev.GetZitadelUserId())
	orgName := lookupOrgName(tx, orgID)
	var clinicName null.String
	if clinicID != nil {
		clinicName = lookupClinicName(tx, *clinicID)
	}
	deptName := lookupDepartmentName(tx, deptID)

	if err := tx.Exec(`
		INSERT INTO projections.employee_cards
		    (employee_id, zitadel_user_id,
		     first_name, last_name, display_name, email,
		     organization_id, organization_name,
		     clinic_id, clinic_name,
		     department_id, department_name,
		     position, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (employee_id) DO NOTHING`,
		empID, ev.GetZitadelUserId(),
		first, last, display, email,
		orgID, orgName,
		clinicID, clinicName,
		deptID, deptName,
		pos, hiredAt,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", aggregateID).
			Wrap(err)
	}

	if err := bumpOrgEmployees(tx, orgID, +1, hiredAt); err != nil {
		return err
	}
	if clinicID != nil {
		if err := bumpClinicEmployees(tx, *clinicID, +1, hiredAt); err != nil {
			return err
		}
	}
	return bumpDepartmentEmployees(tx, deptID, +1, hiredAt)
}

// EmployeeTerminated sets terminated_at and decrements counters.
//
// See: docs/services/Membership.md
func EmployeeTerminated(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *empv1.EmployeeTerminated,
) error {
	empID := uuid.MustParse(aggregateID)
	orgID := uuid.MustParse(ev.GetOrganizationId())
	deptID := uuid.MustParse(ev.GetDepartmentId())
	terminatedAt := ev.GetTerminatedAt().AsTime()

	// Read current clinic_id from projection to decrement the correct counter.
	var clinicID *uuid.UUID
	_ = tx.Raw(`SELECT clinic_id FROM projections.employees WHERE id = ?`, empID).
		Row().Scan(&clinicID)

	if err := tx.Exec(`UPDATE projections.employees SET terminated_at = ?, updated_at = ? WHERE id = ?`,
		terminatedAt, terminatedAt, empID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}
	if err := tx.Exec(`UPDATE projections.employee_cards SET terminated_at = ?, updated_at = ? WHERE employee_id = ?`,
		terminatedAt, terminatedAt, empID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}

	if err := bumpOrgEmployees(tx, orgID, -1, terminatedAt); err != nil {
		return err
	}
	if clinicID != nil {
		if err := bumpClinicEmployees(tx, *clinicID, -1, terminatedAt); err != nil {
			return err
		}
	}
	return bumpDepartmentEmployees(tx, deptID, -1, terminatedAt)
}

// EmployeeDepartmentChanged moves employee to new department and adjusts counters.
//
// See: docs/services/Membership.md
func EmployeeDepartmentChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *empv1.EmployeeDepartmentChanged,
) error {
	empID := uuid.MustParse(aggregateID)
	newDeptID := uuid.MustParse(ev.GetNewDepartmentId())
	updatedAt := ev.GetUpdatedAt().AsTime()

	type prior struct {
		ClinicID     *uuid.UUID
		DepartmentID uuid.UUID
	}
	var pre prior
	if err := tx.Raw(`SELECT clinic_id, department_id FROM projections.employees WHERE id = ?`, empID).
		Row().Scan(&pre.ClinicID, &pre.DepartmentID); err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}

	newClinicID, err := lookupClinicID(tx, newDeptID)
	if err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}

	if err := tx.Exec(`UPDATE projections.employees SET department_id = ?, clinic_id = ?, updated_at = ? WHERE id = ?`,
		newDeptID, newClinicID, updatedAt, empID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}

	var newClinicName null.String
	if newClinicID != nil {
		newClinicName = lookupClinicName(tx, *newClinicID)
	}
	newDeptName := lookupDepartmentName(tx, newDeptID)

	if err := tx.Exec(`UPDATE projections.employee_cards
		   SET department_id = ?, department_name = ?, clinic_id = ?, clinic_name = ?, updated_at = ?
		 WHERE employee_id = ?`,
		newDeptID, newDeptName, newClinicID, newClinicName, updatedAt, empID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}

	if pre.DepartmentID != newDeptID {
		if err := bumpDepartmentEmployees(tx, pre.DepartmentID, -1, updatedAt); err != nil {
			return err
		}
		if err := bumpDepartmentEmployees(tx, newDeptID, +1, updatedAt); err != nil {
			return err
		}
	}
	oldPtr, newPtr := pre.ClinicID, newClinicID
	switch {
	case oldPtr == nil && newPtr == nil:
	case oldPtr != nil && newPtr != nil && *oldPtr == *newPtr:
	default:
		if oldPtr != nil {
			if err := bumpClinicEmployees(tx, *oldPtr, -1, updatedAt); err != nil {
				return err
			}
		}
		if newPtr != nil {
			if err := bumpClinicEmployees(tx, *newPtr, +1, updatedAt); err != nil {
				return err
			}
		}
	}
	return nil
}

// EmployeePositionChanged updates position on employees + employee_cards.
//
// See: docs/services/Membership.md
func EmployeePositionChanged(
	tx *gorm.DB,
	aggregateID string,
	_ time.Time,
	ev *empv1.EmployeePositionChanged,
) error {
	empID := uuid.MustParse(aggregateID)
	updatedAt := ev.GetUpdatedAt().AsTime()
	var pos null.String
	if ev.GetPosition() != "" {
		pos = null.StringFrom(ev.GetPosition())
	}

	if err := tx.Exec(`UPDATE projections.employees SET position = ?, updated_at = ? WHERE id = ?`,
		pos, updatedAt, empID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}
	if err := tx.Exec(`UPDATE projections.employee_cards SET position = ?, updated_at = ? WHERE employee_id = ?`,
		pos, updatedAt, empID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("employee_id", aggregateID).Wrap(err)
	}
	return nil
}

func bumpOrgEmployees(tx *gorm.DB, orgID uuid.UUID, delta int, now time.Time) error {
	if err := tx.Exec(`UPDATE projections.organization_counters SET employees_total = employees_total + ?, updated_at = ? WHERE organization_id = ?`,
		delta, now, orgID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("organization_id", orgID).Wrap(err)
	}
	return nil
}

func bumpClinicEmployees(tx *gorm.DB, clinicID uuid.UUID, delta int, now time.Time) error {
	if err := tx.Exec(`UPDATE projections.clinic_counters SET employees_total = employees_total + ?, updated_at = ? WHERE clinic_id = ?`,
		delta, now, clinicID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("clinic_id", clinicID).Wrap(err)
	}
	return nil
}

func bumpDepartmentEmployees(tx *gorm.DB, deptID uuid.UUID, delta int, now time.Time) error {
	if err := tx.Exec(`UPDATE projections.department_counters SET employees_total = employees_total + ?, updated_at = ? WHERE department_id = ?`,
		delta, now, deptID).Error; err != nil {
		return oops.In("projector.employee").Code(ErrCodeEmployeeProjectionFailed).With("department_id", deptID).Wrap(err)
	}
	return nil
}

// Forwarding methods on *Projectors

func (p *Projectors) EmployeeHired(tx *gorm.DB, id string, t time.Time, ev *empv1.EmployeeHired) error {
	return EmployeeHired(tx, id, t, ev)
}

func (p *Projectors) EmployeeTerminated(tx *gorm.DB, id string, t time.Time, ev *empv1.EmployeeTerminated) error {
	return EmployeeTerminated(tx, id, t, ev)
}

func (p *Projectors) EmployeeDepartmentChanged(tx *gorm.DB, id string, t time.Time, ev *empv1.EmployeeDepartmentChanged) error {
	return EmployeeDepartmentChanged(tx, id, t, ev)
}

func (p *Projectors) EmployeePositionChanged(tx *gorm.DB, id string, t time.Time, ev *empv1.EmployeePositionChanged) error {
	return EmployeePositionChanged(tx, id, t, ev)
}
