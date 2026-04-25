//go:build integration

package incident_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	queryincident "github.com/medincident/medincident-backend/internal/service/query/incident"
)

// incidentIDSet converts a slice of IncidentView to a set of IDs.
func incidentIDSet(views []queryincident.IncidentView) map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool, len(views))
	for _, v := range views {
		m[v.ID] = true
	}
	return m
}

// TestVisibility_RoleMatrix tests the visibility matrix for ListIncidents:
//
//	One org, two clinics (A, B), two departments (A1, B1).
//	I1: clinic A / dept A1, active (registrar X)
//	I2: clinic A / dept A1, cancelled (registrar X)
//	I3: clinic B / dept B1, active (registrar Y)
//
//	SystemAdmin     → I1, I2, I3
//	OrgAdmin        → I1, I2, I3
//	OrgHead         → I1, I3 (no cancelled)
//	ClinicHead(A)   → I1 (clinic A, no cancelled)
//	DeptResp(A1)    → I1 (dept A1, no cancelled)
//	Registrar X     → I1, I2 (own incidents including own cancelled)
//	Patient (no submissions) → []
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

	// -- Seed the role callers. --

	// OrgAdmin.
	const orgAdminID = "orgadmin-vis-zitadel"
	seedUser(t, orgAdminID, "Орг Администратор")
	orgAdminEmpID := seedEmployee(t, orgAdminID, orgID, deptA1)
	seedOrgAdmin(t, orgAdminEmpID, orgID)

	// OrgHead.
	const orgHeadID = "orghead-vis-zitadel"
	seedUser(t, orgHeadID, "Руководитель Орг")
	orgHeadEmpID := seedEmployee(t, orgHeadID, orgID, deptA1)
	seedOrgHead(t, orgHeadEmpID, orgID)

	// ClinicHead of clinic A.
	const clinicHeadAID = "clinicheadA-vis-zitadel"
	seedUser(t, clinicHeadAID, "Главврач Клиника А")
	clinicHeadAEmpID := seedEmployee(t, clinicHeadAID, orgID, deptA1)
	seedClinicHead(t, clinicHeadAEmpID, clinicA)

	// DeptResponsible of dept A1.
	const deptRespA1ID = "deptrespA1-vis-zitadel"
	seedUser(t, deptRespA1ID, "Ответственный Отдела А1")
	deptRespA1EmpID := seedEmployee(t, deptRespA1ID, orgID, deptA1)
	seedDeptResponsible(t, deptRespA1EmpID, deptA1)

	// Patient with no submissions.
	const patientID = "patient-vis-zitadel"
	// Note: patient has no employee row and no incidents.

	filters := &queryincident.ListFilters{}

	// SystemAdmin sees all three.
	views, err := incidentRdr.ListIncidents(ctx, sysadminZitadelID, orgID, filters)
	require.NoError(t, err)
	set := incidentIDSet(views)
	assert.True(t, set[i1ID], "SystemAdmin should see I1")
	assert.True(t, set[i2ID], "SystemAdmin should see I2")
	assert.True(t, set[i3ID], "SystemAdmin should see I3")

	// OrgAdmin sees all three.
	views, err = incidentRdr.ListIncidents(ctx, orgAdminID, orgID, filters)
	require.NoError(t, err)
	set = incidentIDSet(views)
	assert.True(t, set[i1ID], "OrgAdmin should see I1")
	assert.True(t, set[i2ID], "OrgAdmin should see I2")
	assert.True(t, set[i3ID], "OrgAdmin should see I3")

	// OrgHead sees I1 and I3 (no cancelled).
	views, err = incidentRdr.ListIncidents(ctx, orgHeadID, orgID, filters)
	require.NoError(t, err)
	set = incidentIDSet(views)
	assert.True(t, set[i1ID], "OrgHead should see I1")
	assert.False(t, set[i2ID], "OrgHead should NOT see cancelled I2")
	assert.True(t, set[i3ID], "OrgHead should see I3")

	// ClinicHead of clinic A sees I1 only (no cancelled).
	views, err = incidentRdr.ListIncidents(ctx, clinicHeadAID, orgID, filters)
	require.NoError(t, err)
	set = incidentIDSet(views)
	assert.True(t, set[i1ID], "ClinicHead(A) should see I1")
	assert.False(t, set[i2ID], "ClinicHead(A) should NOT see cancelled I2")
	assert.False(t, set[i3ID], "ClinicHead(A) should NOT see I3 (different clinic)")

	// DeptResponsible of dept A1 sees I1 only (no cancelled).
	views, err = incidentRdr.ListIncidents(ctx, deptRespA1ID, orgID, filters)
	require.NoError(t, err)
	set = incidentIDSet(views)
	assert.True(t, set[i1ID], "DeptResp(A1) should see I1")
	assert.False(t, set[i2ID], "DeptResp(A1) should NOT see cancelled I2")
	assert.False(t, set[i3ID], "DeptResp(A1) should NOT see I3 (different dept)")

	// Registrar X sees I1 and I2 (own incidents, including own cancelled).
	views, err = incidentRdr.ListIncidents(ctx, registrarXID, orgID, filters)
	require.NoError(t, err)
	set = incidentIDSet(views)
	assert.True(t, set[i1ID], "Registrar X should see own I1")
	assert.True(t, set[i2ID], "Registrar X should see own cancelled I2")
	assert.False(t, set[i3ID], "Registrar X should NOT see I3 (different registrar)")

	// Patient with no submissions sees nothing.
	views, err = incidentRdr.ListIncidents(ctx, patientID, orgID, filters)
	require.NoError(t, err)
	assert.Empty(t, views, "Patient with no submissions should see no incidents")
}

// TestVisibility_GetIncidentEnforcesVisibility verifies that GetIncident
// respects the same visibility rules: a caller who cannot list an incident
// also cannot fetch it directly.
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

	// Caller B (employee in dept B1) should not see incident from dept A1.
	_, err := incidentRdr.GetIncident(ctx, regBID, incID)
	require.Error(t, err)
	assert.Equal(t, queryincident.ErrCodeIncidentNotFound, codeOf(t, err))

	// Caller A can see own incident.
	view, err := incidentRdr.GetIncident(ctx, regAID, incID)
	require.NoError(t, err)
	assert.Equal(t, model.IncidentStatusPending, view.Status)
}
