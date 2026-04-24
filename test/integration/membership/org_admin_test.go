//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func TestAssignOrganizationAdmin_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_admins WHERE organization_id = ? AND employee_id = ?`,
		f.OrgA, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestAssignOrganizationAdmin_OrganizationNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	err := empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: uuidMustV7().String(),
			EmployeeID:     id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationNotFound, oopsCode(t, err))
}

func TestAssignOrganizationAdmin_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignOrganizationAdmin_EmployeeNotInOrganization(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a → OrgA

	// Target OrgB — Alice is not in that organization.
	err := empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgB.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationAdmin_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	err := empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationAdminAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeOrganizationAdmin_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeOrganizationAdmin(ctxT(t), membership.RevokeOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_admins WHERE organization_id = ? AND employee_id = ?`,
		f.OrgA, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestRevokeOrganizationAdmin_Success_WithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeOrganizationAdmin(ctxT(t), membership.RevokeOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
}

func TestRevokeOrganizationAdmin_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeOrganizationAdmin(ctxT(t), membership.RevokeOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationAdminNotFound, oopsCode(t, err))
}

func TestAssignOrganizationAdminDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))
}

func TestAssignOrganizationAdminDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationAdminNotFound, oopsCode(t, err))
}

func TestAssignOrganizationAdminDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignOrganizationAdminDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignOrganizationAdminDeputy_DeputyNotInOrganization(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f)) // hired into DeptA1a → OrgA

	// Carol is hired into DeptB1a (OrgB — different organization).
	carolRes, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserCarolID,
			DepartmentID:  f.DeptB1a.String(),
		},
	})
	require.NoError(t, err)
	carolID := carolRes.ID

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	// Carol is in OrgB but Alice's OrgAdmin role is in OrgA — expect error.
	err = empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: carolID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationAdminDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	err := empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationAdminDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationAdminDeputy(ctxT(t), membership.AssignOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RemoveOrganizationAdminDeputy(ctxT(t), membership.RemoveOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationAdminDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
}

func TestRemoveOrganizationAdminDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.RemoveOrganizationAdminDeputy(ctxT(t), membership.RemoveOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationAdminDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationAdminDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveOrganizationAdminDeputy(ctxT(t), membership.RemoveOrganizationAdminDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationAdminDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationAdminNotFound, oopsCode(t, err))
}
