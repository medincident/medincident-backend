//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func TestAssignClinicHead_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: id,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		f.ClinicA1, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestAssignClinicHead_ClinicNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	err := empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   uuidMustV7(),
		EmployeeID: id,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeClinicNotFound, oopsCode(t, err))
}

func TestAssignClinicHead_EmployeeNotFound(t *testing.T) {
	f := takeFixture(t)

	err := empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestAssignClinicHead_EmployeeNotInClinic(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // hired into DeptA1a → ClinicA1

	// Target ClinicA2 — Alice is not in that clinic.
	err := empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA2,
		EmployeeID: mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotInClinic, oopsCode(t, err))
}

func TestAssignClinicHead_AlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: id,
	}))

	err := empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: id,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeClinicHeadAlreadyAssigned, oopsCode(t, err))
}

func TestRevokeClinicHead_Success_NoDeputy(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: id,
	}))

	require.NoError(t, empSvc.RevokeClinicHead(ctxT(t), membership.RevokeClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: id,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		f.ClinicA1, id,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestRevokeClinicHead_Success_WithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))
	require.NoError(t, empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	require.NoError(t, empSvc.RevokeClinicHead(ctxT(t), membership.RevokeClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))
}

func TestRevokeClinicHead_NotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RevokeClinicHead(ctxT(t), membership.RevokeClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeClinicHeadNotFound, oopsCode(t, err))
}

func TestAssignClinicHeadDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))
}

func TestAssignClinicHeadDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	err := empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeClinicHeadNotFound, oopsCode(t, err))
}

func TestAssignClinicHeadDeputy_DeputyNotFound(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))

	err := empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotFound, oopsCode(t, err))
}

func TestAssignClinicHeadDeputy_DeputyIsHolder(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))

	err := empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: aliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyIsHolder, oopsCode(t, err))
}

func TestAssignClinicHeadDeputy_DeputyNotInClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	// Carol is hired into DeptA2a (ClinicA2 — different clinic).
	carolRes, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserCarolID,
		DepartmentID:  f.DeptA2a,
	})
	require.NoError(t, err)
	carolID := carolRes.ID

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))

	// Carol is in ClinicA2 but Alice's CH role is in ClinicA1 — expect error.
	err = empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: carolID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotInClinic, oopsCode(t, err))
}

func TestAssignClinicHeadDeputy_DeputyAlreadyAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))
	require.NoError(t, empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	err := empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyAlreadyAssigned, oopsCode(t, err))
}

func TestRemoveClinicHeadDeputy_Success(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))
	require.NoError(t, empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID:         f.ClinicA1,
		EmployeeID:       aliceID,
		DeputyEmployeeID: bobID,
	}))

	require.NoError(t, empSvc.RemoveClinicHeadDeputy(ctxT(t), membership.RemoveClinicHeadDeputyCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))
}

func TestRemoveClinicHeadDeputy_DeputyNotAssigned(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))

	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	}))

	err := empSvc.RemoveClinicHeadDeputy(ctxT(t), membership.RemoveClinicHeadDeputyCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: aliceID,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDeputyNotAssigned, oopsCode(t, err))
}

func TestRemoveClinicHeadDeputy_RoleNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	err := empSvc.RemoveClinicHeadDeputy(ctxT(t), membership.RemoveClinicHeadDeputyCommand{
		ClinicID:   f.ClinicA1,
		EmployeeID: mustParseUUID(t, empID),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeClinicHeadNotFound, oopsCode(t, err))
}
