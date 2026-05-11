//go:build integration

package self_query_integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	selfread "github.com/medincident/medincident-backend/internal/service/query/self"
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
)

// codeOf extracts the oops Code string from any error in a joined error tree.
func codeOf(t *testing.T, err error) string {
	t.Helper()
	var oe oops.OopsError
	if !errors.As(err, &oe) {
		t.Fatalf("expected oops error, got %T: %v", err, err)
	}
	code, ok := oe.Code().(string)
	if !ok {
		t.Fatalf("oops.Code() returned non-string: %T %v", oe.Code(), oe.Code())
	}
	return code
}

// seedOrg inserts a minimal organization through the projector.
func seedOrg(t *testing.T, ctx context.Context, now time.Time) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrganizationCreated(tx, id.String(), now, &orgv1.OrganizationCreated{
			Name:         "Org " + id.String()[:8],
			LegalAddress: &orgv1.Address{Text: "ул. Тестовая, д. 1"},
			CreatedAt:    timestamppb.New(now),
		})
	}))
	return id
}

// seedClinic inserts a minimal clinic through the projector.
func seedClinic(t *testing.T, ctx context.Context, orgID uuid.UUID, now time.Time) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicCreated(tx, id.String(), now, &clinicv1.ClinicCreated{
			OrganizationId:  orgID.String(),
			Name:            "Clinic " + id.String()[:8],
			PhysicalAddress: &clinicv1.Address{Text: "ул. Клиническая, д. 2"},
			CreatedAt:       timestamppb.New(now),
		})
	}))
	return id
}

// seedDept inserts a minimal department through the projector.
func seedDept(t *testing.T, ctx context.Context, clinicID uuid.UUID, now time.Time) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DepartmentCreated(tx, id.String(), now, &deptv1.DepartmentCreated{
			ClinicId:  clinicID.String(),
			Name:      "Dept " + id.String()[:8],
			CreatedAt: timestamppb.New(now),
		})
	}))
	return id
}

// seedEmployee hires an employee through the projector and returns the employee ID.
func seedEmployee(t *testing.T, ctx context.Context, zitadelUserID string, orgID, deptID uuid.UUID, now time.Time) uuid.UUID {
	t.Helper()
	empID := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.EmployeeHired(tx, empID.String(), now, &empv1.EmployeeHired{
			ZitadelUserId:  zitadelUserID,
			OrganizationId: orgID.String(),
			DepartmentId:   deptID.String(),
			HiredAt:        timestamppb.New(now),
		})
	}))
	return empID
}

// seedSystemAdmin inserts a system admin row through the projector.
func seedSystemAdmin(t *testing.T, ctx context.Context, zitadelUserID string, now time.Time) {
	t.Helper()
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.SystemAdminGranted(tx, zitadelUserID, now, &sav1.SystemAdminGranted{
			GrantedAt: timestamppb.New(now),
		})
	}))
}

// ── GetMyIdentity ────────────────────────────────────────────────────────────

func TestSelfReader_GetMyIdentity_IsSystemAdmin(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	seedSystemAdmin(t, ctx, "sysadmin-user", now)

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyIdentity(ctx, "sysadmin-user")
	require.NoError(t, err)
	assert.True(t, view.IsSystemAdmin)
}

func TestSelfReader_GetMyIdentity_NotSystemAdmin(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyIdentity(ctx, "regular-user")
	require.NoError(t, err)
	assert.False(t, view.IsSystemAdmin)
}

// ── ListMyOrganizations ──────────────────────────────────────────────────────

func TestSelfReader_ListMyOrganizations_Empty(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	reader := selfread.NewSelfReader(testDB, &logger)
	list, err := reader.ListMyOrganizations(ctx, "unknown-user")
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestSelfReader_ListMyOrganizations_OneOrg(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	seedEmployee(t, ctx, "user-1", orgID, deptID, now)

	reader := selfread.NewSelfReader(testDB, &logger)
	list, err := reader.ListMyOrganizations(ctx, "user-1")
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, orgID.String(), list[0].ID)
}

