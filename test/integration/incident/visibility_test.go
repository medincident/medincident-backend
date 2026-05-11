//go:build integration

package incident_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
)

// domainIncidentExists returns true if the given incident ID exists in domain.incidents.
func domainIncidentExists(t *testing.T, incidentID uuid.UUID) bool {
	t.Helper()
	var n int
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.incidents WHERE id = ?`, incidentID,
	).Scan(&n).Error)
	return n > 0
}

// domainIncidentStatus returns the status of the given incident from domain.incidents.
func domainIncidentStatus(t *testing.T, incidentID uuid.UUID) string {
	t.Helper()
	var status string
	require.NoError(t, testDB.Raw(
		`SELECT status FROM domain.incidents WHERE id = ?`, incidentID,
	).Row().Scan(&status))
	return status
}

// domainIncidentRegistrar returns the zitadel_user_id of the registrar employee
// for the given incident (via the employees join).
func domainIncidentRegistrar(t *testing.T, incidentID uuid.UUID) string {
	t.Helper()
	var zitadelID string
	require.NoError(t, testDB.Raw(
		`SELECT e.zitadel_user_id
		 FROM domain.incidents i
		 JOIN domain.employees e ON e.id = i.registrar_employee_id
		 WHERE i.id = ?`, incidentID,
	).Row().Scan(&zitadelID))
	return zitadelID
}

// TestVisibility_RoleMatrix tests that incidents are created with the correct
// domain state. One org, two clinics (A, B), two departments (A1, B1).
//
//	I1: clinic A / dept A1, active (registrar X)
//	I2: clinic A / dept A1, cancelled (registrar X)
//	I3: clinic B / dept B1, active (registrar Y)
func TestVisibility_RoleMatrix(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	// Structure: one org, two clinics, two departments.
	orgID := seedOrg(t)
	clinicA := seedClinic(t, orgID)
	clinicB := seedClinic(t, orgID)
	deptA1 := seedDept(t, clinicA)
	deptB1 := seedDept(t, clinicB)

	catID := seedCategory(t, orgID)
	typeID := seedType(t, catID, false)

	// Registrar X — employee in dept A1.
	const registrarXID = "reg-x-zitadel"
	seedUser(t, registrarXID, "Регистратор X")
	seedEmployee(t, registrarXID, orgID, deptA1)
	callerX := authz.Caller{ZitadelUserID: registrarXID}

	// Registrar Y — employee in dept B1.
	const registrarYID = "reg-y-zitadel"
	seedUser(t, registrarYID, "Регистратор Y")
	seedEmployee(t, registrarYID, orgID, deptB1)
	callerY := authz.Caller{ZitadelUserID: registrarYID}

	// I1: active, registrar X, clinic A / dept A1.
	i1ID := createIncident(t, ctx, callerX, deptA1, catID, typeID)

	// I2: cancelled, registrar X, clinic A / dept A1.
	i2ID := createIncident(t, ctx, callerX, deptA1, catID, typeID)
	require.NoError(t, incidentSvc.Cancel(ctx, incidentsvc.CancelIncidentCommand{
		Caller:  callerX,
		Payload: incidentsvc.CancelIncidentPayload{IncidentID: i2ID.String()},
	}))

	// I3: active, registrar Y, clinic B / dept B1.
	i3ID := createIncident(t, ctx, callerY, deptB1, catID, typeID)

	// Verify domain state for all three incidents.
	assert.True(t, domainIncidentExists(t, i1ID), "I1 must exist in domain")
	assert.True(t, domainIncidentExists(t, i2ID), "I2 must exist in domain")
	assert.True(t, domainIncidentExists(t, i3ID), "I3 must exist in domain")

	assert.Equal(t, "pending", domainIncidentStatus(t, i1ID))
	assert.Equal(t, "cancelled", domainIncidentStatus(t, i2ID))
	assert.Equal(t, "pending", domainIncidentStatus(t, i3ID))

	assert.Equal(t, registrarXID, domainIncidentRegistrar(t, i1ID))
	assert.Equal(t, registrarXID, domainIncidentRegistrar(t, i2ID))
	assert.Equal(t, registrarYID, domainIncidentRegistrar(t, i3ID))

	// Seed role holders (used by cascade / authz tests elsewhere; seeded here
	// to ensure the suite stays internally consistent).
	const orgAdminID = "orgadmin-vis-zitadel"
	seedUser(t, orgAdminID, "Орг Администратор")
	orgAdminEmpID := seedEmployee(t, orgAdminID, orgID, deptA1)
	seedOrgAdmin(t, orgAdminEmpID, orgID)

	const orgHeadID = "orghead-vis-zitadel"
	seedUser(t, orgHeadID, "Руководитель Орг")
	orgHeadEmpID := seedEmployee(t, orgHeadID, orgID, deptA1)
	seedOrgHead(t, orgHeadEmpID, orgID)

	const clinicHeadAID = "clinicheadA-vis-zitadel"
	seedUser(t, clinicHeadAID, "Главврач Клиника А")
	clinicHeadAEmpID := seedEmployee(t, clinicHeadAID, orgID, deptA1)
	seedClinicHead(t, clinicHeadAEmpID, clinicA)

	const deptRespA1ID = "deptrespA1-vis-zitadel"
	seedUser(t, deptRespA1ID, "Ответственный Отдела А1")
	deptRespA1EmpID := seedEmployee(t, deptRespA1ID, orgID, deptA1)
	seedDeptResponsible(t, deptRespA1EmpID, deptA1)

	// All role rows exist in domain.
	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_admins WHERE organization_id = ? AND employee_id = ?`,
		orgID, orgAdminEmpID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "OrgAdmin row must exist")

	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.org_heads WHERE organization_id = ? AND employee_id = ?`,
		orgID, orgHeadEmpID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "OrgHead row must exist")

	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		clinicA, clinicHeadAEmpID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "ClinicHead row must exist")

	require.NoError(t, testDB.Raw(
		`SELECT count(*) FROM domain.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		deptA1, deptRespA1EmpID,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count, "DeptResponsible row must exist")
}

// TestVisibility_GetIncidentEnforcesVisibility verifies domain incident state:
// incidents are correctly linked to their registrar and department.
func TestVisibility_GetIncidentEnforcesVisibility(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicA := seedClinic(t, orgID)
	clinicB := seedClinic(t, orgID)
	deptA1 := seedDept(t, clinicA)
	deptB1 := seedDept(t, clinicB)

	catID := seedCategory(t, orgID)
	typeID := seedType(t, catID, false)

	// Registrar of clinic B.
	const regBID = "reg-b-zitadel"
	seedUser(t, regBID, "Регистратор B")
	seedEmployee(t, regBID, orgID, deptB1)

	// Incident in clinic A / dept A1 (different registrar).
	const regAID = "reg-a2-zitadel"
	seedUser(t, regAID, "Регистратор A2")
	seedEmployee(t, regAID, orgID, deptA1)
	callerA := authz.Caller{ZitadelUserID: regAID}

	incID := createIncident(t, ctx, callerA, deptA1, catID, typeID)

	// Incident must be in domain with the correct registrar.
	assert.True(t, domainIncidentExists(t, incID))
	assert.Equal(t, regAID, domainIncidentRegistrar(t, incID))

	// Department of the incident must be deptA1 (not deptB1).
	var deptID uuid.UUID
	require.NoError(t, testDB.Raw(
		`SELECT department_id FROM domain.incidents WHERE id = ?`, incID,
	).Row().Scan(&deptID))
	assert.Equal(t, deptA1, deptID)
	assert.NotEqual(t, deptB1, deptID)
}
