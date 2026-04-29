//go:build integration

package analytics_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

func TestGetSnapshot_ReturnsIncidents(t *testing.T) {
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

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	result, err := analyticsRdr.GetSnapshot(ctx, caller, orgID.String(), from, to, nil, nil, false)
	require.NoError(t, err)
	assert.Len(t, result.Incidents, 1)
	assert.Equal(t, "pending", result.Incidents[0].Status)
	assert.Empty(t, result.Requests)
	assert.Empty(t, result.PatientBuffer)
}

func TestGetSnapshot_ClinicFilter_ReturnsOnlyClinicIncidents(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicA := seedClinic(t, orgID)
	clinicB := seedClinic(t, orgID)
	deptA := seedDept(t, clinicA)
	deptB := seedDept(t, clinicB)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "User A")
	empA := seedEmployee(t, userAZitadelID, orgID, deptA)
	seedOrgAdmin(t, empA, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, caller, deptA, categoryID, typeID)
	createIncident(t, caller, deptB, categoryID, typeID)

	from := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	clinicAStr := clinicA.String()

	result, err := analyticsRdr.GetSnapshot(ctx, caller, orgID.String(), from, to, &clinicAStr, nil, false)
	require.NoError(t, err)
	assert.Len(t, result.Incidents, 1)
	assert.Equal(t, clinicA, result.Incidents[0].ClinicID)
}

func TestGetSnapshot_CrossTenantClinicRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgA := seedOrg(t)
	orgB := seedOrg(t)
	clinicB := seedClinic(t, orgB)

	from := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	to := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	clinicBStr := clinicB.String()

	result, err := analyticsRdr.GetSnapshot(ctx, sysadminCaller, orgA.String(), from, to, &clinicBStr, nil, false)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetSnapshot_PeriodTooLarge_Rejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	from := time.Now().Add(-400 * 24 * time.Hour).Format(time.RFC3339)
	to := time.Now().Format(time.RFC3339)

	result, err := analyticsRdr.GetSnapshot(ctx, sysadminCaller, orgID.String(), from, to, nil, nil, false)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetSnapshot_PermissionDenied_ForNonManager(t *testing.T) {
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

	_, err := analyticsRdr.GetSnapshot(ctx, caller, orgID.String(), from, to, nil, nil, false)
	require.Error(t, err)
}

func TestGetSnapshot_ClinicHead_CanReadClinicScope(t *testing.T) {
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

	result, err := analyticsRdr.GetSnapshot(ctx, caller, orgID.String(), from, to, &clinicStr, nil, false)
	require.NoError(t, err)
	assert.Len(t, result.Incidents, 1)
}

func TestGetSnapshot_ClinicHead_CannotReadOrgScope(t *testing.T) {
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

	_, err := analyticsRdr.GetSnapshot(ctx, caller, orgID.String(), from, to, nil, nil, false)
	require.Error(t, err)
}
