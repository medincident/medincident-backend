package authz

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchCtx_SequentialPlaceholders(t *testing.T) {
	bc := &branchCtx{callerID: "alice"}

	c1 := bc.addCaller()
	s1 := bc.addScope(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	c2 := bc.addCaller()

	assert.Equal(t, "caller0", c1)
	assert.Equal(t, "scope1", s1)
	assert.Equal(t, "caller2", c2)
	require.Len(t, bc.args, 3)

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
