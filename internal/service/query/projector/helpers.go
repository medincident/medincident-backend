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

// resolveRegistrarLocation returns (orgID, clinicID, deptID) for the registrar.
// When registrar location fields are present in the event (post-fix events),
// they are parsed directly. When they are absent (pre-fix events that carry
// empty strings), the function falls back to querying projections.employees by
// zitadel_user_id so those events can be replayed without a permanent Term().
// If the employee projection is also absent, the incident's own values are
// used as a last resort.
func resolveRegistrarLocation(
	tx *gorm.DB,
	rawOrgID, rawClinicID, rawDeptID string,
	zitadelUserID string,
	incidentOrgID, incidentClinicID, incidentDeptID uuid.UUID,
) (orgID, clinicID, deptID uuid.UUID) {
	if rawOrgID != "" && rawClinicID != "" && rawDeptID != "" {
		orgID, _ = uuid.Parse(rawOrgID)
		clinicID, _ = uuid.Parse(rawClinicID)
		deptID, _ = uuid.Parse(rawDeptID)
		if orgID != uuid.Nil && clinicID != uuid.Nil && deptID != uuid.Nil {
			return orgID, clinicID, deptID
		}
	}
	// Pre-fix event: look up from the employee projection.
	orgID = incidentOrgID
	var row struct {
		DepartmentID uuid.UUID
		ClinicID     *uuid.UUID
	}
	_ = tx.Raw(
		`SELECT department_id, clinic_id FROM projections.employees
		  WHERE zitadel_user_id = ? AND organization_id = ? AND terminated_at IS NULL
		  LIMIT 1`,
		zitadelUserID, incidentOrgID,
	).Scan(&row)
	if row.DepartmentID != uuid.Nil {
		deptID = row.DepartmentID
		if row.ClinicID != nil {
			clinicID = *row.ClinicID
		} else {
			clinicID = incidentClinicID
		}
		return orgID, clinicID, deptID
	}
	return incidentOrgID, incidentClinicID, incidentDeptID
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
