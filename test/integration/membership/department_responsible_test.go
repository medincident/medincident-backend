//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func TestAssignDepartmentResponsible_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   id.String(),
		},
	}))

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
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: uuidMustV7().String(),
			EmployeeID:   id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsible_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsible_EmployeeNotInDepartment(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1b.String(), // different department
			EmployeeID:   mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInDepartment, oopsCode(t, err))
}

func TestAssignDepartmentResponsible_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   id.String(),
		},
	}))

	err := empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   id.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeDepartmentResponsible_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   id.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeDepartmentResponsible(ctxT(t), membership.RevokeDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   id.String(),
		},
	}))

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
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RevokeDepartmentResponsible(ctxT(t), membership.RevokeDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))
}

func TestRevokeDepartmentResponsible_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeDepartmentResponsible(ctxT(t), membership.RevokeDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.RevokeDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))

	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))
}

func TestAssignDepartmentResponsibleDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyNotInDepartment(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	// Carol is hired into DeptA1b (different department).
	carolRes, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserCarolID,
			DepartmentID:  f.DeptA1b.String(),
		},
	})
	require.NoError(t, err)
	carolID := carolRes.ID

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))

	err = empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: carolID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInDepartment, oopsCode(t, err))
}

func TestAssignDepartmentResponsibleDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	err := empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveDepartmentResponsibleDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsibleDeputyPayload{
			DepartmentID:     f.DeptA1a.String(),
			EmployeeID:       aliceID.String(),
			DeputyEmployeeID: bobID.String(),
		},
	}))

	require.NoError(t, empSvc.RemoveDepartmentResponsibleDeputy(ctxT(t), membership.RemoveDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveDepartmentResponsibleDeputyPayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))
}

func TestRemoveDepartmentResponsibleDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		Caller: sysadminCaller,
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	}))

	err := empSvc.RemoveDepartmentResponsibleDeputy(ctxT(t), membership.RemoveDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveDepartmentResponsibleDeputyPayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   aliceID.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveDepartmentResponsibleDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveDepartmentResponsibleDeputy(ctxT(t), membership.RemoveDepartmentResponsibleDeputyCommand{
		Caller: sysadminCaller,
		Payload: membership.RemoveDepartmentResponsibleDeputyPayload{
			DepartmentID: f.DeptA1a.String(),
			EmployeeID:   mustParseUUID(t, empID).String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentResponsibleNotFound, oopsCode(t, err))
}
