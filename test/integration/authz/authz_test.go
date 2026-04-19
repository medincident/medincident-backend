//go:build integration

package authz_integration_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/authz"
)

// Tests in this file share the requireDenied helper and public-message
// constants with policy_test.go (same package).

// ---------------------------------------------------------------------------
// RequireSystemAdmin
// ---------------------------------------------------------------------------

func TestRequireSystemAdmin_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.SystemAdmin)
	require.NoError(t, err)
}

func TestRequireSystemAdmin_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "alice", authz.SystemAdmin)
	requireDenied(t, err, publicSystemAdmin)
}

func TestRequireSystemAdmin_UnknownCaller(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "nonexistent", authz.SystemAdmin)
	requireDenied(t, err, publicSystemAdmin)
}

// ---------------------------------------------------------------------------
// RequireOrgAdmin
// ---------------------------------------------------------------------------

func TestRequireOrgAdmin_SystemAdmin_AnyOrg(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	require.NoError(t, authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Organization(orgA)))
	require.NoError(t, authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Organization(orgB)))
}

func TestRequireOrgAdmin_DirectAdmin(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Organization(orgA))
	require.NoError(t, err)
}

func TestRequireOrgAdmin_CrossOrgDenied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Organization(orgB))
	requireDenied(t, err, publicAdminOf)
}

func TestRequireOrgAdmin_DeputyOnVacation(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	// Bob is on vacation (active), Carol is his deputy -> allowed.
	err := authzSvc.Require(ctxT(t), "carol", authz.AdminOf.Organization(orgA))
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

	err := authzSvc.Require(ctxT(t), "carol", authz.AdminOf.Organization(orgA))
	requireDenied(t, err, publicAdminOf)
}

func TestRequireOrgAdmin_RegularEmployee(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "alice", authz.AdminOf.Organization(orgA))
	requireDenied(t, err, publicAdminOf)
}

func TestRequireOrgAdmin_UnknownCaller(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "ghost", authz.AdminOf.Organization(orgA))
	requireDenied(t, err, publicAdminOf)
}

// Non-sysadmin caller targeting a random org id must get
// permission_denied — the zero-row EXISTS result is the same as
// "wrong org", so enumeration via differentiated error codes is not
// possible.
func TestRequireOrgAdmin_NonSysAdminNonexistentOrg(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Organization(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaClinic
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaClinic_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Clinic(clinicA1))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaClinic_CrossOrg(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	// Bob is OrgA admin, clinicB1 belongs to OrgB -> denied.
	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Clinic(clinicB1))
	requireDenied(t, err, publicAdminOf)
}

// Sysadmin authorizes regardless of whether the clinic exists — the
// service layer is the authoritative source of not-found for them.
func TestRequireOrgAdminViaClinic_SysAdminNotFoundIsAllowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Clinic(uuid.Must(uuid.NewV7())))
	require.NoError(t, err)
}

// A non-sysadmin caller probing a random clinic id must be
// indistinguishable from a cross-tenant access attempt.
func TestRequireOrgAdminViaClinic_NonSysAdminNotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Clinic(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaDepartment
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaDepartment_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Department(deptA1a))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaDepartment_SysAdminNotFoundIsAllowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Department(uuid.Must(uuid.NewV7())))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaDepartment_NonSysAdminNotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Department(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaEmployee
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaEmployee_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Employee(empAlice))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaEmployee_SysAdminNotFoundIsAllowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Employee(uuid.Must(uuid.NewV7())))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaEmployee_NonSysAdminNotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Employee(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}

// Bob (OrgA admin) targeting Dave (OrgB employee) must be denied.
func TestRequireOrgAdminViaEmployee_CrossOrg(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Employee(empDave))
	requireDenied(t, err, publicAdminOf)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaCategory
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaCategory_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Category(categoryA))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaCategory_SysAdminNotFoundIsAllowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Category(uuid.Must(uuid.NewV7())))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaCategory_NonSysAdminNotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Category(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaType
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaType_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.IncidentType(typeA))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaType_SysAdminNotFoundIsAllowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.IncidentType(uuid.Must(uuid.NewV7())))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaType_NonSysAdminNotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.IncidentType(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}

// ---------------------------------------------------------------------------
// RequireOrgAdminViaVacation
// ---------------------------------------------------------------------------

func TestRequireOrgAdminViaVacation_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Vacation(aliceVacID))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaVacation_SysAdminNotFoundIsAllowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Vacation(uuid.Must(uuid.NewV7())))
	require.NoError(t, err)
}

func TestRequireOrgAdminViaVacation_NonSysAdminNotFound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	err := authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Vacation(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicAdminOf)
}
