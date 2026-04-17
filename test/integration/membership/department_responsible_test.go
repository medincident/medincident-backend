//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/membership"
	departmentv1 "github.com/medincident/medincident-command-service/pkg/event/department/v1"
)

func TestAssignDepartmentResponsible_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   id,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectDepartmentResponsibleAssigned, rows[0].Subject)
	ev := &departmentv1.DepartmentResponsibleAssigned{}
	env := decodePayload(t, rows[0], ev)
	assert.Equal(t, membership.AggregateTypeDepartment, env.AggregateType)
	assert.Equal(t, f.DeptA1a.String(), env.AggregateId)
	assert.Equal(t, id.String(), ev.EmployeeId)

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		f.DeptA1a, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestAssignDepartmentResponsible_DepartmentNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: uuidMustV7(),
		EmployeeID:   id,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsible_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsible_EmployeeNotInDepartment(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1b, // different department
		EmployeeID:   mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInDepartment, oopsCode(t, err))
}

func TestAssignDepartmentResponsible_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   id,
	}))

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   id,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeDepartmentResponsible_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   id,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RevokeDepartmentResponsible(ctxT(t), membership.RevokeDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   id,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	assert.Equal(t, membership.SubjectDepartmentResponsibleRevoked, rows[0].Subject)

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		f.DeptA1a, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestRevokeDepartmentResponsible_Success_WithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RevokeDepartmentResponsible(ctxT(t), membership.RevokeDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 2)
	assert.Equal(t, membership.SubjectDepartmentResponsibleDeputyRemoved, rows[0].Subject, "Rule 2: cleanup before terminate")
	assert.Equal(t, membership.SubjectDepartmentResponsibleRevoked, rows[1].Subject)
}

func TestRevokeDepartmentResponsible_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeDepartmentResponsible(ctxT(t), membership.RevokeDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	assert.Equal(t, membership.SubjectDepartmentResponsibleDeputyAssigned, rows[0].Subject)
	ev := &departmentv1.DepartmentResponsibleDeputyAssigned{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, aliceID.String(), ev.EmployeeId)
	assert.Equal(t, bobID.String(), ev.DeputyEmployeeId)
}

func TestAssignDepartmentResponsibleDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: aliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyNotInDepartment(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	// Carol is hired into DeptA1b (different department).
	carolRes, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserCarolID,
		DepartmentID:  f.DeptA1b,
	})
	require.NoError(t, err)
	carolID := carolRes.ID

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))

	err = empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: carolID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInDepartment, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveDepartmentResponsibleDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     f.DeptA1a,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.RemoveDepartmentResponsibleDeputy(ctxT(t), membership.RemoveDepartmentResponsibleDeputyCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	assert.Equal(t, membership.SubjectDepartmentResponsibleDeputyRemoved, rows[0].Subject)
	ev := &departmentv1.DepartmentResponsibleDeputyRemoved{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, aliceID.String(), ev.EmployeeId)
}

func TestRemoveDepartmentResponsibleDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	}))

	err := empSvc.RemoveDepartmentResponsibleDeputy(ctxT(t), membership.RemoveDepartmentResponsibleDeputyCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   aliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveDepartmentResponsibleDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveDepartmentResponsibleDeputy(ctxT(t), membership.RemoveDepartmentResponsibleDeputyCommand{
		DepartmentID: f.DeptA1a,
		EmployeeID:   mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleNotFound, oopsCode(t, err))
}
