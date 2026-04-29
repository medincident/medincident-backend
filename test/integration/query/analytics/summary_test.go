//go:build integration

package analytics_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
)

func seedCategoryNamed(t *testing.T, orgID uuid.UUID, name string) uuid.UUID {
	t.Helper()
	res, err := categorySvc.Create(context.Background(), classifiersvc.CreateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{OrganizationID: orgID.String(), Name: name},
	})
	require.NoError(t, err)
	return res.ID
}

func seedIncidentTypeNamed(t *testing.T, categoryID uuid.UUID, name string) uuid.UUID {
	t.Helper()
	res, err := typeSvc.Create(context.Background(), classifiersvc.CreateIncidentTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{CategoryID: categoryID.String(), Name: name},
	})
	require.NoError(t, err)
	return res.ID
}

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

	// Insert 2 published + 1 rejected directly into domain.patient_incident_buffer
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
	seedUser(t, userAZitadelID, "Dept Responsible")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedDeptResponsible(t, empID, deptID)
	deptResponsibleCaller := authz.Caller{ZitadelUserID: userAZitadelID}

	// create incident as dept responsible caller (who is an org member)
	createIncident(t, deptResponsibleCaller, deptID, categoryID, typeID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	deptIDStr := deptID.String()

	result, err := analyticsRdr.GetSummary(ctx, deptResponsibleCaller, orgID.String(), from, to, nil, &deptIDStr)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(1), result.IncidentAgg.Total)
}
