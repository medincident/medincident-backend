//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/membership"
	employeev1 "github.com/medincident/medincident-command-service/pkg/event/employee/v1"
)

func TestForceEndVacation_Success_UnlimitedRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	truncateOutbox(t)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{VacationID: mustParseUUID(t, vacID)}))

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectVacationEnded, rows[0].Subject)
	ev := &employeev1.VacationEnded{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, vacID, ev.VacationId)
	require.NotNil(t, ev.EndsAt)
}

func TestForceEndVacation_Success_BoundedRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	end := time.Now().Add(72 * time.Hour)
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		EmployeeID: mustParseUUID(t, empID),
		EndsAt:     &end,
	})
	require.NoError(t, err)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{VacationID: res.ID}))
}

func TestForceEndVacation_NotStarted(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		EmployeeID: mustParseUUID(t, empID),
		StartsAt:   start,
	})
	require.NoError(t, err)
	err = empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{VacationID: res.ID})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationNotStarted, oopsCode(t, err))
}

func TestForceEndVacation_AlreadyEnded(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{VacationID: mustParseUUID(t, vacID)}))
	// Second attempt.
	err := empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{VacationID: mustParseUUID(t, vacID)})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationAlreadyEnded, oopsCode(t, err))
}

// TestUpdateVacationEndDate_AlreadyEnded cross-validates F7 precondition:
// UpdateVacationEndDate cannot resurrect a force-ended vacation.
func TestUpdateVacationEndDate_AlreadyEnded(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	end1 := time.Now().Add(24 * time.Hour)
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		EmployeeID: mustParseUUID(t, empID),
		EndsAt:     &end1,
	})
	require.NoError(t, err)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{VacationID: res.ID}))
	// Now try to change end date.
	err = empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: res.ID,
		EndsAt:     time.Now().Add(72 * time.Hour),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationAlreadyEnded, oopsCode(t, err))
}
