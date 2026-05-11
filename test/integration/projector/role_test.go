//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
)

// TestClinicHead_AssignAndRevoke verifies the assign/revoke pair for
// clinic_head projections.
func TestClinicHead_AssignAndRevoke(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org, clinic := seedClinic(t, ctx, now)
	_ = org
	empID := uuid.Must(uuid.NewV7())

	role := &model.ClinicHead{
		ClinicID:   clinic.ID,
		EmployeeID: empID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicHeadAssigned(tx, role.ClinicID.String(), role.CreatedAt, &clinicv1.ClinicHeadAssigned{
			EmployeeId: role.EmployeeID.String(),
			AssignedAt: timestamppb.New(role.CreatedAt),
		})
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		clinic.ID, empID,
	).Row().Scan(&count))
	require.Equal(t, int64(1), count)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicHeadRevoked(tx, clinic.ID.String(), time.Now(), &clinicv1.ClinicHeadRevoked{
			EmployeeId: empID.String(),
		})
	}))

	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		clinic.ID, empID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}

// TestDepartmentResponsible_AssignAndRevoke verifies department_responsible
// projection writes.
func TestDepartmentResponsible_AssignAndRevoke(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	_, _, dept := seedDepartment(t, ctx, now)
	empID := uuid.Must(uuid.NewV7())

	role := &model.DepartmentResponsible{
		DepartmentID: dept.ID,
		EmployeeID:   empID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DeptResponsibleAssigned(tx, role.DepartmentID.String(), role.CreatedAt, &deptv1.DeptResponsibleAssigned{
			EmployeeId: role.EmployeeID.String(),
			AssignedAt: timestamppb.New(role.CreatedAt),
		})
	}))
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DeptResponsibleRevoked(tx, dept.ID.String(), time.Now(), &deptv1.DeptResponsibleRevoked{
			EmployeeId: empID.String(),
		})
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.department_responsibles WHERE department_id = ?`, dept.ID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}

// TestOrgAdmin_AssignWithDeputyAndRevoke verifies org_admin projection,
// including the deputy slot.
func TestOrgAdmin_AssignWithDeputyAndRevoke(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org := seedOrganization(t, ctx, now)
	empID := uuid.Must(uuid.NewV7())
	deputyID := uuid.Must(uuid.NewV7())

	role := &model.OrgAdmin{
		OrganizationID: org.ID,
		EmployeeID:     empID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgAdminAssigned(tx, role.OrganizationID.String(), role.CreatedAt, &orgv1.OrgAdminAssigned{
			EmployeeId: role.EmployeeID.String(),
			AssignedAt: timestamppb.New(role.CreatedAt),
		})
	}))

	role.DeputyEmployeeID = null.ValueFrom(deputyID)
	role.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgAdminDeputyAssigned(tx, role.OrganizationID.String(), role.UpdatedAt, &orgv1.OrgAdminDeputyAssigned{
			EmployeeId:       role.EmployeeID.String(),
			DeputyEmployeeId: role.DeputyEmployeeID.V.String(),
			UpdatedAt:        timestamppb.New(role.UpdatedAt),
		})
	}))

	var gotDeputy *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT deputy_employee_id FROM projections.org_admins WHERE organization_id = ? AND employee_id = ?`,
		org.ID, empID,
	).Row().Scan(&gotDeputy))
	require.NotNil(t, gotDeputy)
	require.Equal(t, deputyID, *gotDeputy)

	removedAt := now.Add(2 * time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgAdminDeputyRemoved(tx, org.ID.String(), removedAt, &orgv1.OrgAdminDeputyRemoved{
			EmployeeId: empID.String(),
			UpdatedAt:  timestamppb.New(removedAt),
		})
	}))

	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT deputy_employee_id FROM projections.org_admins WHERE organization_id = ? AND employee_id = ?`,
		org.ID, empID,
	).Row().Scan(&gotDeputy))
	require.Nil(t, gotDeputy)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgAdminRevoked(tx, org.ID.String(), time.Now(), &orgv1.OrgAdminRevoked{
			EmployeeId: empID.String(),
		})
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.org_admins WHERE organization_id = ?`, org.ID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}

// TestOrgDispatcher_AssignAndRevoke verifies org_dispatcher projection writes.
func TestOrgDispatcher_AssignAndRevoke(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org := seedOrganization(t, ctx, now)
	empID := uuid.Must(uuid.NewV7())

	role := &model.OrgDispatcher{
		OrganizationID: org.ID,
		EmployeeID:     empID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgDispatcherAssigned(tx, role.OrganizationID.String(), role.CreatedAt, &orgv1.OrgDispatcherAssigned{
			EmployeeId: role.EmployeeID.String(),
			AssignedAt: timestamppb.New(role.CreatedAt),
		})
	}))
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgDispatcherRevoked(tx, org.ID.String(), time.Now(), &orgv1.OrgDispatcherRevoked{
			EmployeeId: empID.String(),
		})
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.org_dispatchers WHERE organization_id = ?`, org.ID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}

// TestOrgHead_AssignAndRevoke verifies org_head projection writes.
func TestOrgHead_AssignAndRevoke(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	org := seedOrganization(t, ctx, now)
	empID := uuid.Must(uuid.NewV7())

	role := &model.OrgHead{
		OrganizationID: org.ID,
		EmployeeID:     empID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgHeadAssigned(tx, role.OrganizationID.String(), role.CreatedAt, &orgv1.OrgHeadAssigned{
			EmployeeId: role.EmployeeID.String(),
			AssignedAt: timestamppb.New(role.CreatedAt),
		})
	}))
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgHeadRevoked(tx, org.ID.String(), time.Now(), &orgv1.OrgHeadRevoked{
			EmployeeId: empID.String(),
		})
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.org_heads WHERE organization_id = ?`, org.ID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}

// TestSystemAdmin_GrantAndRevoke verifies the system_admin projection
// pair — simplest of the role projections (no deputy, no scope).
func TestSystemAdmin_GrantAndRevoke(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	zitadelID := "zit-admin-1"

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.SystemAdminGranted(tx, zitadelID, now, &sav1.SystemAdminGranted{
			GrantedAt: timestamppb.New(now),
		})
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.system_admins WHERE zitadel_user_id = ?`, zitadelID,
	).Row().Scan(&count))
	require.Equal(t, int64(1), count)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.SystemAdminRevoked(tx, zitadelID, time.Now(), &sav1.SystemAdminRevoked{})
	}))

	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.system_admins WHERE zitadel_user_id = ?`, zitadelID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}