func TestSelfReader_ListMyOrganizations_MultipleOrgs(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID1 := seedOrg(t, ctx, now)
	clinicID1 := seedClinic(t, ctx, orgID1, now)
	deptID1 := seedDept(t, ctx, clinicID1, now)

	orgID2 := seedOrg(t, ctx, now)
	clinicID2 := seedClinic(t, ctx, orgID2, now)
	deptID2 := seedDept(t, ctx, clinicID2, now)

	seedEmployee(t, ctx, "multi-user", orgID1, deptID1, now)
	seedEmployee(t, ctx, "multi-user", orgID2, deptID2, now)

	reader := selfread.NewSelfReader(testDB, &logger)
	list, err := reader.ListMyOrganizations(ctx, "multi-user")
	require.NoError(t, err)
	require.Len(t, list, 2)

	ids := []string{list[0].ID, list[1].ID}
	assert.Contains(t, ids, orgID1.String())
	assert.Contains(t, ids, orgID2.String())
}

func TestSelfReader_ListMyOrganizations_TerminatedExcluded(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	empID := seedEmployee(t, ctx, "terminated-user", orgID, deptID, now)

	terminatedAt := now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.EmployeeTerminated(tx, empID.String(), terminatedAt, &empv1.EmployeeTerminated{
			OrganizationId: orgID.String(),
			DepartmentId:   deptID.String(),
			TerminatedAt:   timestamppb.New(terminatedAt),
		})
	}))

	reader := selfread.NewSelfReader(testDB, &logger)
	list, err := reader.ListMyOrganizations(ctx, "terminated-user")
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ── GetMyEmployment ──────────────────────────────────────────────────────────

func TestSelfReader_GetMyEmployment_Found(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	empID := seedEmployee(t, ctx, "emp-user", orgID, deptID, now)

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyEmployment(ctx, "emp-user", orgID)
	require.NoError(t, err)
	assert.Equal(t, empID, view.EmployeeID)
	assert.Equal(t, "emp-user", view.ZitadelUserID)
	assert.Equal(t, orgID, view.OrganizationID)
	assert.Equal(t, deptID, view.DepartmentID)
}

func TestSelfReader_GetMyEmployment_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	reader := selfread.NewSelfReader(testDB, &logger)
	_, err := reader.GetMyEmployment(ctx, "ghost-user", uuid.Must(uuid.NewV7()))
	require.Error(t, err)
	assert.Equal(t, selfread.ErrCodeSelfEmploymentNotFound, codeOf(t, err))
}

func TestSelfReader_GetMyEmployment_WrongOrg(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	seedEmployee(t, ctx, "user-wrong-org", orgID, deptID, now)

	otherOrgID := uuid.Must(uuid.NewV7())

	reader := selfread.NewSelfReader(testDB, &logger)
	_, err := reader.GetMyEmployment(ctx, "user-wrong-org", otherOrgID)
	require.Error(t, err)
	assert.Equal(t, selfread.ErrCodeSelfEmploymentNotFound, codeOf(t, err))
}

// ── GetMyOrganizationRole ────────────────────────────────────────────────────

func TestSelfReader_GetMyOrganizationRole_NoRoles(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	seedEmployee(t, ctx, "plain-user", orgID, deptID, now)

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyOrganizationRole(ctx, "plain-user", orgID)
	require.NoError(t, err)
	assert.False(t, view.IsOrgAdmin)
	assert.False(t, view.IsOrgHead)
	assert.False(t, view.IsOrgDispatcher)
}

func TestSelfReader_GetMyOrganizationRole_WithRoles(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	empID := seedEmployee(t, ctx, "admin-user", orgID, deptID, now)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgAdminAssigned(tx, orgID.String(), now, &orgv1.OrgAdminAssigned{
			EmployeeId: empID.String(),
			AssignedAt: timestamppb.New(now),
		})
	}))

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyOrganizationRole(ctx, "admin-user", orgID)
	require.NoError(t, err)
	assert.True(t, view.IsOrgAdmin)
	assert.False(t, view.IsOrgHead)
	assert.False(t, view.IsOrgDispatcher)
}

