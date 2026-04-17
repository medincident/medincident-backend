//go:build integration

package authz_integration_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/authz"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func requirePermissionDenied(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe), "expected oops error, got: %v", err)
	assert.Equal(t, authz.ErrCodePermissionDenied, oe.Code())
}

func requireScopeResolveFailed(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe), "expected oops error, got: %v", err)
	assert.Equal(t, authz.ErrCodeScopeResolveFailed, oe.Code())
}

// ---------------------------------------------------------------------------
// RequireSystemAdmin
// ---------------------------------------------------------------------------

func TestRequireSystemAdmin_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireSystemAdmin(ctxT(t), "sysadmin")
	require.NoError(t, err)
}

func TestRequireSystemAdmin_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireSystemAdmin(ctxT(t), "alice")
	requirePermissionDenied(t, err)
}

func TestRequireSystemAdmin_UnknownCaller(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireSystemAdmin(ctxT(t), "nonexistent")
	requirePermissionDenied(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdmin
// ---------------------------------------------------------------------------

func TestRequireOrgAdmin_SystemAdmin_AnyOrg(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	require.NoError(t, authzSvc.RequireOrgAdmin(ctxT(t), "sysadmin", orgA))
	require.NoError(t, authzSvc.RequireOrgAdmin(ctxT(t), "sysadmin", orgB))
}

func TestRequireOrgAdmin_DirectAdmin(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdmin(ctxT(t), "bob", orgA)
	require.NoError(t, err)
}

func TestRequireOrgAdmin_CrossOrgDenied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdmin(ctxT(t), "bob", orgB)
	requirePermissionDenied(t, err)
}

func TestRequireOrgAdmin_DeputyOnVacation(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	// Bob is on vacation (active), Carol is his deputy -> allowed.
	err := authzSvc.RequireOrgAdmin(ctxT(t), "carol", orgA)
	require.NoError(t, err)
}

func TestRequireOrgAdmin_DeputyNoVacation(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	// End Bob's vacation so Carol no longer has deputy access.
	require.NoError(t, testDB.Exec(
		`UPDATE domain.employee_vacations SET ends_at = now() - interval '1 minute' WHERE id = ?`,
		bobVacationID,
	).Error)

	err := authzSvc.RequireOrgAdmin(ctxT(t), "carol", orgA)
	requirePermissionDenied(t, err)
}

func TestRequireOrgAdmin_RegularEmployee(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdmin(ctxT(t), "alice", orgA)
	requirePermissionDenied(t, err)
}

func TestRequireOrgAdmin_UnknownCaller(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdmin(ctxT(t), "ghost", orgA)
	requirePermissionDenied(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaClinic
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaClinic_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaClinic(ctxT(t), "sysadmin", clinicA1)
	require.NoError(t, err)
}

func TestRequireOrgAdminViaClinic_CrossOrg(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	// Bob is OrgA admin, clinicB1 belongs to OrgB -> denied.
	err := authzSvc.RequireOrgAdminViaClinic(ctxT(t), "bob", clinicB1)
	requirePermissionDenied(t, err)
}

func TestRequireOrgAdminViaClinic_NotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaClinic(ctxT(t), "sysadmin", uuid.Must(uuid.NewV7()))
	requireScopeResolveFailed(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaDepartment
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaDepartment_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaDepartment(ctxT(t), "sysadmin", deptA1a)
	require.NoError(t, err)
}

func TestRequireOrgAdminViaDepartment_NotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaDepartment(ctxT(t), "sysadmin", uuid.Must(uuid.NewV7()))
	requireScopeResolveFailed(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaEmployee
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaEmployee_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaEmployee(ctxT(t), "sysadmin", empAlice)
	require.NoError(t, err)
}

func TestRequireOrgAdminViaEmployee_NotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaEmployee(ctxT(t), "sysadmin", uuid.Must(uuid.NewV7()))
	requireScopeResolveFailed(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaCategory
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaCategory_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaCategory(ctxT(t), "sysadmin", categoryA)
	require.NoError(t, err)
}

func TestRequireOrgAdminViaCategory_NotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaCategory(ctxT(t), "sysadmin", uuid.Must(uuid.NewV7()))
	requireScopeResolveFailed(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaType
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaType_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaType(ctxT(t), "sysadmin", typeA)
	require.NoError(t, err)
}

func TestRequireOrgAdminViaType_NotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaType(ctxT(t), "sysadmin", uuid.Must(uuid.NewV7()))
	requireScopeResolveFailed(t, err)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaVacation
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaVacation_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaVacation(ctxT(t), "sysadmin", aliceVacID)
	require.NoError(t, err)
}

func TestRequireOrgAdminViaVacation_NotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.RequireOrgAdminViaVacation(ctxT(t), "sysadmin", uuid.Must(uuid.NewV7()))
	requireScopeResolveFailed(t, err)
}
