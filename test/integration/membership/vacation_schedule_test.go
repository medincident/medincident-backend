//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func TestScheduleVacation_Success_Unlimited(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ScheduleVacationPayload{EmployeeID: id.String(), StartsAt: start},
	})
	require.NoError(t, err)

	var domainCount int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employee_vacations WHERE id = ?`, res.ID).Scan(&domainCount).Error)
	require.Equal(t, int64(1), domainCount)
}

func TestScheduleVacation_Success_WithEnd(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(72 * time.Hour)
	end := start.Add(48 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ScheduleVacationPayload{EmployeeID: id.String(), StartsAt: start, EndsAt: &end},
	})
	require.NoError(t, err)
}

func TestScheduleVacation_StartInPast(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(-1 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ScheduleVacationPayload{EmployeeID: id.String(), StartsAt: start},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationStartInPast, oopsCode(t, err))
}

func TestScheduleVacation_EndBeforeStart(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(72 * time.Hour)
	end := start.Add(-1 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ScheduleVacationPayload{EmployeeID: id.String(), StartsAt: start, EndsAt: &end},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationEndBeforeStart, oopsCode(t, err))
}

func TestScheduleVacation_Overlap(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start1 := time.Now().Add(72 * time.Hour)
	end1 := start1.Add(48 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ScheduleVacationPayload{EmployeeID: id.String(), StartsAt: start1, EndsAt: &end1},
	})
	require.NoError(t, err)
	// Overlapping second vacation.
	start2 := start1.Add(24 * time.Hour)
	end2 := start2.Add(48 * time.Hour)
	_, err = empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller:  sysadminCaller,
		Payload: membership.ScheduleVacationPayload{EmployeeID: id.String(), StartsAt: start2, EndsAt: &end2},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}

// Regression test for the StartVacationNow overlap path that depends on ScheduleVacation existing.
func TestStartVacationNow_Overlap_WithScheduled(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	futureStart := time.Now().Add(1 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		Caller: sysadminCaller,
		Payload: membership.ScheduleVacationPayload{
			EmployeeID: id.String(),
			StartsAt:   futureStart,
		},
	})
	require.NoError(t, err)
	// Start now with unlimited end → overlaps with the scheduled one.
	_, err = empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: id.String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}
