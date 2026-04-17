//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	organizationv1 "github.com/medincident/medincident-command-service/pkg/event/organization/v1"
)

func TestAssignOrganizationHead_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     id,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectOrganizationHeadAssigned, rows[0].Subject)
	ev := &organizationv1.OrganizationHeadAssigned{}
	env := decodePayload(t, rows[0], ev)
	assert.Equal(t, membership.AggregateTypeOrganization, env.AggregateType)
	assert.Equal(t, f.OrgA.String(), env.AggregateId)
	assert.Equal(t, id.String(), ev.EmployeeId)

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
		OrganizationID: uuidMustV7(),
		EmployeeID:     id,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHead_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHead_EmployeeNotInOrganization(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a → OrgA

	// Target OrgB — Alice is not in that organization.
	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgB,
		EmployeeID:     mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationHead_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     id,
	}))

	err := empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     id,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeOrganizationHead_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     id,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RevokeOrganizationHead(ctxT(t), membership.RevokeOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     id,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	assert.Equal(t, membership.SubjectOrganizationHeadRevoked, rows[0].Subject)

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
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))
	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RevokeOrganizationHead(ctxT(t), membership.RevokeOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 2)
	assert.Equal(t, membership.SubjectOrganizationHeadDeputyRemoved, rows[0].Subject, "Rule 2: cleanup before terminate")
	assert.Equal(t, membership.SubjectOrganizationHeadRevoked, rows[1].Subject)
}

func TestRevokeOrganizationHead_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeOrganizationHead(ctxT(t), membership.RevokeOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	assert.Equal(t, membership.SubjectOrganizationHeadDeputyAssigned, rows[0].Subject)
	ev := &organizationv1.OrganizationHeadDeputyAssigned{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, aliceID.String(), ev.EmployeeId)
	assert.Equal(t, bobID.String(), ev.DeputyEmployeeId)
}

func TestAssignOrganizationHeadDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: aliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyNotInOrganization(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f)) // hired into DeptA1a → OrgA

	// Carol is hired into DeptB1a (OrgB — different organization).
	carolRes, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserCarolID,
		DepartmentID:  f.DeptB1a,
	})
	require.NoError(t, err)
	carolID := carolRes.ID

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))

	// Carol is in OrgB but Alice's OrgHead role is in OrgA — expect error.
	err = empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: carolID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInOrganization, oopsCode(t, err))
}

func TestAssignOrganizationHeadDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))
	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	err := empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationHeadDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))
	require.NoError(t, empSvc.AssignOrganizationHeadDeputy(ctxT(t), membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   f.OrgA,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RemoveOrganizationHeadDeputy(ctxT(t), membership.RemoveOrganizationHeadDeputyCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	assert.Equal(t, membership.SubjectOrganizationHeadDeputyRemoved, rows[0].Subject)
	ev := &organizationv1.OrganizationHeadDeputyRemoved{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, aliceID.String(), ev.EmployeeId)
}

func TestRemoveOrganizationHeadDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	}))

	err := empSvc.RemoveOrganizationHeadDeputy(ctxT(t), membership.RemoveOrganizationHeadDeputyCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     aliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveOrganizationHeadDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveOrganizationHeadDeputy(ctxT(t), membership.RemoveOrganizationHeadDeputyCommand{
		OrganizationID: f.OrgA,
		EmployeeID:     mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeOrganizationHeadNotFound, oopsCode(t, err))
}
