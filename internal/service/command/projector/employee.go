package projector

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// EmployeeHired writes projections.employees + projections.employee_cards
// and bumps employees_total on the organization, clinic, and department
// counter tables. The clinic_id is resolved via the parent department's
// projection; all name fields on the employee_card are populated from
// the existing organisation / clinic / department / users projections.
// Called by EmployeeService.Hire inside the same transaction as the
// domain.employees INSERT.
func EmployeeHired(tx *gorm.DB, e *model.Employee) error {
	clinicID, err := lookupClinicID(tx, e.DepartmentID)
	if err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			With("department_id", e.DepartmentID).
			Wrap(err)
	}

	if err := tx.Exec(`
		INSERT INTO projections.employees
		    (id, zitadel_user_id, organization_id, clinic_id, department_id,
		     position, hired_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.ZitadelUserID, e.OrganizationID, clinicID, e.DepartmentID,
		e.Position, e.CreatedAt, e.CreatedAt, e.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	// Name look-ups from existing projections. Missing rows leave the
	// mirror columns null; the respective rename handlers back-fill
	// lazily so the eventual state is consistent.
	first, last, display, email := lookupUserName(tx, e.ZitadelUserID)
	orgName := lookupOrgName(tx, e.OrganizationID)
	var clinicName null.String
	if clinicID != nil {
		clinicName = lookupClinicName(tx, *clinicID)
	}
	deptName := lookupDepartmentName(tx, e.DepartmentID)

	if err := tx.Exec(`
		INSERT INTO projections.employee_cards
		    (employee_id, zitadel_user_id,
		     first_name, last_name, display_name, email,
		     organization_id, organization_name,
		     clinic_id, clinic_name,
		     department_id, department_name,
		     position, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.ZitadelUserID,
		first, last, display, email,
		e.OrganizationID, orgName,
		clinicID, clinicName,
		e.DepartmentID, deptName,
		e.Position, e.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	if err := bumpOrgEmployees(tx, e.OrganizationID, +1, e.UpdatedAt); err != nil {
		return err
	}
	if clinicID != nil {
		if err := bumpClinicEmployees(tx, *clinicID, +1, e.UpdatedAt); err != nil {
			return err
		}
	}
	if err := bumpDepartmentEmployees(tx, e.DepartmentID, +1, e.UpdatedAt); err != nil {
		return err
	}
	return nil
}

// EmployeeTerminated soft-deletes the employee: sets terminated_at on
// both projections.employees and projections.employee_cards, and
// decrements the three employee-counter tables. The terminated_at stamp
// is taken from the caller-provided `now` to match the domain write.
func EmployeeTerminated(tx *gorm.DB, e *model.Employee, terminatedAt time.Time) error {
	// Load prior state from the projection so we can decrement the
	// correct counters (clinic_id may have changed vs. hire time).
	var clinicID *uuid.UUID
	if err := tx.Raw(
		`SELECT clinic_id FROM projections.employees WHERE id = ?`, e.ID,
	).Row().Scan(&clinicID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.employees
		   SET terminated_at = ?, updated_at = ?
		 WHERE id = ?`,
		terminatedAt, terminatedAt, e.ID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET terminated_at = ?, updated_at = ?
		 WHERE employee_id = ?`,
		terminatedAt, terminatedAt, e.ID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	if err := bumpOrgEmployees(tx, e.OrganizationID, -1, terminatedAt); err != nil {
		return err
	}
	if clinicID != nil {
		if err := bumpClinicEmployees(tx, *clinicID, -1, terminatedAt); err != nil {
			return err
		}
	}
	if err := bumpDepartmentEmployees(tx, e.DepartmentID, -1, terminatedAt); err != nil {
		return err
	}
	return nil
}

// EmployeeDepartmentChanged moves the employee between departments
// (and possibly clinics). The projector reads the prior department_id
// and clinic_id from projections.employees, applies counter deltas
// (decrement old, increment new) for any level whose value actually
// changed, and rewrites department_id/clinic_id + their name mirrors
// on both employees and employee_cards.
func EmployeeDepartmentChanged(tx *gorm.DB, e *model.Employee) error {
	type prior struct {
		ClinicID     *uuid.UUID
		DepartmentID uuid.UUID `validate:"required"`
	}
	var pre prior
	if err := tx.Raw(
		`SELECT clinic_id, department_id FROM projections.employees WHERE id = ?`, e.ID,
	).Row().Scan(&pre.ClinicID, &pre.DepartmentID); err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	newClinicID, err := lookupClinicID(tx, e.DepartmentID)
	if err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			With("department_id", e.DepartmentID).
			Wrap(err)
	}

	if err := tx.Exec(`
		UPDATE projections.employees
		   SET department_id = ?, clinic_id = ?, updated_at = ?
		 WHERE id = ?`,
		e.DepartmentID, newClinicID, e.UpdatedAt, e.ID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	// Update employee_card mirrors. clinic_name is (re)read; nil clinic
	// clears the mirror.
	var newClinicName null.String
	if newClinicID != nil {
		newClinicName = lookupClinicName(tx, *newClinicID)
	}
	newDeptName := lookupDepartmentName(tx, e.DepartmentID)

	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET department_id = ?, department_name = ?,
		       clinic_id = ?, clinic_name = ?,
		       updated_at = ?
		 WHERE employee_id = ?`,
		e.DepartmentID, newDeptName,
		newClinicID, newClinicName,
		e.UpdatedAt, e.ID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}

	// Department counter: always a move (-1 old, +1 new) because the
	// service only publishes when DepartmentID actually changed.
	if pre.DepartmentID != e.DepartmentID {
		if err := bumpDepartmentEmployees(tx, pre.DepartmentID, -1, e.UpdatedAt); err != nil {
			return err
		}
		if err := bumpDepartmentEmployees(tx, e.DepartmentID, +1, e.UpdatedAt); err != nil {
			return err
		}
	}

	// Clinic counter: move only if the clinic actually changed.
	oldPtr, newPtr := pre.ClinicID, newClinicID
	switch {
	case oldPtr == nil && newPtr == nil:
		// no-op
	case oldPtr != nil && newPtr != nil && *oldPtr == *newPtr:
		// no-op
	default:
		if oldPtr != nil {
			if err := bumpClinicEmployees(tx, *oldPtr, -1, e.UpdatedAt); err != nil {
				return err
			}
		}
		if newPtr != nil {
			if err := bumpClinicEmployees(tx, *newPtr, +1, e.UpdatedAt); err != nil {
				return err
			}
		}
	}
	return nil
}

// EmployeePositionChanged rewrites the position column on both
// projections.employees and projections.employee_cards. No counter
// bumps (position is metadata only).
func EmployeePositionChanged(tx *gorm.DB, e *model.Employee) error {
	if err := tx.Exec(`
		UPDATE projections.employees
		   SET position = ?, updated_at = ?
		 WHERE id = ?`,
		e.Position, e.UpdatedAt, e.ID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}
	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET position = ?, updated_at = ?
		 WHERE employee_id = ?`,
		e.Position, e.UpdatedAt, e.ID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("employee_id", e.ID).
			Wrap(err)
	}
	return nil
}

// lookupClinicID reads clinic_id from projections.departments for the
// given department. The department MUST already be projected because
// the sync-projector contract orders parent writes before child writes
// inside the same transaction — a missing row is a bug, not a valid
// state, and is surfaced as an oops error.
//
// A NULL clinic_id column (department directly under the organization)
// is a legitimate result and returns (nil, nil).
func lookupClinicID(tx *gorm.DB, deptID uuid.UUID) (*uuid.UUID, error) {
	var clinicID uuid.UUID
	err := tx.Raw(
		`SELECT clinic_id FROM projections.departments WHERE id = ?`, deptID,
	).Row().Scan(&clinicID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("projector.employee").
				Code(ErrCodeEmployeeProjectionFailed).
				With("department_id", deptID).
				Errorf("department projection not found")
		}
		return nil, err
	}
	if clinicID == uuid.Nil {
		return nil, nil //nolint:nilnil // NULL clinic_id — department directly under the org
	}
	return &clinicID, nil
}

