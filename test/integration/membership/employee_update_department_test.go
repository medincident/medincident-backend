//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/command/membership"
)

func TestUpdateEmployeeDepartment_Success_SameClinic(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptA1b,
	}))
}

func TestUpdateEmployeeDepartment_Success_DifferentClinic_SameOrg(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // Alice in DeptA1a (Clinic A1)
	id := mustParseUUID(t, empID)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptA2a, // in Clinic A2 — same org A
	}))

	var projDeptID uuid.UUID
	require.NoError(t, testDB.Raw(`SELECT department_id FROM projections.employees WHERE id = ?`, id).Row().Scan(&projDeptID))
	assert.Equal(t, f.DeptA2a, projDeptID)
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
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           id,
		DepartmentID: f.DeptA1a, // same as current
	}))
	// No-op means projection still shows the original department.
	var projDeptID uuid.UUID
	require.NoError(t, testDB.Raw(`SELECT department_id FROM projections.employees WHERE id = ?`, id).Row().Scan(&projDeptID))
	assert.Equal(t, f.DeptA1a, projDeptID)
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
