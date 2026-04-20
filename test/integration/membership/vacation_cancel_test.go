//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
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

	require.NoError(t, empSvc.CancelScheduledVacation(ctxT(t), membership.CancelScheduledVacationCommand{VacationID: res.ID}))

	var count int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employee_vacations WHERE id = ?`, res.ID).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
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
