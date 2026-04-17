//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestScheduleVacation_Success_Unlimited(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{EmployeeID: id, StartsAt: start})
	require.NoError(t, err)

	var projState string
	require.NoError(t, testDB.Raw(`SELECT state FROM projections.employee_vacations WHERE id = ?`, res.ID).Row().Scan(&projState))
	require.Equal(t, "scheduled", projState)
}

func TestScheduleVacation_Success_WithEnd(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(72 * time.Hour)
	end := start.Add(48 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{EmployeeID: id, StartsAt: start, EndsAt: &end})
	require.NoError(t, err)
}

func TestScheduleVacation_StartInPast(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(-1 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{EmployeeID: id, StartsAt: start})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationStartInPast, oopsCode(t, err))
}

func TestScheduleVacation_EndBeforeStart(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start := time.Now().Add(72 * time.Hour)
	end := start.Add(-1 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{EmployeeID: id, StartsAt: start, EndsAt: &end})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationEndBeforeStart, oopsCode(t, err))
}

func TestScheduleVacation_Overlap(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	start1 := time.Now().Add(72 * time.Hour)
	end1 := start1.Add(48 * time.Hour)
	_, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{EmployeeID: id, StartsAt: start1, EndsAt: &end1})
	require.NoError(t, err)
	// Overlapping second vacation.
	start2 := start1.Add(24 * time.Hour)
	end2 := start2.Add(48 * time.Hour)
	_, err = empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{EmployeeID: id, StartsAt: start2, EndsAt: &end2})
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
		EmployeeID: id,
		StartsAt:   futureStart,
	})
	require.NoError(t, err)
	// Start now with unlimited end → overlaps with the scheduled one.
	_, err = empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{EmployeeID: id})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}
