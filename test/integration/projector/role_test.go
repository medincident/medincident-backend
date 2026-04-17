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
		return projector.ClinicHeadAssigned(tx, role)
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		clinic.ID, empID,
	).Row().Scan(&count))
	require.Equal(t, int64(1), count)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.ClinicHeadRevoked(tx, clinic.ID, empID)
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
		return projector.DepartmentResponsibleAssigned(tx, role)
	}))
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.DepartmentResponsibleRevoked(tx, dept.ID, empID)
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
		return projector.OrgAdminAssigned(tx, role)
	}))

	role.DeputyEmployeeID = null.ValueFrom(deputyID)
	role.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrgAdminDeputyAssigned(tx, role)
	}))

	var gotDeputy *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT deputy_employee_id FROM projections.org_admins WHERE organization_id = ? AND employee_id = ?`,
		org.ID, empID,
	).Row().Scan(&gotDeputy))
	require.NotNil(t, gotDeputy)
	require.Equal(t, deputyID, *gotDeputy)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrgAdminDeputyRemoved(tx, org.ID, empID, now.Add(2*time.Hour))
	}))

	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT deputy_employee_id FROM projections.org_admins WHERE organization_id = ? AND employee_id = ?`,
		org.ID, empID,
	).Row().Scan(&gotDeputy))
	require.Nil(t, gotDeputy)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrgAdminRevoked(tx, org.ID, empID)
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
		return projector.OrgDispatcherAssigned(tx, role)
	}))
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrgDispatcherRevoked(tx, org.ID, empID)
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
		return projector.OrgHeadAssigned(tx, role)
	}))
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.OrgHeadRevoked(tx, org.ID, empID)
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
		return projector.SystemAdminGranted(tx, zitadelID, now)
	}))

	var count int64
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.system_admins WHERE zitadel_user_id = ?`, zitadelID,
	).Row().Scan(&count))
	require.Equal(t, int64(1), count)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.SystemAdminRevoked(tx, zitadelID)
	}))

	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.system_admins WHERE zitadel_user_id = ?`, zitadelID,
	).Row().Scan(&count))
	require.Equal(t, int64(0), count)
}
