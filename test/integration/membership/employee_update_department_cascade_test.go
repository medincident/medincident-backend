//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/services/membership"
)

func TestUpdateEmployeeDepartment_CascadeRevokesDR_SameClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           aliceID,
		DepartmentID: f.DeptA1b,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		f.DeptA1a, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)

	rows := latestOutbox(t)
	require.GreaterOrEqual(t, len(rows), 2)
	assert.Equal(t, "medincident.event.employee.v1.department_changed", rows[0].Subject)
	assert.Equal(t, membership.SubjectDepartmentResponsibleRevoked, rows[1].Subject)
}

func TestUpdateEmployeeDepartment_CascadeRevokesDRWithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID,
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID, DeputyEmployeeID: bobID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           aliceID,
		DepartmentID: f.DeptA1b,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 3, "expected 3 events: DepartmentChanged, DeputyRemoved, Revoked")
	assert.Equal(t, "medincident.event.employee.v1.department_changed", rows[0].Subject, "Rule 1: cause first")
	assert.Equal(t, membership.SubjectDepartmentResponsibleDeputyRemoved, rows[1].Subject, "Rule 2: cleanup before terminate")
	assert.Equal(t, membership.SubjectDepartmentResponsibleRevoked, rows[2].Subject)
}
