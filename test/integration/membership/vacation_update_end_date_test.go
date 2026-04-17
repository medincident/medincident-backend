//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestUpdateVacationEndDate_ExtendRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	end1 := time.Now().Add(24 * time.Hour)
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		EmployeeID: mustParseUUID(t, empID),
		EndsAt:     &end1,
	})
	require.NoError(t, err)
	newEnd := time.Now().Add(72 * time.Hour)
	require.NoError(t, empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: res.ID,
		EndsAt:     newEnd,
	}))
}

func TestUpdateVacationEndDate_ShortenRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	end1 := time.Now().Add(72 * time.Hour)
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		EmployeeID: mustParseUUID(t, empID),
		EndsAt:     &end1,
	})
	require.NoError(t, err)
	newEnd := time.Now().Add(1 * time.Hour)
	require.NoError(t, empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: res.ID,
		EndsAt:     newEnd,
	}))
}

func TestUpdateVacationEndDate_SetOnUnlimitedRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	end := time.Now().Add(72 * time.Hour)
	require.NoError(t, empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: mustParseUUID(t, vacID),
		EndsAt:     end,
	}))
}

func TestUpdateVacationEndDate_SetOnUnlimitedScheduled(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		EmployeeID: mustParseUUID(t, empID),
		StartsAt:   start,
	})
	require.NoError(t, err)
	end := start.Add(48 * time.Hour)
	require.NoError(t, empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: res.ID,
		EndsAt:     end,
	}))
}

func TestUpdateVacationEndDate_EndInPast(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	err := empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: mustParseUUID(t, vacID),
		EndsAt:     time.Now().Add(-1 * time.Hour),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationEndInPast, oopsCode(t, err))
}

func TestUpdateVacationEndDate_EndBeforeStart(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		EmployeeID: mustParseUUID(t, empID),
		StartsAt:   start,
	})
	require.NoError(t, err)
	// end < start (but still in the future relative to now)
	err = empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: res.ID,
		EndsAt:     start.Add(-1 * time.Hour),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationEndBeforeStart, oopsCode(t, err))
}

func TestUpdateVacationEndDate_Overlap(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	// v1 = running [now, now+24h)
	end1 := time.Now().Add(24 * time.Hour)
	v1, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id, EndsAt: &end1})
	require.NoError(t, err)
	// v2 = scheduled [now+30h, now+40h)
	start2 := time.Now().Add(30 * time.Hour)
	end2 := time.Now().Add(40 * time.Hour)
	_, err = empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		EmployeeID: id, StartsAt: start2, EndsAt: &end2,
	})
	require.NoError(t, err)
	// Extend v1 into v2's range.
	err = empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: v1.ID,
		EndsAt:     time.Now().Add(35 * time.Hour),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}

func TestUpdateVacationEndDate_NoOp(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	end1 := time.Now().Add(24 * time.Hour)
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		EmployeeID: mustParseUUID(t, empID),
		EndsAt:     &end1,
	})
	require.NoError(t, err)
	// Same end as current → no-op.
	require.NoError(t, empSvc.UpdateVacationEndDate(ctxT(t), membership.UpdateVacationEndDateCommand{
		VacationID: res.ID,
		EndsAt:     end1,
	}))
}
