//go:build integration

package membership_integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestTerminateEmployee_Success_NoVacations(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: id.String()},
	}))

	var count int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employees WHERE id = ?`, id).Scan(&count).Error)
	assert.Equal(t, int64(0), count)

	// Projection keeps the row but stamps terminated_at.
	var terminatedAt *time.Time
	require.NoError(t, testDB.Raw(`SELECT terminated_at FROM projections.employees WHERE id = ?`, id).Row().Scan(&terminatedAt))
	require.NotNil(t, terminatedAt)
}

func TestTerminateEmployee_Success_WithVacations_CascadeDeletesAll(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)

	// Insert a vacation row directly via DB (vacation service comes in F5+).
	vacID := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.Exec(`
		INSERT INTO domain.employee_vacations (id, employee_id, starts_at, ends_at)
		VALUES (?, ?, ?, NULL)`, vacID, id, time.Now()).Error)

	var vCount int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employee_vacations WHERE employee_id = ?`, id).Scan(&vCount).Error)
	require.Equal(t, int64(1), vCount)

	require.NoError(t, empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: id.String()},
	}))

	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employee_vacations WHERE employee_id = ?`, id).Scan(&vCount).Error)
	assert.Equal(t, int64(0), vCount, "cascade delete should have removed the vacation")
}

func TestTerminateEmployee_NotFound(t *testing.T) {
	_ = takeFixture(t)
	err := empSvc.Terminate(ctxT(t), membership.TerminateEmployeeCommand{
		Caller:  sysadminCaller,
		Payload: membership.TerminateEmployeePayload{ID: uuidMustV7().String()},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeNotFound, oopsCode(t, err))
}
