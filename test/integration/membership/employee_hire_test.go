//go:build integration

package membership_integration_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

func TestHireEmployee_Success_WithPosition(t *testing.T) {
	f := takeFixture(t)
	truncateOutbox(t)
	pos := "Senior nurse"
	res, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
		Position:      &pos,
	})
	require.NoError(t, err)
	require.NotEqual(t, "", res.ID.String())

	var count int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employees WHERE id = ?`, res.ID).Scan(&count).Error)
	require.Equal(t, int64(1), count)

	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectEmployeeHired, rows[0].Subject)
	ev := &employeev1.EmployeeHired{}
	env := decodePayload(t, rows[0], ev)
	assert.Equal(t, "employee", env.AggregateType)
	assert.Equal(t, res.ID.String(), env.AggregateId)
	assert.Equal(t, testUserAliceID, ev.ZitadelUserId)
	assert.Equal(t, f.DeptA1a.String(), ev.DepartmentId)
	assert.Equal(t, f.OrgA.String(), ev.OrganizationId)
	require.NotNil(t, ev.Position)
	assert.Equal(t, "Senior nurse", *ev.Position)
}

func TestHireEmployee_Success_NoPosition(t *testing.T) {
	f := takeFixture(t)
	truncateOutbox(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserBobID,
		DepartmentID:  f.DeptA1a,
	})
	require.NoError(t, err)
	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	ev := &employeev1.EmployeeHired{}
	decodePayload(t, rows[0], ev)
	assert.Nil(t, ev.Position)
}

func TestHireEmployee_Success_EmptyPositionTreatedAsUnset(t *testing.T) {
	f := takeFixture(t)
	empty := "   "
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
		Position:      &empty,
	})
	require.NoError(t, err)
	var position *string
	require.NoError(t, testDB.Raw(`SELECT position FROM domain.employees`).Scan(&position).Error)
	assert.Nil(t, position)
}

func TestHireEmployee_ZitadelUserNotFound(t *testing.T) {
	f := takeFixture(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: "nonexistent-user-id",
		DepartmentID:  f.DeptA1a,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeZitadelUserNotFound, oopsCode(t, err))
}

func TestHireEmployee_DepartmentNotFound(t *testing.T) {
	_ = takeFixture(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotFound, oopsCode(t, err))
}

func TestHireEmployee_AlreadyHired(t *testing.T) {
	f := takeFixture(t)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
	})
	require.NoError(t, err)
	_, err = empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeEmployeeAlreadyHired, oopsCode(t, err))
}

func TestHireEmployee_PositionTooShort(t *testing.T) {
	f := takeFixture(t)
	short := "A"
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
		Position:      &short,
	})
	require.Error(t, err)
	assert.Contains(t, oopsCode(t, err), "position_too_short")
}

func TestHireEmployee_PositionTooLong(t *testing.T) {
	f := takeFixture(t)
	long := strings.Repeat("x", 257)
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
		Position:      &long,
	})
	require.Error(t, err)
	assert.Contains(t, oopsCode(t, err), "position_too_long")
}

func TestHireEmployee_MultiErrorReturnsAllViolations(t *testing.T) {
	tooShort := "A"
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: "",
		Position:      &tooShort,
	})
	codes := oopsCodes(t, err)
	assert.Contains(t, codes, membership.ErrCodeEmployeeZitadelUserIDEmpty)
	assert.Contains(t, codes, membership.ErrCodeEmployeeDepartmentIDEmpty)
	assert.Contains(t, codes, membership.ErrCodeEmployeePositionTooShort)
}
