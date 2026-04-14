//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// hireAlice is a shared helper for tests that need a baseline employee.
func hireAlice(t *testing.T, f fixture) (employeeID string) {
	t.Helper()
	res, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
	})
	require.NoError(t, err)
	return res.ID.String()
}

func TestUpdateEmployeePosition_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)
	newPos := "Head nurse"
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		ID:       id,
		Position: &newPos,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectEmployeePositionChanged, rows[0].Subject)
	ev := &employeev1.EmployeePositionChanged{}
	decodePayload(t, rows[0], ev)
	require.NotNil(t, ev.Position)
	assert.Equal(t, "Head nurse", *ev.Position)
}

func TestUpdateEmployeePosition_ClearPosition(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	pos := "Nurse"
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{ID: id, Position: &pos}))

	truncateOutbox(t)
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{ID: id, Position: nil}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	ev := &employeev1.EmployeePositionChanged{}
	decodePayload(t, rows[0], ev)
	assert.Nil(t, ev.Position)
}

func TestUpdateEmployeePosition_NoOp(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)
	// Alice was hired with no position; updating with nil → no-op.
	require.NoError(t, empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{ID: id, Position: nil}))
	require.Empty(t, latestOutbox(t))
}

func TestUpdateEmployeePosition_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.UpdatePosition(ctxT(t), membership.UpdateEmployeePositionCommand{
		ID: uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}
