//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

// seedDepartment is a test helper that creates an organization + clinic
// + department via the projector so employee tests can attach a parent
// dept.
func seedDepartment(t *testing.T, ctx context.Context, now time.Time) (*model.Organization, *model.Clinic, *model.Department) {
	t.Helper()
	org, clinic := seedClinic(t, ctx, now)
	dept := &model.Department{
		ID:        uuid.Must(uuid.NewV7()),
		ClinicID:  clinic.ID,
		Name:      "Parent Department",
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.DepartmentCreated(tx, dept)
	}))
	return org, clinic, dept
}

// TestEmployeeHired_WritesEmployeeCardAndBumpsCounters confirms the
// projector inserts rows into employees and employee_cards, and bumps
// all three employee counter columns.
func TestEmployeeHired_WritesEmployeeCardAndBumpsCounters(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, clinic, dept := seedDepartment(t, ctx, now)

	emp := &model.Employee{
		ID:             uuid.Must(uuid.NewV7()),
		ZitadelUserID:  "zit-user-1",
		OrganizationID: org.ID,
		DepartmentID:   dept.ID,
		Position:       null.StringFrom("Nurse"),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeHired(tx, emp)
	}))

	var zit, position string
	var clinicID *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT zitadel_user_id, clinic_id, position FROM projections.employees WHERE id = ?`, emp.ID).
		Row().Scan(&zit, &clinicID, &position))
	require.Equal(t, "zit-user-1", zit)
	require.NotNil(t, clinicID)
	require.Equal(t, clinic.ID, *clinicID)
	require.Equal(t, "Nurse", position)

	var cardZit string
	var cardOrgName *string
	var cardDeptName *string
	require.NoError(t, testDB.WithContext(ctx).Raw(`
		SELECT zitadel_user_id, organization_name, department_name
		FROM projections.employee_cards WHERE employee_id = ?`, emp.ID).
		Row().Scan(&cardZit, &cardOrgName, &cardDeptName))
	require.Equal(t, "zit-user-1", cardZit)
	require.NotNil(t, cardOrgName)
	require.Equal(t, org.Name, *cardOrgName)
	require.NotNil(t, cardDeptName)
	require.Equal(t, dept.Name, *cardDeptName)

	var orgEmps, clinicEmps, deptEmps int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.organization_counters WHERE organization_id = ?`, org.ID,
	).Row().Scan(&orgEmps))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.clinic_counters WHERE clinic_id = ?`, clinic.ID,
	).Row().Scan(&clinicEmps))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.department_counters WHERE department_id = ?`, dept.ID,
	).Row().Scan(&deptEmps))
	require.Equal(t, int64(1), orgEmps)
	require.Equal(t, int64(1), clinicEmps)
	require.Equal(t, int64(1), deptEmps)
}

// TestEmployeeTerminated_StampsTerminatedAtAndDecrementsCounters
// confirms the projector soft-deletes the employee/card and moves all
// three counters back to zero.
func TestEmployeeTerminated_StampsTerminatedAtAndDecrementsCounters(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, clinic, dept := seedDepartment(t, ctx, now)

	emp := &model.Employee{
		ID:             uuid.Must(uuid.NewV7()),
		ZitadelUserID:  "zit-user-2",
		OrganizationID: org.ID,
		DepartmentID:   dept.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeHired(tx, emp)
	}))

	terminatedAt := now.Add(24 * time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeTerminated(tx, emp, terminatedAt)
	}))

	var empTerminated, cardTerminated *time.Time
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT terminated_at FROM projections.employees WHERE id = ?`, emp.ID,
	).Row().Scan(&empTerminated))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT terminated_at FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&cardTerminated))
	require.NotNil(t, empTerminated)
	require.NotNil(t, cardTerminated)

	var orgEmps, clinicEmps, deptEmps int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.organization_counters WHERE organization_id = ?`, org.ID,
	).Row().Scan(&orgEmps))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.clinic_counters WHERE clinic_id = ?`, clinic.ID,
	).Row().Scan(&clinicEmps))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.department_counters WHERE department_id = ?`, dept.ID,
	).Row().Scan(&deptEmps))
	require.Equal(t, int64(0), orgEmps)
	require.Equal(t, int64(0), clinicEmps)
	require.Equal(t, int64(0), deptEmps)
}

// TestEmployeeDepartmentChanged_CrossClinicMove confirms counters move
// between clinics and departments when the target department lives in
// a different clinic.
func TestEmployeeDepartmentChanged_CrossClinicMove(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, clinic, dept := seedDepartment(t, ctx, now)

	// Create a second clinic + department in the same org for the move target.
	clinic2 := &model.Clinic{
		ID:              uuid.Must(uuid.NewV7()),
		OrganizationID:  org.ID,
		Name:            "Clinic B",
		PhysicalAddress: model.Address{Text: "99 Other Street"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.ClinicCreated(tx, clinic2)
	}))
	dept2 := &model.Department{
		ID:        uuid.Must(uuid.NewV7()),
		ClinicID:  clinic2.ID,
		Name:      "Other Dept",
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.DepartmentCreated(tx, dept2)
	}))

	emp := &model.Employee{
		ID:             uuid.Must(uuid.NewV7()),
		ZitadelUserID:  "zit-user-3",
		OrganizationID: org.ID,
		DepartmentID:   dept.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeHired(tx, emp)
	}))

	// Move employee to dept2 (different clinic).
	emp.DepartmentID = dept2.ID
	emp.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeDepartmentChanged(tx, emp)
	}))

	var c1, c2, d1, d2 int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.clinic_counters WHERE clinic_id = ?`, clinic.ID,
	).Row().Scan(&c1))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.clinic_counters WHERE clinic_id = ?`, clinic2.ID,
	).Row().Scan(&c2))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.department_counters WHERE department_id = ?`, dept.ID,
	).Row().Scan(&d1))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT employees_total FROM projections.department_counters WHERE department_id = ?`, dept2.ID,
	).Row().Scan(&d2))
	require.Equal(t, int64(0), c1)
	require.Equal(t, int64(1), c2)
	require.Equal(t, int64(0), d1)
	require.Equal(t, int64(1), d2)

	var newClinicID *uuid.UUID
	var newDeptName *string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT clinic_id, department_name FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&newClinicID, &newDeptName))
	require.NotNil(t, newClinicID)
	require.Equal(t, clinic2.ID, *newClinicID)
	require.NotNil(t, newDeptName)
	require.Equal(t, dept2.Name, *newDeptName)
}

// TestEmployeePositionChanged_UpdatesRowAndCard confirms the projector
// updates position on both employees and employee_cards.
func TestEmployeePositionChanged_UpdatesRowAndCard(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, _, dept := seedDepartment(t, ctx, now)

	emp := &model.Employee{
		ID:             uuid.Must(uuid.NewV7()),
		ZitadelUserID:  "zit-user-4",
		OrganizationID: org.ID,
		DepartmentID:   dept.ID,
		Position:       null.StringFrom("Nurse"),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeHired(tx, emp)
	}))

	emp.Position = null.StringFrom("Senior Nurse")
	emp.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeePositionChanged(tx, emp)
	}))

	var empPos, cardPos *string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT position FROM projections.employees WHERE id = ?`, emp.ID,
	).Row().Scan(&empPos))
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT position FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&cardPos))
	require.NotNil(t, empPos)
	require.Equal(t, "Senior Nurse", *empPos)
	require.NotNil(t, cardPos)
	require.Equal(t, "Senior Nurse", *cardPos)
}
