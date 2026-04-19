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

func requireDenied(t *testing.T, err error, wantPublic string) {
	t.Helper()
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe), "expected oops error, got: %v", err)
	assert.Equal(t, authz.ErrCodePermissionDenied, oe.Code())
	assert.Equal(t, wantPublic, oe.Public())
}

const (
	publicSystemAdmin = "Access denied: requires system administrator privileges."
	publicOrgAdmin    = "Access denied: requires organization administrator privileges."
	publicAdminOf     = "Access denied: requires system administrator or organization administrator privileges."
)

func TestRequire_SystemAdmin_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	require.NoError(t, authzSvc.Require(ctxT(t), "sysadmin", authz.SystemAdmin))
}

func TestRequire_SystemAdmin_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	err := authzSvc.Require(ctxT(t), "alice", authz.SystemAdmin)
	requireDenied(t, err, publicSystemAdmin)
}

func TestRequire_OrgAdminOf_Organization_DirectAdmin(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	require.NoError(t, authzSvc.Require(ctxT(t), "bob", authz.OrgAdminOf.Organization(orgA)))
}

func TestRequire_OrgAdminOf_Organization_NonAdmin_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	err := authzSvc.Require(ctxT(t), "alice", authz.OrgAdminOf.Organization(orgA))
	requireDenied(t, err, publicOrgAdmin)
}

func TestRequire_OrgAdminOf_Organization_DeputyOnVacation(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	// Bob is on vacation, Carol is his deputy — authorized.
	require.NoError(t, authzSvc.Require(ctxT(t), "carol", authz.OrgAdminOf.Organization(orgA)))
}

func TestRequire_OrgAdminOf_Clinic_CrossOrg_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	// Bob admins OrgA; clinicB1 belongs to OrgB.
	err := authzSvc.Require(ctxT(t), "bob", authz.OrgAdminOf.Clinic(clinicB1))
	requireDenied(t, err, publicOrgAdmin)
}

func TestRequire_OrgAdminOf_Clinic_Unknown_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	err := authzSvc.Require(ctxT(t), "bob", authz.OrgAdminOf.Clinic(uuid.Must(uuid.NewV7())))
	requireDenied(t, err, publicOrgAdmin)
}

func TestRequire_AdminOf_Clinic_SysAdminBypass(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	// sysadmin passes even for a clinic they don't directly admin.
	require.NoError(t, authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Clinic(clinicA1)))
	// and even for a non-existent clinic (service layer is the authoritative not-found source for sysadmins).
	require.NoError(t, authzSvc.Require(ctxT(t), "sysadmin", authz.AdminOf.Clinic(uuid.Must(uuid.NewV7()))))
}

func TestRequire_AdminOf_Clinic_OrgAdmin_Allowed(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	require.NoError(t, authzSvc.Require(ctxT(t), "bob", authz.AdminOf.Clinic(clinicA1)))
}

func TestRequire_AdminOf_Clinic_Unrelated_Denied(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	err := authzSvc.Require(ctxT(t), "alice", authz.AdminOf.Clinic(clinicA1))
	requireDenied(t, err, publicAdminOf)
}

func TestRequire_AnyOf_ExplicitComposition_EquivalentToAdminOf(t *testing.T) {
	resetDB(t)
	seedFixtures(t)
	explicit := authz.AnyOf(authz.SystemAdmin, authz.OrgAdminOf.Clinic(clinicA1))
	require.NoError(t, authzSvc.Require(ctxT(t), "bob", explicit))
	require.NoError(t, authzSvc.Require(ctxT(t), "sysadmin", explicit))
	err := authzSvc.Require(ctxT(t), "alice", explicit)
	requireDenied(t, err, publicAdminOf)
}
