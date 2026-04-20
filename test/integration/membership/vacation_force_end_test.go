//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestForceEndVacation_Success_UnlimitedRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ForceEndVacationPayload{VacationID: mustParseUUID(t, vacID).String()},
	}))
}

func TestForceEndVacation_Success_BoundedRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	end := time.Now().Add(72 * time.Hour)
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller: sysadminCaller,
		Payload: membership.StartVacationNowPayload{
			EmployeeID: mustParseUUID(t, empID).String(),
			EndsAt:     &end,
		},
	})
	require.NoError(t, err)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ForceEndVacationPayload{VacationID: res.ID.String()},
	}))
}

func TestForceEndVacation_NotStarted(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: mustParseUUID(t, empID).String(),
			StartsAt:   start,
		},
	})
	require.NoError(t, err)
	err = empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ForceEndVacationPayload{VacationID: res.ID.String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationNotStarted, oopsCode(t, err))
}

func TestForceEndVacation_AlreadyEnded(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ForceEndVacationPayload{VacationID: mustParseUUID(t, vacID).String()},
	}))
	// Second attempt.
	err := empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ForceEndVacationPayload{VacationID: mustParseUUID(t, vacID).String()},
	})
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
		Caller: sysadminCaller,
		Payload: membership.StartVacationNowPayload{
			EmployeeID: mustParseUUID(t, empID).String(),
			EndsAt:     &end1,
		},
	})
	require.NoError(t, err)
	require.NoError(t, empSvc.ForceEndVacation(ctxT(t), membership.ForceEndVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ForceEndVacationPayload{VacationID: res.ID.String()},
	}))
	// Now try to change end date.
	err = empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateVacationEndDatePayload{
			VacationID: res.ID.String(),
			EndsAt:     time.Now().Add(72 * time.Hour),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationAlreadyEnded, oopsCode(t, err))
}