func TestSelfReader_GetMyOrganizationRole_RolesScopedToOrg(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	// Org A: user is org admin.
	orgAID := seedOrg(t, ctx, now)
	clinicAID := seedClinic(t, ctx, orgAID, now)
	deptAID := seedDept(t, ctx, clinicAID, now)
	empAID := seedEmployee(t, ctx, "cross-org-user", orgAID, deptAID, now)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.OrgAdminAssigned(tx, orgAID.String(), now, &orgv1.OrgAdminAssigned{
			EmployeeId: empAID.String(),
			AssignedAt: timestamppb.New(now),
		})
	}))

	// Org B: user is a plain employee, no roles.
	orgBID := seedOrg(t, ctx, now)
	clinicBID := seedClinic(t, ctx, orgBID, now)
	deptBID := seedDept(t, ctx, clinicBID, now)
	seedEmployee(t, ctx, "cross-org-user", orgBID, deptBID, now)

	reader := selfread.NewSelfReader(testDB, &logger)

	// Role in org A is visible when querying org A.
	viewA, err := reader.GetMyOrganizationRole(ctx, "cross-org-user", orgAID)
	require.NoError(t, err)
	assert.True(t, viewA.IsOrgAdmin, "expected is_org_admin for org A")

	// Role must NOT bleed into org B.
	viewB, err := reader.GetMyOrganizationRole(ctx, "cross-org-user", orgBID)
	require.NoError(t, err)
	assert.False(t, viewB.IsOrgAdmin, "is_org_admin must be false for org B")
	assert.False(t, viewB.IsOrgHead)
	assert.False(t, viewB.IsOrgDispatcher)
}

func TestSelfReader_GetMyOrganizationRole_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	reader := selfread.NewSelfReader(testDB, &logger)
	_, err := reader.GetMyOrganizationRole(ctx, "ghost-user", uuid.Must(uuid.NewV7()))
	require.Error(t, err)
	assert.Equal(t, selfread.ErrCodeSelfOrgRoleNotFound, codeOf(t, err))
}

// ── GetMyClinicRole ──────────────────────────────────────────────────────────

func TestSelfReader_GetMyClinicRole_NotHead(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	seedEmployee(t, ctx, "regular-clinic-user", orgID, deptID, now)

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyClinicRole(ctx, "regular-clinic-user", clinicID)
	require.NoError(t, err)
	assert.False(t, view.IsClinicHead)
}

func TestSelfReader_GetMyClinicRole_IsHead(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	empID := seedEmployee(t, ctx, "head-user", orgID, deptID, now)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.ClinicHeadAssigned(tx, clinicID.String(), now, &clinicv1.ClinicHeadAssigned{
			EmployeeId: empID.String(),
			AssignedAt: timestamppb.New(now),
		})
	}))

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyClinicRole(ctx, "head-user", clinicID)
	require.NoError(t, err)
	assert.True(t, view.IsClinicHead)
}

func TestSelfReader_GetMyClinicRole_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	reader := selfread.NewSelfReader(testDB, &logger)
	_, err := reader.GetMyClinicRole(ctx, "ghost-user", uuid.Must(uuid.NewV7()))
	require.Error(t, err)
	assert.Equal(t, selfread.ErrCodeSelfClinicRoleNotFound, codeOf(t, err))
}

// ── GetMyDepartmentRole ──────────────────────────────────────────────────────

func TestSelfReader_GetMyDepartmentRole_NotResponsible(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	seedEmployee(t, ctx, "plain-dept-user", orgID, deptID, now)

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyDepartmentRole(ctx, "plain-dept-user", deptID)
	require.NoError(t, err)
	assert.False(t, view.IsDepartmentResponsible)
}

func TestSelfReader_GetMyDepartmentRole_IsResponsible(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()
	now := time.Now().UTC().Truncate(time.Second)

	orgID := seedOrg(t, ctx, now)
	clinicID := seedClinic(t, ctx, orgID, now)
	deptID := seedDept(t, ctx, clinicID, now)
	empID := seedEmployee(t, ctx, "responsible-user", orgID, deptID, now)

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.DeptResponsibleAssigned(tx, deptID.String(), now, &deptv1.DeptResponsibleAssigned{
			EmployeeId: empID.String(),
			AssignedAt: timestamppb.New(now),
		})
	}))

	reader := selfread.NewSelfReader(testDB, &logger)
	view, err := reader.GetMyDepartmentRole(ctx, "responsible-user", deptID)
	require.NoError(t, err)
	assert.True(t, view.IsDepartmentResponsible)
}

func TestSelfReader_GetMyDepartmentRole_NotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	reader := selfread.NewSelfReader(testDB, &logger)
	_, err := reader.GetMyDepartmentRole(ctx, "ghost-user", uuid.Must(uuid.NewV7()))
	require.Error(t, err)
	assert.Equal(t, selfread.ErrCodeSelfDeptRoleNotFound, codeOf(t, err))
}
