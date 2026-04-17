//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	employeev1 "github.com/medincident/medincident-command-service/pkg/event/employee/v1"
)

func TestUpdateEmployeeDepartment_Success_SameClinic(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptA1b,
	}))
	rows := latestOutbox(t)
	require.Len(t, rows, 1)
	require.Equal(t, membership.SubjectEmployeeDepartmentChanged, rows[0].Subject)
	ev := &employeev1.EmployeeDepartmentChanged{}
	decodePayload(t, rows[0], ev)
	assert.Equal(t, f.DeptA1b.String(), ev.DepartmentId)
}

func TestUpdateEmployeeDepartment_Success_DifferentClinic_SameOrg(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // Alice in DeptA1a (Clinic A1)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptA2a, // in Clinic A2 — same org A
	}))
	require.Len(t, latestOutbox(t), 1)
}

func TestUpdateEmployeeDepartment_DifferentOrganization(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	err := empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptB1a, // different org
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotInSameOrganization, oopsCode(t, err))
}

func TestUpdateEmployeeDepartment_NoOp(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	truncateOutbox(t)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptA1a, // same as current
	}))
	require.Empty(t, latestOutbox(t))
}

func TestUpdateEmployeeDepartment_DepartmentNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	err := empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: uuidMustV7(),
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotFound, oopsCode(t, err))
}
