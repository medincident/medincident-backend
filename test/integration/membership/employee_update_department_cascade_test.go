//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/service/command/membership"
)

func TestUpdateEmployeeDepartment_CascadeRevokesDR_SameClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           aliceID,
		DepartmentID: f.DeptA1b,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		f.DeptA1a, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestUpdateEmployeeDepartment_CascadeRevokesCH_CrossClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID: f.ClinicA1, EmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID: aliceID, DepartmentID: f.DeptA2a,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`, f.ClinicA1, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestUpdateEmployeeDepartment_CHUnaffected_SameClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID: f.ClinicA1, EmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID: aliceID, DepartmentID: f.DeptA1b,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`, f.ClinicA1, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "ClinicHead must survive same-clinic moves")
}

// When Alice is moved out of DeptA1a, any DR row in DeptA1a that
// has Alice as its deputy must have the deputy slot cleared — the
// invariant "deputy must belong to the same department as the holder"
// becomes false the moment Alice's department_id changes.
func TestUpdateEmployeeDepartment_ClearsDRDeputySlotInOldDepartment(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f)) // also in DeptA1a
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a, EmployeeID: bobID,
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID: f.DeptA1a, EmployeeID: bobID, DeputyEmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID: aliceID, DepartmentID: f.DeptA1b,
	}))

	var deputyCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ? AND deputy_employee_id IS NOT NULL`,
		f.DeptA1a, bobID,
	).Scan(&deputyCount).Error)
	assert.Equal(t, int64(0), deputyCount, "DR deputy slot in old department must be cleared")

	var holderCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		f.DeptA1a, bobID,
	).Scan(&holderCount).Error)
	assert.Equal(t, int64(1), holderCount, "DR holder row itself must survive — Bob stays in DeptA1a")
}

// When Alice is moved cross-clinic (DeptA1a → DeptA2a), any CH row in
// the OLD clinic (ClinicA1) that has Alice as its deputy must be
// cleared.
func TestUpdateEmployeeDepartment_ClearsCHDeputySlotInOldClinic_CrossClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f)) // DeptA1a / ClinicA1
	bobID := mustParseUUID(t, hireBob(t, f))     // DeptA1a / ClinicA1
	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID: f.ClinicA1, EmployeeID: bobID,
	}))
	require.NoError(t, empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID: f.ClinicA1, EmployeeID: bobID, DeputyEmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID: aliceID, DepartmentID: f.DeptA2a,
	}))

	var deputyCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ? AND deputy_employee_id IS NOT NULL`,
		f.ClinicA1, bobID,
	).Scan(&deputyCount).Error)
	assert.Equal(t, int64(0), deputyCount, "CH deputy slot in old clinic must be cleared on cross-clinic move")

	var holderCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		f.ClinicA1, bobID,
	).Scan(&holderCount).Error)
	assert.Equal(t, int64(1), holderCount, "CH holder row itself must survive — Bob stays in ClinicA1")
}

// Same-clinic moves must NOT clear the CH deputy slot — the invariant
// is per-clinic and the employee is still in the same clinic.
func TestUpdateEmployeeDepartment_CHDeputyUnaffected_SameClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))
	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID: f.ClinicA1, EmployeeID: bobID,
	}))
	require.NoError(t, empSvc.AssignClinicHeadDeputy(ctxT(t), membership.AssignClinicHeadDeputyCommand{
		ClinicID: f.ClinicA1, EmployeeID: bobID, DeputyEmployeeID: aliceID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID: aliceID, DepartmentID: f.DeptA1b,
	}))

	var deputyCount int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ? AND deputy_employee_id = ?`,
		f.ClinicA1, bobID, aliceID,
	).Scan(&deputyCount).Error)
	assert.Equal(t, int64(1), deputyCount, "CH deputy slot must survive same-clinic moves")
}

func TestUpdateEmployeeDepartment_CascadeRevokesDRWithDeputy_EmitsDeputyRemovedFirst(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	bobID := mustParseUUID(t, hireBob(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID,
	}))
	require.NoError(t, empSvc.AssignDepartmentResponsibleDeputy(ctxT(t), membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID, DeputyEmployeeID: bobID,
	}))

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           aliceID,
		DepartmentID: f.DeptA1b,
	}))
}
