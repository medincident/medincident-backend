//go:build integration

package analytics_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/service/authz"
	analyticsread "github.com/medincident/medincident-backend/internal/service/query/analytics"
)

func TestGetTimeSeries_DayGranularity_ZeroFilled(t *testing.T) {
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

	// 3-day window: from 2 days ago truncated to day, to tomorrow truncated to day
	now := time.Now().UTC()
	from := now.AddDate(0, 0, -2).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)

	buckets, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, nil, nil, analyticsread.GranularityDay)
	require.NoError(t, err)
	assert.NotEmpty(t, buckets)

	var hasIncidents bool
	for _, b := range buckets {
		if b.IncidentTotal > 0 {
			hasIncidents = true
			break
		}
	}
	assert.True(t, hasIncidents, "expected at least one bucket with IncidentTotal > 0")
}

func TestGetTimeSeries_EmptyPeriod_ReturnsZeroBuckets(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -2).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)

	buckets, err := analyticsRdr.GetTimeSeries(ctx, sysadminCaller, orgID.String(), from, to, nil, nil, analyticsread.GranularityDay)
	require.NoError(t, err)
	assert.NotEmpty(t, buckets)

	for _, b := range buckets {
		assert.Equal(t, int64(0), b.IncidentTotal, "expected zero IncidentTotal in bucket %v", b.BucketStart)
		assert.Equal(t, int64(0), b.ReqTotal, "expected zero ReqTotal in bucket %v", b.BucketStart)
	}
}

func TestGetTimeSeries_WeekGranularity(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -14).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.Truncate(24 * time.Hour).Format(time.RFC3339)

	buckets, err := analyticsRdr.GetTimeSeries(ctx, sysadminCaller, orgID.String(), from, to, nil, nil, analyticsread.GranularityWeek)
	require.NoError(t, err)
	assert.NotEmpty(t, buckets)
	assert.GreaterOrEqual(t, len(buckets), 2)
}

func TestGetTimeSeries_CrossTenantDeptRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgA := seedOrg(t)
	orgB := seedOrg(t)
	clinicB := seedClinic(t, orgB)
	deptB := seedDept(t, clinicB)

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -1).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)
	deptBStr := deptB.String()

	_, err := analyticsRdr.GetTimeSeries(ctx, sysadminCaller, orgA.String(), from, to, nil, &deptBStr, analyticsread.GranularityDay)
	require.Error(t, err)
}

func TestGetTimeSeries_MonthGranularity(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -60).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.Truncate(24 * time.Hour).Format(time.RFC3339)

	buckets, err := analyticsRdr.GetTimeSeries(ctx, sysadminCaller, orgID.String(), from, to, nil, nil, analyticsread.GranularityMonth)
	require.NoError(t, err)
	assert.NotEmpty(t, buckets)
	assert.GreaterOrEqual(t, len(buckets), 2)
}

func TestGetTimeSeries_ClinicFilter_ReturnsOnlyClinicData(t *testing.T) {
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

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -2).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)
	clinicAStr := clinicA.String()

	buckets, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, &clinicAStr, nil, analyticsread.GranularityDay)
	require.NoError(t, err)

	var total int64
	for _, b := range buckets {
		total += b.IncidentTotal
	}
	assert.Equal(t, int64(1), total)
}

func TestGetTimeSeries_DeptFilter_ReturnsOnlyDeptData(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptA := seedDept(t, clinicID)
	deptB := seedDept(t, clinicID)
	categoryID := seedCategory(t, orgID)
	typeID := seedIncidentType(t, categoryID)
	seedUser(t, userAZitadelID, "User A")
	empA := seedEmployee(t, userAZitadelID, orgID, deptA)
	seedOrgAdmin(t, empA, orgID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	createIncident(t, caller, deptA, categoryID, typeID)
	createIncident(t, caller, deptA, categoryID, typeID)
	createIncident(t, caller, deptB, categoryID, typeID)

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -2).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)
	deptAStr := deptA.String()

	buckets, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, nil, &deptAStr, analyticsread.GranularityDay)
	require.NoError(t, err)

	var total int64
	for _, b := range buckets {
		total += b.IncidentTotal
	}
	assert.Equal(t, int64(2), total)
}

func TestGetTimeSeries_PermissionDenied_ForNonManager(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	seedUser(t, userAZitadelID, "Regular User")
	seedEmployee(t, userAZitadelID, orgID, deptID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -1).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)

	_, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, nil, nil, analyticsread.GranularityDay)
	require.Error(t, err)
}

func TestGetTimeSeries_ClinicHead_CanReadClinicScope(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	seedUser(t, userAZitadelID, "Clinic Head")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedClinicHead(t, empID, clinicID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -2).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)
	clinicStr := clinicID.String()

	_, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, &clinicStr, nil, analyticsread.GranularityDay)
	require.NoError(t, err)
}

func TestGetTimeSeries_ClinicHead_CannotReadOrgScope(t *testing.T) {
	resetDB(t)
	ctx := context.Background()

	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	seedUser(t, userAZitadelID, "Clinic Head")
	empID := seedEmployee(t, userAZitadelID, orgID, deptID)
	seedClinicHead(t, empID, clinicID)
	caller := authz.Caller{ZitadelUserID: userAZitadelID}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -1).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)

	_, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, nil, nil, analyticsread.GranularityDay)
	require.Error(t, err)
}

func TestGetTimeSeries_CountsRequests(t *testing.T) {
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

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -2).Truncate(24 * time.Hour).Format(time.RFC3339)
	to := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Format(time.RFC3339)

	buckets, err := analyticsRdr.GetTimeSeries(ctx, caller, orgID.String(), from, to, nil, nil, analyticsread.GranularityDay)
	require.NoError(t, err)

	var total int64
	for _, b := range buckets {
		total += b.ReqTotal
	}
	assert.Equal(t, int64(1), total)
}
