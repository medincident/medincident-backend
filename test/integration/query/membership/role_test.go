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
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
	vacv1 "github.com/medincident/medincident-backend/pkg/event/vacation/v1"
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
		if err := qprojector.EmployeeHired(tx, emp.ID.String(), emp.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  emp.ZitadelUserID,
			OrganizationId: emp.OrganizationID.String(),
			DepartmentId:   emp.DepartmentID.String(),
			HiredAt:        timestamppb.New(emp.CreatedAt),
		}); err != nil {
			return err
		}
		return qprojector.ClinicHeadAssigned(tx, head.ClinicID.String(), head.CreatedAt, &clinicv1.ClinicHeadAssigned{
			EmployeeId: head.EmployeeID.String(),
			AssignedAt: timestamppb.New(head.CreatedAt),
		})
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
		if err := qprojector.EmployeeHired(tx, holder.ID.String(), holder.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  holder.ZitadelUserID,
			OrganizationId: holder.OrganizationID.String(),
			DepartmentId:   holder.DepartmentID.String(),
			HiredAt:        timestamppb.New(holder.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.EmployeeHired(tx, deputy.ID.String(), deputy.CreatedAt, &empv1.EmployeeHired{
			ZitadelUserId:  deputy.ZitadelUserID,
			OrganizationId: deputy.OrganizationID.String(),
			DepartmentId:   deputy.DepartmentID.String(),
			HiredAt:        timestamppb.New(deputy.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.VacationStarted(tx, vac.EmployeeID.String(), vac.CreatedAt, &vacv1.VacationStarted{
			VacationId: vac.ID.String(),
			StartsAt:   timestamppb.New(vac.StartsAt),
			EndsAt:     timestamppb.New(vac.EndsAt.Time),
			CreatedAt:  timestamppb.New(vac.CreatedAt),
		}); err != nil {
			return err
		}
		if err := qprojector.ClinicHeadAssigned(tx, head.ClinicID.String(), head.CreatedAt, &clinicv1.ClinicHeadAssigned{
			EmployeeId: head.EmployeeID.String(),
			AssignedAt: timestamppb.New(head.CreatedAt),
		}); err != nil {
			return err
		}
		return qprojector.ClinicHeadDeputyAssigned(tx, head.ClinicID.String(), head.UpdatedAt, &clinicv1.ClinicHeadDeputyAssigned{
			EmployeeId:       head.EmployeeID.String(),
			DeputyEmployeeId: head.DeputyEmployeeID.V.String(),
			UpdatedAt:        timestamppb.New(head.UpdatedAt),
		})
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
		if err := qprojector.SystemAdminGranted(tx, "zit-a", now, &sav1.SystemAdminGranted{
			GrantedAt: timestamppb.New(now),
		}); err != nil {
			return err
		}
		return qprojector.SystemAdminGranted(tx, "zit-b", now.Add(time.Second), &sav1.SystemAdminGranted{
			GrantedAt: timestamppb.New(now.Add(time.Second)),
		})
	}))

	reader := memberread.NewRoleReader(testDB, authzSvc, &logger)
	items, err := reader.ListSystemAdmins(ctx, sysadminCaller, memberread.ListQuery{})
	require.NoError(t, err)
	require.Len(t, items.Items, 2)
	require.Equal(t, "zit-b", items.Items[0].ZitadelUserID)
}
