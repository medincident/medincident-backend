package projector

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
)

// lookupClinicID reads clinic_id from projections.departments for the
// given department. Returns (nil, nil) when the department has no clinic
// (directly under the org). Returns an error when the department row is
// absent — callers must ensure parent projections are applied before
// children, which the strict-seq delivery guarantees.
func lookupClinicID(tx *gorm.DB, deptID uuid.UUID) (*uuid.UUID, error) {
	var clinicID *uuid.UUID
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
	return clinicID, nil // nil when SQL NULL (dept directly under org)
}

// lookupUserName reads the four user-display columns from projections.users.
// Missing rows return zero-valued null.Strings — lazily back-filled later.
func lookupUserName(tx *gorm.DB, zitadelID string) (first, last, display, email null.String) {
	_ = tx.Raw(
		`SELECT first_name, last_name, display_name, email FROM projections.users WHERE id = ?`,
		zitadelID,
	).Row().Scan(&first, &last, &display, &email)
	return first, last, display, email
}

// lookupUserDisplayName returns projections.users.display_name, or empty
// string when the row is absent.
func lookupUserDisplayName(tx *gorm.DB, zitadelID string) string {
	var name string
	_ = tx.Raw(`SELECT display_name FROM projections.users WHERE id = ?`, zitadelID).
		Row().Scan(&name)
	return name
}

func lookupOrgName(tx *gorm.DB, orgID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.organizations WHERE id = ?`, orgID).
		Row().Scan(&name)
	return name
}

func lookupClinicName(tx *gorm.DB, clinicID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.clinics WHERE id = ?`, clinicID).
		Row().Scan(&name)
	return name
}

func lookupDepartmentName(tx *gorm.DB, deptID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.departments WHERE id = ?`, deptID).
		Row().Scan(&name)
	return name
}

func lookupOrgIDForClinic(tx *gorm.DB, clinicID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := tx.Raw(
		`SELECT organization_id FROM projections.clinics WHERE id = ?`, clinicID,
	).Row().Scan(&orgID); err != nil {
		return uuid.Nil, oops.In("projector.clinic").
			Code(ErrCodeClinicProjectionFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	return orgID, nil
}

// registrarSnapshot holds the projection-side employee fields needed for
// the projections.incidents insert. Fetched via JOIN so clinic_id is
// always resolved even when projections.employees.clinic_id is NULL.
type registrarSnapshot struct {
	organizationID uuid.UUID
	clinicID       uuid.UUID
	departmentID   uuid.UUID
	position       null.String
}

// lookupRegistrarSnapshot fetches org/clinic/dept/position for an employee
// from projections.employees (JOIN projections.departments for clinic_id).
// Returns a transient (_failed) error when the row is absent so the
// consumer NAKs and retries once the EmployeeHired projection lands.
func lookupRegistrarSnapshot(tx *gorm.DB, employeeID uuid.UUID) (registrarSnapshot, error) {
	var snap registrarSnapshot
	err := tx.Raw(`
		SELECT e.organization_id,
		       COALESCE(e.clinic_id, d.clinic_id) AS clinic_id,
		       e.department_id,
		       e.position
		FROM projections.employees e
		JOIN projections.departments d ON d.id = e.department_id
		WHERE e.id = ?`, employeeID,
	).Row().Scan(&snap.organizationID, &snap.clinicID, &snap.departmentID, &snap.position)
	if err != nil {
		return registrarSnapshot{}, oops.In("projector.incident_lifecycle").
			Code(ErrCodeIncidentLifecycleProjectionFailed).
			With("registrar_employee_id", employeeID).
			Wrap(err)
	}
	return snap, nil
}

// parseUUID parses s as a UUID and returns a permanent (_malformed) error on
// failure so the dispatcher calls Term() instead of NakWithDelay.
func parseUUID(s, fieldName, oopsIn, oopsCode string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, oops.In(oopsIn).
			Code(oopsCode+"_malformed").
			With(fieldName, s).
			Wrap(err)
	}
	return id, nil
}
