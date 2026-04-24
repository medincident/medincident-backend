//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func TestAssignOrganizationDispatcher_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_dispatchers WHERE organization_id = ? AND employee_id = ?`,
		f.OrgA, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestAssignOrganizationDispatcher_OrganizationNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	err := empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: uuidMustV7().String(),
			EmployeeID:     id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationNotFound, oopsCode(t, err))
}

func TestAssignOrganizationDispatcher_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignOrganizationDispatcher_EmployeeNotInOrganization(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a → OrgA

	// Target OrgB — Alice is not in that organization.
	err := empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgB.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationDispatcher_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	err := empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationDispatcherAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeOrganizationDispatcher_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeOrganizationDispatcher(ctxT(t), membership.RevokeOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     id.String(),
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_dispatchers WHERE organization_id = ? AND employee_id = ?`,
		f.OrgA, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestRevokeOrganizationDispatcher_Success_WithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeOrganizationDispatcher(ctxT(t), membership.RevokeOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
}

func TestRevokeOrganizationDispatcher_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeOrganizationDispatcher(ctxT(t), membership.RevokeOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationDispatcherNotFound, oopsCode(t, err))
}

func TestAssignOrganizationDispatcherDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))
}

func TestAssignOrganizationDispatcherDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationDispatcherNotFound, oopsCode(t, err))
}

func TestAssignOrganizationDispatcherDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignOrganizationDispatcherDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignOrganizationDispatcherDeputy_DeputyNotInOrganization(t *testing.T) {
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

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	// Carol is in OrgB but Alice's OrgDispatcher role is in OrgA — expect error.
	err = empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: carolID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationDispatcherDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	err := empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationDispatcherDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignOrganizationDispatcherDeputy(ctxT(t), membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   f.OrgA.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RemoveOrganizationDispatcherDeputy(ctxT(t), membership.RemoveOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationDispatcherDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))
}

func TestRemoveOrganizationDispatcherDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	}))

	err := empSvc.RemoveOrganizationDispatcherDeputy(ctxT(t), membership.RemoveOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationDispatcherDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationDispatcherDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveOrganizationDispatcherDeputy(ctxT(t), membership.RemoveOrganizationDispatcherDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveOrganizationDispatcherDeputyPayload{
			OrganizationID: f.OrgA.String(),
			EmployeeID:     mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationDispatcherNotFound, oopsCode(t, err))
}
