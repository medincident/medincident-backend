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
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeeDepartmentPayload{
			ID:           id.String(),
			DepartmentID: f.DeptA1b.String(),
		},
	}))
}

func TestUpdateEmployeeDepartment_Success_DifferentClinic_SameOrg(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f) // Alice in DeptA1a (Clinic A1)
	id := mustParseUUID(t, empID)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeeDepartmentPayload{
			ID:           id.String(),
			DepartmentID: f.DeptA2a.String(), // in Clinic A2 — same org A
		},
	}))

	var domainDeptID uuid.UUID
	require.NoError(t, testDB.Raw(`SELECT department_id FROM domain.employees WHERE id = ?`, id).Row().Scan(&domainDeptID))
	assert.Equal(t, f.DeptA2a, domainDeptID)
}

func TestUpdateEmployeeDepartment_DifferentOrganization(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	err := empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeeDepartmentPayload{
			ID:           id.String(),
			DepartmentID: f.DeptB1a.String(), // different org
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotInSameOrganization, oopsCode(t, err))
}

func TestUpdateEmployeeDepartment_NoOp(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeeDepartmentPayload{
			ID:           id.String(),
			DepartmentID: f.DeptA1a.String(), // same as current
		},
	}))
	// No-op means domain still shows the original department.
	var domainDeptID uuid.UUID
	require.NoError(t, testDB.Raw(`SELECT department_id FROM domain.employees WHERE id = ?`, id).Row().Scan(&domainDeptID))
	assert.Equal(t, f.DeptA1a, domainDeptID)
}

func TestUpdateEmployeeDepartment_DepartmentNotFound(t *testing.T) {
	f := takeFixture(t)
	empID := hireAlice(t, f)
	id := mustParseUUID(t, empID)
	err := empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		Caller: sysadminCaller,
		Payload: membership.UpdateEmployeeDepartmentPayload{
			ID:           id.String(),
			DepartmentID: uuidMustV7().String(),
		},
	})
	require.Error(t, err)
	assert.Equal(t, membership.ErrCodeDepartmentNotFound, oopsCode(t, err))
}