// lookupUserName reads the four user-display columns from the existing
// Zitadel users projection. Missing rows return zero-valued null
// strings — the respective rename handlers back-fill lazily.
func lookupUserName(tx *gorm.DB, zitadelID string) (first, last, display, email null.String) {
	_ = tx.Raw(
		`SELECT first_name, last_name, display_name, email FROM projections.users WHERE id = ?`,
		zitadelID,
	).Row().Scan(&first, &last, &display, &email)
	return first, last, display, email
}

// lookupOrgName reads projections.organizations.name; returns an
// Invalid null.String when the row is absent.
func lookupOrgName(tx *gorm.DB, orgID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.organizations WHERE id = ?`, orgID).
		Row().Scan(&name)
	return name
}

// lookupClinicName reads projections.clinics.name; returns Invalid on miss.
func lookupClinicName(tx *gorm.DB, clinicID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.clinics WHERE id = ?`, clinicID).
		Row().Scan(&name)
	return name
}

// lookupDepartmentName reads projections.departments.name; returns
// Invalid on miss.
func lookupDepartmentName(tx *gorm.DB, deptID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.departments WHERE id = ?`, deptID).
		Row().Scan(&name)
	return name
}

// bumpOrgEmployees applies a delta to the employees_total column on
// projections.organization_counters for the given organization. The
// row must already exist (created by OrganizationCreated); delta
// may be negative.
func bumpOrgEmployees(tx *gorm.DB, orgID uuid.UUID, delta int, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.organization_counters
		   SET employees_total = employees_total + ?, updated_at = ?
		 WHERE organization_id = ?`,
		delta, now, orgID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	return nil
}

// bumpClinicEmployees applies a delta to the employees_total column on
// projections.clinic_counters. The row must already exist.
func bumpClinicEmployees(tx *gorm.DB, clinicID uuid.UUID, delta int, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.clinic_counters
		   SET employees_total = employees_total + ?, updated_at = ?
		 WHERE clinic_id = ?`,
		delta, now, clinicID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	return nil
}

// bumpDepartmentEmployees applies a delta to the employees_total column
// on projections.department_counters. The row must already exist.
func bumpDepartmentEmployees(tx *gorm.DB, deptID uuid.UUID, delta int, now time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.department_counters
		   SET employees_total = employees_total + ?, updated_at = ?
		 WHERE department_id = ?`,
		delta, now, deptID,
	).Error; err != nil {
		return oops.In("projector.employee").
			Code(ErrCodeEmployeeProjectionFailed).
			With("department_id", deptID).
			Wrap(err)
	}
	return nil
}
