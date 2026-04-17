//go:build integration

package membership_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	memberread "github.com/medincident/medincident-command-service/internal/service/query/membership"
)

// seedOrgClinicDept inserts an organization + clinic + department
// through the projector so the employee projection has parents.
func seedOrgClinicDept(t *testing.T, ctx context.Context, now time.Time) (uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	orgID := uuid.Must(uuid.NewV7())
	clinicID := uuid.Must(uuid.NewV7())
	deptID := uuid.Must(uuid.NewV7())
	org := &model.Organization{ID: orgID, Name: "O", LegalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now}
	clinic := &model.Clinic{ID: clinicID, OrganizationID: orgID, Name: "C", PhysicalAddress: model.Address{Text: "S"}, CreatedAt: now, UpdatedAt: now}
	dept := &model.Department{ID: deptID, ClinicID: clinicID, Name: "D", CreatedAt: now, UpdatedAt: now}

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.OrganizationCreated(tx, org); err != nil {
			return err
		}
		if err := projector.ClinicCreated(tx, clinic); err != nil {
			return err
		}
		return projector.DepartmentCreated(tx, dept)
	}))
	return orgID, clinicID, deptID
}

// TestEmployeeReader_Get_AndListByDepartment seeds an employee through
// the projector and reads it back.
func TestEmployeeReader_Get_AndListByDepartment(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID, _, deptID := seedOrgClinicDept(t, ctx, now)

	empID := uuid.Must(uuid.NewV7())
	emp := &model.Employee{
		ID:             empID,
		ZitadelUserID:  "zit-1",
		OrganizationID: orgID,
		DepartmentID:   deptID,
		Position:       null.StringFrom("Nurse"),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeHired(tx, emp)
	}))

	reader := memberread.NewEmployeeReader(testDB, &logger)
	got, err := reader.Get(ctx, empID)
	require.NoError(t, err)
	require.Equal(t, empID, got.EmployeeID)
	require.Equal(t, "zit-1", got.ZitadelUserID)
	require.Equal(t, deptID, got.DepartmentID)
	require.NotNil(t, got.Position)
	require.Equal(t, "Nurse", *got.Position)

	list, err := reader.ListByDepartment(ctx, deptID, memberread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, empID, list[0].EmployeeID)
}

// TestEmployeeReader_Get_NotFound surfaces the expected code.
func TestEmployeeReader_Get_NotFound(t *testing.T) {
	resetProjections(t)
	logger := zerolog.Nop()
	reader := memberread.NewEmployeeReader(testDB, &logger)
	_, err := reader.Get(context.Background(), uuid.Must(uuid.NewV7()))
	require.Error(t, err)
}

// TestEmployeeReader_ListVacationsByEmployee verifies scheduled
// vacations are visible + the state filter works.
func TestEmployeeReader_ListVacationsByEmployee(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID, _, deptID := seedOrgClinicDept(t, ctx, now)

	empID := uuid.Must(uuid.NewV7())
	emp := &model.Employee{
		ID: empID, ZitadelUserID: "zit-x", OrganizationID: orgID, DepartmentID: deptID,
		CreatedAt: now, UpdatedAt: now,
	}
	vacID := uuid.Must(uuid.NewV7())
	vac := &model.EmployeeVacation{
		ID: vacID, EmployeeID: empID,
		StartsAt: now.Add(24 * time.Hour), EndsAt: null.TimeFrom(now.Add(72 * time.Hour)),
		CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.EmployeeHired(tx, emp); err != nil {
			return err
		}
		return projector.VacationScheduled(tx, vac)
	}))

	reader := memberread.NewEmployeeReader(testDB, &logger)
	all, err := reader.ListVacationsByEmployee(ctx, empID, "")
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Equal(t, "scheduled", all[0].State)

	scheduled, err := reader.ListVacationsByEmployee(ctx, empID, "scheduled")
	require.NoError(t, err)
	require.Len(t, scheduled, 1)

	none, err := reader.ListVacationsByEmployee(ctx, empID, "active")
	require.NoError(t, err)
	require.Empty(t, none)
}
