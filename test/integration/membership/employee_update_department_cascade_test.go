//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/services/membership"
)

func TestUpdateEmployeeDepartment_CascadeRevokesDR_SameClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignDepartmentResponsible(ctxT(t), membership.AssignDepartmentResponsibleCommand{
		DepartmentID: f.DeptA1a, EmployeeID: aliceID,
	}))
	truncateOutbox(t)

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

	rows := latestOutbox(t)
	require.GreaterOrEqual(t, len(rows), 2)
	assert.Equal(t, "medincident.event.employee.v1.department_changed", rows[0].Subject)
	assert.Equal(t, membership.SubjectDepartmentResponsibleRevoked, rows[1].Subject)
}

func TestUpdateEmployeeDepartment_CascadeRevokesCH_CrossClinic(t *testing.T) {
	f := takeFixture(t)
	aliceID := mustParseUUID(t, hireAlice(t, f))
	require.NoError(t, empSvc.AssignClinicHead(ctxT(t), membership.AssignClinicHeadCommand{
		ClinicID: f.ClinicA1, EmployeeID: aliceID,
	}))
	truncateOutbox(t)

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID: aliceID, DepartmentID: f.DeptA2a,
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`, f.ClinicA1, aliceID,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)

	rows := latestOutbox(t)
	require.GreaterOrEqual(t, len(rows), 2)
	assert.Equal(t, "medincident.event.employee.v1.department_changed", rows[0].Subject)
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
	truncateOutbox(t)

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

	rows := latestOutbox(t)
	subjects := make([]string, 0, len(rows))
	for _, r := range rows {
		subjects = append(subjects, r.Subject)
	}
	assert.Contains(t, subjects, membership.SubjectDepartmentResponsibleDeputyRemoved,
		"expected DepartmentResponsibleDeputyRemoved in outbox, got: %v", subjects)
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
	truncateOutbox(t)

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

	rows := latestOutbox(t)
	subjects := make([]string, 0, len(rows))
	for _, r := range rows {
		subjects = append(subjects, r.Subject)
	}
	assert.Contains(t, subjects, membership.SubjectClinicHeadDeputyRemoved,
		"expected ClinicHeadDeputyRemoved in outbox, got: %v", subjects)
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
	truncateOutbox(t)

	require.NoError(t, empSvc.UpdateDepartment(ctxT(t), membership.UpdateEmployeeDepartmentCommand{
		ID:           aliceID,
		DepartmentID: f.DeptA1b,
	}))

	rows := latestOutbox(t)
	require.Len(t, rows, 3, "expected 3 events: DepartmentChanged, DeputyRemoved, Revoked")
	assert.Equal(t, "medincident.event.employee.v1.department_changed", rows[0].Subject, "Rule 1: cause first")
	assert.Equal(t, membership.SubjectDepartmentResponsibleDeputyRemoved, rows[1].Subject, "Rule 2: cleanup before terminate")
	assert.Equal(t, membership.SubjectDepartmentResponsibleRevoked, rows[2].Subject)
}
