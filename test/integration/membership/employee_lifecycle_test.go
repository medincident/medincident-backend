//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func isEmployeeActive(t *testing.T, employeeID string) bool {
	t.Helper()
	var active bool
	require.NoError(t, testDB.Raw(`SELECT is_active FROM domain.employees WHERE id = ?`, employeeID).Row().Scan(&active))
	return active
}

func TestDeactivateEmployee_HappyPath(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	require.NoError(t, empSvc.Deactivate(ctxT(t), membership.DeactivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.DeactivateEmployeePayload{ID: empID},
	}))

	assert.False(t, isEmployeeActive(t, empID))
}

func TestDeactivateEmployee_EmitsEvent(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	var beforeCount int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM outbox.events WHERE subject = 'medincident.event.employee.v1.deactivated'`,
	).Scan(&beforeCount).Error)

	require.NoError(t, empSvc.Deactivate(ctxT(t), membership.DeactivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.DeactivateEmployeePayload{ID: empID},
	}))

	var afterCount int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM outbox.events WHERE subject = 'medincident.event.employee.v1.deactivated'`,
	).Scan(&afterCount).Error)
	assert.Equal(t, beforeCount+1, afterCount)
}

func TestDeactivateEmployee_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.Deactivate(ctxT(t), membership.DeactivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.DeactivateEmployeePayload{ID: uuidMustV7().String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}

func TestActivateEmployee_HappyPath(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	require.NoError(t, empSvc.Deactivate(ctxT(t), membership.DeactivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.DeactivateEmployeePayload{ID: empID},
	}))
	require.False(t, isEmployeeActive(t, empID))

	require.NoError(t, empSvc.Activate(ctxT(t), membership.ActivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.ActivateEmployeePayload{ID: empID},
	}))

	assert.True(t, isEmployeeActive(t, empID))
}

func TestActivateEmployee_ParentDeptInactive(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)

	// Deactivate the department directly (bypassing the service to avoid
	// cascade event noise and keep the test focused on activate behaviour).
	require.NoError(t, testDB.Exec(
		`UPDATE domain.departments SET is_active = FALSE WHERE id = ?`, f.DeptA1a,
	).Error)
	// Also deactivate the employee so we can try to re-activate it.
	require.NoError(t, testDB.Exec(
		`UPDATE domain.employees SET is_active = FALSE WHERE id = ?`, empID,
	).Error)

	err := empSvc.Activate(ctxT(t), membership.ActivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.ActivateEmployeePayload{ID: empID},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeActivateParentInactive, oopsCode(t, err))
}

func TestActivateEmployee_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.Activate(ctxT(t), membership.ActivateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.ActivateEmployeePayload{ID: uuidMustV7().String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}
