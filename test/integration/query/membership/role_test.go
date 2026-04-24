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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
)

// TestRoleReader_GetClinicHead returns ErrRoleVacant when the clinic
// has no head, an assignment view with the joined holder card when
// assigned.
func TestRoleReader_GetClinicHead(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID, clinicID, deptID := seedOrgClinicDept(t, ctx, now)

	reader := memberread.NewRoleReader(testDB, authzSvc, &logger)

	view, err := reader.GetClinicHead(ctx, sysadminCaller, clinicID)
	require.ErrorIs(t, err, memberread.ErrRoleVacant)
	require.Nil(t, view)

	// Seed the holder employee first — the JOIN requires the
	// employee_cards row to exist.
	empID := uuid.Must(uuid.NewV7())
	emp := &model.Employee{
		ID: empID, ZitadelUserID: "zit-head",
		OrganizationID: orgID, DepartmentID: deptID,
		CreatedAt: now, UpdatedAt: now,
	}
	head := &model.ClinicHead{
		ClinicID: clinicID, EmployeeID: empID,
		CreatedAt: now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.EmployeeHired(tx, emp); err != nil {
			return err
		}
		return projector.ClinicHeadAssigned(tx, head)
	}))

	view, err = reader.GetClinicHead(ctx, sysadminCaller, clinicID)
	require.NoError(t, err)
	require.NotNil(t, view)
	require.Equal(t, empID, view.Holder.EmployeeID)
	require.Equal(t, "zit-head", view.Holder.ZitadelUserID)
	require.Nil(t, view.Deputy)
}

// TestRoleReader_GetClinicHead_WithDeputy seeds a clinic head with a
// deputy plus an active vacation on the main holder and asserts both
// .Holder and .Deputy come back populated through the JOIN.
func TestRoleReader_GetClinicHead_WithDeputy(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	orgID, clinicID, deptID := seedOrgClinicDept(t, ctx, now)

	holderID := uuid.Must(uuid.NewV7())
	deputyID := uuid.Must(uuid.NewV7())
	holder := &model.Employee{
		ID: holderID, ZitadelUserID: "zit-holder",
		OrganizationID: orgID, DepartmentID: deptID,
		CreatedAt: now, UpdatedAt: now,
	}
	deputy := &model.Employee{
		ID: deputyID, ZitadelUserID: "zit-deputy",
		OrganizationID: orgID, DepartmentID: deptID,
		CreatedAt: now, UpdatedAt: now,
	}
	vac := &model.EmployeeVacation{
		ID: uuid.Must(uuid.NewV7()), EmployeeID: holderID,
		StartsAt: now.Add(-time.Hour), EndsAt: null.TimeFrom(now.Add(72 * time.Hour)),
		CreatedAt: now, UpdatedAt: now,
	}
	head := &model.ClinicHead{
		ClinicID:         clinicID,
		EmployeeID:       holderID,
		DeputyEmployeeID: null.ValueFrom(deputyID),
		CreatedAt:        now, UpdatedAt: now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.EmployeeHired(tx, holder); err != nil {
			return err
		}
		if err := projector.EmployeeHired(tx, deputy); err != nil {
			return err
		}
		if err := projector.VacationStarted(tx, vac); err != nil {
			return err
		}
		return projector.ClinicHeadAssigned(tx, head)
	}))

	reader := memberread.NewRoleReader(testDB, authzSvc, &logger)
	view, err := reader.GetClinicHead(ctx, sysadminCaller, clinicID)
	require.NoError(t, err)
	require.NotNil(t, view)

	require.Equal(t, holderID, view.Holder.EmployeeID)
	require.Equal(t, "zit-holder", view.Holder.ZitadelUserID)
	require.NotNil(t, view.Holder.CurrentVacationEndsAt,
		"holder's active vacation should surface on the joined card")

	require.NotNil(t, view.Deputy)
	require.Equal(t, deputyID, view.Deputy.EmployeeID)
	require.Equal(t, "zit-deputy", view.Deputy.ZitadelUserID)
}

// TestRoleReader_ListSystemAdmins returns every row, newest first.
func TestRoleReader_ListSystemAdmins(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := projector.SystemAdminGranted(tx, "zit-a", now); err != nil {
			return err
		}
		return projector.SystemAdminGranted(tx, "zit-b", now.Add(time.Second))
	}))

	reader := memberread.NewRoleReader(testDB, authzSvc, &logger)
	items, err := reader.ListSystemAdmins(ctx, sysadminCaller, memberread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, "zit-b", items[0].ZitadelUserID)
}
