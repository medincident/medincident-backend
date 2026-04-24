//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

// hireAlice is a shared helper for tests that need a baseline employee.
func hireAlice(t *testing.T, f fixture) (employeeID string) {
	t.Helper()
	res, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
		},
	})
	require.NoError(t, err)
	return res.ID.String()
}

func TestUpdateEmployeePosition_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	newPos := "Head nurse"
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeePositionPayload{
			ID:       id.String(),
			Position: &newPos,
		},
	}))

	var projPos string
	require.NoError(t, testDB.Raw(`SELECT position FROM projections.employees WHERE id = ?`, id).Row().Scan(&projPos))
	assert.Equal(t, "Head nurse", projPos)
}

func TestUpdateEmployeePosition_ClearPosition(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	pos := "Nurse"
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		Caller:  sysadminCaller,
		Payload: membership.UpdateEmployeePositionPayload{ID: id.String(), Position: &pos},
	}))

	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		Caller:  sysadminCaller,
		Payload: membership.UpdateEmployeePositionPayload{ID: id.String(), Position: nil},
	}))
}

func TestUpdateEmployeePosition_NoOp(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	// Alice was hired with no position; updating with nil → no-op.
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		Caller:  sysadminCaller,
		Payload: membership.UpdateEmployeePositionPayload{ID: id.String(), Position: nil},
	}))
	var projPos *string
	require.NoError(t, testDB.Raw(`SELECT position FROM projections.employees WHERE id = ?`, id).Row().Scan(&projPos))
	assert.Nil(t, projPos)
}

func TestUpdateEmployeePosition_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeePositionPayload{
			ID: uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}
