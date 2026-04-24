//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestTerminateEmployee_CascadeRevokesDRAsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(), EmployeeID: aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: aliceID.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE employee_id = ?`, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestTerminateEmployee_CascadeRevokesCHAsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignClinicHeadPayload{
			ClinicID: f.ClinicA1.String(), EmployeeID: aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: aliceID.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE employee_id = ?`, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestTerminateEmployee_CascadeRevokesOrgAdminAsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignOrganizationAdmin(ctxT(t), membership.AssignOrganizationAdminCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationAdminPayload{
			OrganizationID: f.OrgA.String(), EmployeeID: aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: aliceID.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_admins WHERE employee_id = ?`, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestTerminateEmployee_CascadeRevokesOrgHeadAsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignOrganizationHead(ctxT(t), membership.AssignOrganizationHeadCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: f.OrgA.String(), EmployeeID: aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: aliceID.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_heads WHERE employee_id = ?`, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestTerminateEmployee_CascadeRevokesOrgDispatcherAsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignOrganizationDispatcher(ctxT(t), membership.AssignOrganizationDispatcherCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: f.OrgA.String(), EmployeeID: aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: aliceID.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_dispatchers WHERE employee_id = ?`, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestTerminateEmployee_CascadeClearsDRDeputy(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(), EmployeeID: aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID: f.DeptA1a.String(), EmployeeID: aliceID.String(), DeputyEmployeeID: bobID.String(),
		},
	}))

	// Terminate Bob — he was the deputy, not the holder. Alice's role remains.
	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: bobID.String()},
	}))

	var deputyCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE deputy_employee_id = ?`, bobID,
	).Scan(&deputyCount).Error)
	assert.Equal(t, int64(0), deputyCount, "deputy slot must be cleared")

	var roleCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE employee_id = ? AND department_id = ?`, aliceID, f.DeptA1a,
	).Scan(&roleCount).Error)
	assert.Equal(t, int64(1), roleCount, "Alice's role row must remain")
}
