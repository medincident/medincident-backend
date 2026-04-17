//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// startUnlimitedVacation is a helper reused by F6-F9 tests.
func startUnlimitedVacation(t *testing.T, empID string) string {
	t.Helper()
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: mustParseUUID(t, empID)})
	require.NoError(t, err)
	return res.ID.String()
}

func TestStartVacationNow_Success_Unlimited(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)

	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id})
	require.NoError(t, err)
	require.NotEqual(t, "", res.ID.String())

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectVacationStarted, rows[0].Subject)
	ev := &employeev1.VacationStarted{}
	env := decodePayload(t, rows[0], ev)
	assert.Equal(t, id.String(), env.AggregateId) // aggregate is the employee
	assert.Equal(t, res.ID.String(), ev.VacationId)
	assert.Nil(t, ev.EndsAt)
}

func TestStartVacationNow_Success_WithEnd(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	end := time.Now().Add(24 * time.Hour)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id, EndsAt: &end})
	require.NoError(t, err)
}

func TestStartVacationNow_EndBeforeStart(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	past := time.Now().Add(-1 * time.Hour)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id, EndsAt: &past})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationEndBeforeStart, oopsCode(t, err))
}

func TestStartVacationNow_Overlap_WithRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id})
	require.NoError(t, err)
	_, err = empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}

func TestStartVacationNow_EmployeeNotFound(t *testing.T) {
	_ = takeFixture(t)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: uuidMustV7()})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}
