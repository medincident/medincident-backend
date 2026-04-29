//go:build integration

package analytics_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
)

func TestGetSummary_BasicCounts(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "User A")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedOrgAdmin(t, empID, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, caller, deptID, categoryID, typeID)
	createIncident(t, caller, deptID, categoryID, typeID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(2), result.IncidentAgg.Total)
	assert.Equal(t, int64(2), result.IncidentAgg.Pending)
	assert.Equal(t, int64(0), result.IncidentAgg.Done)
	assert.False(t, result.IncidentPercentiles.P50.Valid)
}

func TestGetSummary_TopCategories(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	catA := seedCategoryNamed(t, orgID, "Категория А")
	catB := seedCategoryNamed(t, orgID, "Категория Б")
	typeA := seedIncidentTypeNamed(t, catA, "Тип А")
	typeB := seedIncidentTypeNamed(t, catB, "Тип Б")
	seedUser(t, userAZitadelID, "User A")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedOrgAdmin(t, empID, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	// 3 incidents in catA, 1 in catB
	createIncident(t, caller, deptID, catA, typeA)
	createIncident(t, caller, deptID, catA, typeA)
	createIncident(t, caller, deptID, catA, typeA)
	createIncident(t, caller, deptID, catB, typeB)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.IncidentCategories)

	assert.Equal(t, catA, result.IncidentCategories[0].ID)
	assert.Equal(t, int64(3), result.IncidentCategories[0].Count)
}

func TestGetSummary_PatientBufferRates(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)

	raw, err := testDB.DB()
	require.NoError(t, err)

	for range 2 {
		_, err = raw.Exec(
			`INSERT INTO domain.patient_incident_buffer (id, organization_id, patient_zitadel_user_id, status, created_at)
			 VALUES ($1, $2, $3, 'published', now())`,
			uuidNew(), orgID, uuidNew(),
		)
		require.NoError(t, err)
	}
	_, err = raw.Exec(
		`INSERT INTO domain.patient_incident_buffer (id, organization_id, patient_zitadel_user_id, status, created_at)
		 VALUES ($1, $2, $3, 'rejected', now())`,
		uuidNew(), orgID, uuidNew(),
	)
	require.NoError(t, err)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, sysadminCaller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(3), result.BufferAgg.Total)
	assert.Equal(t, int64(2), result.BufferAgg.Published)
	assert.Equal(t, int64(1), result.BufferAgg.Rejected)
}

func TestGetSummary_DeptResponsible_CanReadDeptScope(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)

	seedUser(t, userBZitadelID, "Org Admin")
	orgAdminEmpID := seedEmployee(t, userBZitadelID, orgID, deptID)
	seedOrgAdmin(t, orgAdminEmpID, orgID)
	orgAdminCaller := authz.Caller{ZitadelUserID: userBZitadelID}

	seedUser(t, userAZitadelID, "Dept Responsible")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedDeptResponsible(t, empID, deptID)
	deptResponsibleCaller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, orgAdminCaller, deptID, categoryID, typeID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	deptIDStr := deptID.String()

	result, err := analyticsRdr.GetSummary(ctx, deptResponsibleCaller, orgID.String(), from, to, nil, &deptIDStr)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(1), result.IncidentAgg.Total)
}

func TestGetSummary_PermissionDenied_ForNonManager(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	seedUser(t, userAZitadelID, "Regular User")
	seedEmployee(t, userAZitadelID, orgID, deptID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	_, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.Error(t, err)
}

func TestGetSummary_ClinicHead_CanReadClinicScope(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "Clinic Head")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedClinicHead(t, empID, clinicID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, caller, deptID, categoryID, typeID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	clinicStr := clinicID.String()

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, &clinicStr, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(1), result.IncidentAgg.Total)
}

func TestGetSummary_ClinicHead_CannotReadOrgScope(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	seedUser(t, userAZitadelID, "Clinic Head")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedClinicHead(t, empID, clinicID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	_, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.Error(t, err)
}

func TestGetSummary_CrossTenantClinicRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgA := seedOrg(t)
	orgB := seedOrg(t)
	clinicB := seedClinic(t, orgB)

	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	clinicBStr := clinicB.String()

	_, err := analyticsRdr.GetSummary(ctx, sysadminCaller, orgA.String(), from, to, &clinicBStr, nil)
	require.Error(t, err)
}

func TestGetSummary_CrossTenantDeptRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgA := seedOrg(t)
	orgB := seedOrg(t)
	clinicB := seedClinic(t, orgB)
	deptB := seedDept(t, clinicB)

	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	deptBStr := deptB.String()

	_, err := analyticsRdr.GetSummary(ctx, sysadminCaller, orgA.String(), from, to, nil, &deptBStr)
	require.Error(t, err)
}

func TestGetSummary_TopTypes(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	catID := seedCategory(t, orgID)
	typeA := seedIncidentTypeNamed(t, catID, "Тип А")
	typeB := seedIncidentTypeNamed(t, catID, "Тип Б")
	seedUser(t, userAZitadelID, "User A")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedOrgAdmin(t, empID, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, caller, deptID, catID, typeA)
	createIncident(t, caller, deptID, catID, typeA)
	createIncident(t, caller, deptID, catID, typeA)
	createIncident(t, caller, deptID, catID, typeB)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, result.IncidentTypes)

	assert.Equal(t, typeA, result.IncidentTypes[0].ID)
	assert.Equal(t, int64(3), result.IncidentTypes[0].Count)
}

func TestGetSummary_TopDepartments(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptA := seedDept(t, clinicID)
	deptB := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "User A")
	empID := seedEmployee(t, userAZitadelID, orgID, deptA)
	seedOrgAdmin(t, empID, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, caller, deptA, categoryID, typeID)
	createIncident(t, caller, deptA, categoryID, typeID)
	createIncident(t, caller, deptA, categoryID, typeID)
	createIncident(t, caller, deptB, categoryID, typeID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, result.IncidentDepartments)

	assert.Equal(t, deptA, result.IncidentDepartments[0].ID)
	assert.Equal(t, int64(3), result.IncidentDepartments[0].Count)
}

func TestGetSummary_PriorityBreakdown(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "User A")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedOrgAdmin(t, empID, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	inc1 := createIncident(t, caller, deptID, categoryID, typeID)
	createIncident(t, caller, deptID, categoryID, typeID)

	require.NoError(t, incidentSvc.UpdatePriority(context.Background(), incidentsvc.UpdateIncidentPriorityCommand{
		Caller:  caller,
		Payload: incidentsvc.UpdateIncidentPriorityPayload{IncidentID: inc1.String(), Priority: "high"},
	}))

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(1), result.IncidentAgg.Normal)
	assert.Equal(t, int64(1), result.IncidentAgg.High)
}

func TestGetSummary_ResolutionStats_PopulatedWhenDoneExists(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "User A")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedOrgAdmin(t, empID, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	incID := createIncident(t, caller, deptID, categoryID, typeID)
	doneIncident(t, caller, incID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(1), result.IncidentAgg.Done)
	assert.Equal(t, int64(0), result.IncidentAgg.Pending)
	assert.True(t, result.IncidentPercentiles.P50.Valid)
}

func TestGetSummary_RequestCounts(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	requestTypeID := seedRequestType(t, orgID)
	seedUser(t, userAZitadelID, "User A")
	empA := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedOrgAdmin(t, empA, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createRequest(t, caller, deptID, requestTypeID, empA)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSummary(ctx, caller, orgID.String(), from, to, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(1), result.RequestAgg.Total)
	assert.Equal(t, int64(1), result.RequestAgg.Created)
}
