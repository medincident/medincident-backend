//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestAssignOrganizationHead_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_heads WHERE organization_id = ? AND employee_id = ?`,
		f.OrgA, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestAssignOrganizationHead_OrganizationNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: uuidMustV7().String(),
			EmployeeID:     id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHead_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHead_EmployeeNotInOrganization(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a → OrgA

	// Target OrgB — Alice is not in that organization.
	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgB.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationHead_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeOrganizationHead_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeOrganizationHead(ctxT(t), membership.RevokeOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_heads WHERE organization_id = ? AND employee_id = ?`,
		f.OrgA, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestRevokeOrganizationHead_Success_WithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeOrganizationHead(ctxT(t), membership.RevokeOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
}

func TestRevokeOrganizationHead_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeOrganizationHead(ctxT(t), membership.RevokeOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))
}

func TestAssignOrganizationHeadDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyNotInOrganization(t *testing.T) {
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

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	// Carol is in OrgB but Alice's OrgHead role is in OrgA — expect error.
	err = empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: carolID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationHeadDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RemoveOrganizationHeadDeputy(ctxT(t), membership.RemoveOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationHeadDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
}

func TestRemoveOrganizationHeadDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.RemoveOrganizationHeadDeputy(ctxT(t), membership.RemoveOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationHeadDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationHeadDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveOrganizationHeadDeputy(ctxT(t), membership.RemoveOrganizationHeadDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationHeadDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadNotFound, oopsCode(t, err))
}
