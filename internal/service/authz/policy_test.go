package authz

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchCtx_IndependentCounters(t *testing.T) {
	bc := &branchCtx{callerID: "alice"}

	c1 := bc.addCaller()
	s1 := bc.addScope(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	c2 := bc.addCaller()
	s2 := bc.addScope(uuid.MustParse("22222222-2222-2222-2222-222222222222"))

	// Counters are independent so names are natural sequences per type.
	assert.Equal(t, "caller0", c1)
	assert.Equal(t, "scope0", s1)
	assert.Equal(t, "caller1", c2)
	assert.Equal(t, "scope1", s2)
	require.Len(t, bc.args, 4)

	first, ok := bc.args[0].(sql.NamedArg)
	require.True(t, ok)
	assert.Equal(t, "caller0", first.Name)
	assert.Equal(t, "alice", first.Value)
}

func TestSystemAdmin_Describe(t *testing.T) {
	assert.Equal(t, "system administrator", SystemAdmin.describe())
}

func TestSystemAdmin_Branches_SingleBranch(t *testing.T) {
	bc := &branchCtx{callerID: "bob"}
	bs := SystemAdmin.branches(bc)
	require.Len(t, bs, 1)
	assert.Contains(t, bs[0], "domain.system_admins")
	assert.Contains(t, bs[0], "@caller0")
}

func TestAnyOf_ConcatenatesBranches(t *testing.T) {
	bc := &branchCtx{callerID: "bob"}
	bs := AnyOf(SystemAdmin, SystemAdmin).branches(bc)
	assert.Len(t, bs, 2)
	// Each SystemAdmin branch requests its own caller placeholder.
	assert.Contains(t, bs[0], "@caller0")
	assert.Contains(t, bs[1], "@caller1")
}

func TestAnyOf_Describe_JoinsWithOr(t *testing.T) {
	assert.Equal(t,
		"system administrator or system administrator",
		AnyOf(SystemAdmin, SystemAdmin).describe())
}

func TestOrgAdminOf_Clinic_Branches_DirectAndDeputy(t *testing.T) {
	bc := &branchCtx{callerID: "bob"}
	bs := OrgAdminOf.Clinic(uuid.MustParse("22222222-2222-2222-2222-222222222222")).branches(bc)

	require.Len(t, bs, 2) // direct + deputy
	assert.Contains(t, bs[0], "oa.employee_id")
	assert.Contains(t, bs[0], "JOIN domain.clinics c ON c.organization_id = oa.organization_id")
	assert.Contains(t, bs[0], "WHERE c.id = @scope0")
	assert.Contains(t, bs[0], "@caller0")

	assert.Contains(t, bs[1], "oa.deputy_employee_id")
	assert.Contains(t, bs[1], "domain.employee_vacations")
	assert.Contains(t, bs[1], activeVacationPredicate)
}

func TestOrgAdminOf_Organization_NoJoin(t *testing.T) {
	bc := &branchCtx{callerID: "bob"}
	bs := OrgAdminOf.Organization(uuid.MustParse("33333333-3333-3333-3333-333333333333")).branches(bc)
	require.Len(t, bs, 2)
	assert.Contains(t, bs[0], "WHERE oa.organization_id = @scope0")
	assert.NotContains(t, bs[0], "JOIN domain.clinics")
}

func TestOrgAdminOf_Describe(t *testing.T) {
	id := uuid.New()
	assert.Equal(t, "organization administrator", OrgAdminOf.Organization(id).describe())
	assert.Equal(t, "organization administrator", OrgAdminOf.Clinic(id).describe())
	assert.Equal(t, "organization administrator", OrgAdminOf.IncidentType(id).describe())
}

func TestAdminOf_Clinic_IsAnyOf_SystemAdmin_Plus_OrgAdmin(t *testing.T) {
	id := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	bc1 := &branchCtx{callerID: "x"}
	bc2 := &branchCtx{callerID: "x"}

	got := AdminOf.Clinic(id).branches(bc1)
	want := AnyOf(SystemAdmin, OrgAdminOf.Clinic(id)).branches(bc2)

	assert.Equal(t, want, got)
}

func TestAdminOf_Describe(t *testing.T) {
	id := uuid.New()
	assert.Equal(t,
		"system administrator or organization administrator",
		AdminOf.Clinic(id).describe())
}

func TestRequire_NilPolicy_ReturnsAuthzCheckFailed(t *testing.T) {
	a := &Authz{}
	err := a.Require(t.Context(), "x", nil)
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe))
	assert.Equal(t, ErrCodeAuthzCheckFailed, oe.Code())
}

func TestRequire_EmptyAnyOf_ReturnsAuthzCheckFailed(t *testing.T) {
	a := &Authz{}
	err := a.Require(t.Context(), "x", AnyOf())
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe))
	assert.Equal(t, ErrCodeAuthzCheckFailed, oe.Code())
}

func TestRequire_ZeroScope_ReturnsAuthzCheckFailed(t *testing.T) {
	a := &Authz{}
	err := a.Require(t.Context(), "x", OrgAdminOf.Clinic(uuid.Nil))
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe))
	assert.Equal(t, ErrCodeAuthzCheckFailed, oe.Code())
}

func TestBranchCtx_HasZeroScope(t *testing.T) {
	bc := &branchCtx{callerID: "x"}
	bc.addCaller()
	bc.addScope(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	assert.False(t, bc.hasZeroScope())

	bc.addScope(uuid.Nil)
	assert.True(t, bc.hasZeroScope())
}
