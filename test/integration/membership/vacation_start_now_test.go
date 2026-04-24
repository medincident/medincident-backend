//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

// startUnlimitedVacation is a helper reused by F6-F9 tests.
func startUnlimitedVacation(t *testing.T, empID string) string {
	t.Helper()
	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: mustParseUUID(t, empID).String()},
	})
	require.NoError(t, err)
	return res.ID.String()
}

func TestStartVacationNow_Success_Unlimited(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	res, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: id.String()},
	})
	require.NoError(t, err)
	require.NotEqual(t, "", res.ID.String())
}

func TestStartVacationNow_Success_WithEnd(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	end := time.Now().Add(24 * time.Hour)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: id.String(), EndsAt: &end},
	})
	require.NoError(t, err)
}

func TestStartVacationNow_EndBeforeStart(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	past := time.Now().Add(-1 * time.Hour)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: id.String(), EndsAt: &past},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationEndBeforeStart, oopsCode(t, err))
}

func TestStartVacationNow_Overlap_WithRunning(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: id.String()},
	})
	require.NoError(t, err)
	_, err = empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: id.String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeVacationOverlap, oopsCode(t, err))
}

func TestStartVacationNow_EmployeeNotFound(t *testing.T) {
	_ = takeFixture(t)
	_, err := empSvc.StartVacationNow(ctxT(t), membership.StartVacationNowCommand{
		Caller:  sysadminCaller,
		Payload: membership.StartVacationNowPayload{EmployeeID: uuidMustV7().String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}
