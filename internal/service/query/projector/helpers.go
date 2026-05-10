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
//
//nolint:unused // used by membership projectors added in a later task
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
//
//nolint:unused // used by membership projectors added in a later task
func lookupUserName(tx *gorm.DB, zitadelID string) (first, last, display, email null.String) {
	_ = tx.Raw(
		`SELECT first_name, last_name, display_name, email FROM projections.users WHERE id = ?`,
		zitadelID,
	).Row().Scan(&first, &last, &display, &email)
	return first, last, display, email
}

// lookupUserDisplayName returns projections.users.display_name, or empty
// string when the row is absent.
//
//nolint:unused // used by membership projectors added in a later task
func lookupUserDisplayName(tx *gorm.DB, zitadelID string) string {
	var name string
	_ = tx.Raw(`SELECT display_name FROM projections.users WHERE id = ?`, zitadelID).
		Row().Scan(&name)
	return name
}

//nolint:unused // used by membership projectors added in a later task
func lookupOrgName(tx *gorm.DB, orgID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.organizations WHERE id = ?`, orgID).
		Row().Scan(&name)
	return name
}

//nolint:unused // used by membership projectors added in a later task
func lookupClinicName(tx *gorm.DB, clinicID uuid.UUID) null.String {
	var name null.String
	_ = tx.Raw(`SELECT name FROM projections.clinics WHERE id = ?`, clinicID).
		Row().Scan(&name)
	return name
}

//nolint:unused // used by membership projectors added in a later task
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
