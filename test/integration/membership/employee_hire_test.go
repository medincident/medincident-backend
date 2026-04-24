//go:build integration

package membership_integration_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func hireBob(t *testing.T, f fixture) (employeeID string) {
	t.Helper()
	res, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserBobID,
			DepartmentID:  f.DeptA1a.String(),
		},
	})
	require.NoError(t, err)
	return res.ID.String()
}

func TestHireEmployee_Success_WithPosition(t *testing.T) {
	f := takeFixture(t)
	pos := "Senior nurse"
	res, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
			Position:      &pos,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, "", res.ID.String())

	var count int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employees WHERE id = ?`, res.ID).Scan(&count).Error)
	require.Equal(t, int64(1), count)

	// Projection row written atomically.
	var projCount int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM projections.employees WHERE id = ?`, res.ID).Scan(&projCount).Error)
	require.Equal(t, int64(1), projCount)

	var projPos string
	require.NoError(t, testDB.Raw(`SELECT position FROM projections.employees WHERE id = ?`, res.ID).Row().Scan(&projPos))
	require.Equal(t, "Senior nurse", projPos)
}

func TestHireEmployee_Success_NoPosition(t *testing.T) {
	f := takeFixture(t)
	res, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserBobID,
			DepartmentID:  f.DeptA1a.String(),
		},
	})
	require.NoError(t, err)

	var projPos *string
	require.NoError(t, testDB.Raw(`SELECT position FROM projections.employees WHERE id = ?`, res.ID).Row().Scan(&projPos))
	assert.Nil(t, projPos)
}

func TestHireEmployee_WhitespaceOnlyPositionRejected(t *testing.T) {
	f := takeFixture(t)
	empty := "   "
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
			Position:      &empty,
		},
	})
	require.Error(t, err)
	assertHasViolation(t, validationViolations(t, err), "position", "min")
}

func TestHireEmployee_ZitadelUserNotFound(t *testing.T) {
	f := takeFixture(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: "nonexistent-user-id",
			DepartmentID:  f.DeptA1a.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeZitadelUserNotFound, oopsCode(t, err))
}

func TestHireEmployee_DepartmentNotFound(t *testing.T) {
	_ = takeFixture(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotFound, oopsCode(t, err))
}

func TestHireEmployee_AlreadyHired(t *testing.T) {
	f := takeFixture(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
		},
	})
	require.NoError(t, err)
	_, err = empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeAlreadyHired, oopsCode(t, err))
}

func TestHireEmployee_PositionTooShort(t *testing.T) {
	f := takeFixture(t)
	short := "A"
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
			Position:      &short,
		},
	})
	require.Error(t, err)
	assertHasViolation(t, validationViolations(t, err), "position", "min")
}

func TestHireEmployee_PositionTooLong(t *testing.T) {
	f := takeFixture(t)
	long := strings.Repeat("x", 257)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: testUserAliceID,
			DepartmentID:  f.DeptA1a.String(),
			Position:      &long,
		},
	})
	require.Error(t, err)
	assertHasViolation(t, validationViolations(t, err), "position", "max")
}

func TestHireEmployee_MultiErrorReturnsAllViolations(t *testing.T) {
	tooShort := "A"
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		Caller: sysadminCaller,
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: "",
			Position:      &tooShort,
		},
	})
	// ZitadelUserID and DepartmentID are empty strings with a `required`
	// rule; Position="A" trips `min`. The translator collapses all
	// violations into a single validation_failed error carrying the
	// full list in context.
	violations := validationViolations(t, err)
	assertHasViolation(t, violations, "zitadel_user_id", "required")
	assertHasViolation(t, violations, "department_id", "required")
	assertHasViolation(t, violations, "position", "min")
}
