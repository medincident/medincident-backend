//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	employeev1 "github.com/medincident/medincident-command-service/pkg/event/employee/v1"
)

func TestCancelScheduledVacation_Success(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	start := time.Now().Add(72 * time.Hour)
	res, err := empSvc.ScheduleVacation(ctxT(t), membership.ScheduleVacationCommand{
		EmployeeID: mustParseUUID(t, empID),
		StartsAt:   start,
	})
	require.NoError(t, err)
	truncateOutbox(t)

	require.NoError(t, empSvc.CancelScheduledVacation(ctxT(t), membership.CancelScheduledVacationCommand{VacationID: res.ID}))

	var count int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employee_vacations WHERE id = ?`, res.ID).Scan(&count).Error)
	assert.Equal(t, int64(0), count)

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectVacationCancelled, rows[0].Subject)
	ev := &employeev1.VacationCancelled{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, res.ID.String(), ev.VacationId)
}

func TestCancelScheduledVacation_AlreadyStarted(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	vacID := startUnlimitedVacation(t, empID)
	err := empSvc.CancelScheduledVacation(ctxT(t), membership.CancelScheduledVacationCommand{VacationID: mustParseUUID(t, vacID)})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationAlreadyStarted, oopsCode(t, err))
}

func TestCancelScheduledVacation_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.CancelScheduledVacation(ctxT(t), membership.CancelScheduledVacationCommand{VacationID: uuidMustV7()})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationNotFound, oopsCode(t, err))
}
