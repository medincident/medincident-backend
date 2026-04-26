//go:build integration

package request_query_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestread "github.com/medincident/medincident-backend/internal/service/query/request"
)

func TestGetServiceRequest_HappyPath(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID)

	got, err := requestReader.GetServiceRequest(ctx, sysadminCaller, reqID)
	require.NoError(t, err)
	assert.Equal(t, reqID, got.ID)
	assert.Equal(t, model.ServiceRequestStatusCreated, got.Status)
	assert.Len(t, got.Executors, 1)
	assert.Equal(t, empID, got.Executors[0].EmployeeID)
}

func TestGetServiceRequest_NotFound(t *testing.T) {
	resetProjections(t)
	_, err := requestReader.GetServiceRequest(context.Background(), sysadminCaller, uuid.New())
	require.Error(t, err)
}

func TestListServiceRequests_HappyPath(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	createRequest(t, ctx, deptID, typeID, empID)
	createRequest(t, ctx, deptID, typeID, empID)

	list, err := requestReader.ListServiceRequests(ctx, sysadminCaller, orgID, requestread.ListQuery{Limit: 10})
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestListServiceRequestsByIncident_HappyPath(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)

	incidentID, err := uuid.NewV7()
	require.NoError(t, err)
	// projections.incidents has no FK constraints — insert directly.
	require.NoError(t, testDB.WithContext(ctx).Exec(`
		INSERT INTO projections.incidents (
			id, organization_id, clinic_id, department_id,
			category_id, type_id, status, priority,
			occurred_at, registrar_employee_id, registrar_display_name,
			registrar_organization_id, registrar_clinic_id, registrar_department_id,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, gen_random_uuid(), gen_random_uuid(),
		          'pending', 'normal', now(), gen_random_uuid(), 'Registrar',
		          ?, ?, ?, now(), now())`,
		incidentID, orgID, clinicID, deptID, orgID, clinicID, deptID,
	).Error)

	reqID, err := uuid.NewV7()
	require.NoError(t, err)
	// projections.service_requests has no FK constraints — insert directly.
	require.NoError(t, testDB.WithContext(ctx).Exec(`
		INSERT INTO projections.service_requests
			(id, organization_id, clinic_id, department_id, type_id,
			 incident_id, description, status, author_id, author_display_name,
			 created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 'Заявка по инциденту', 'created',
		        ?, 'System Admin', now(), now())`,
		reqID, orgID, clinicID, deptID, typeID, incidentID, sysadminZitadelID,
	).Error)

	list, err := requestReader.ListServiceRequestsByIncident(
		ctx, sysadminCaller, incidentID, requestread.ListQuery{Limit: 10},
	)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, reqID, list[0].ID)
}

func TestGetServiceRequestHistory_HappyPath(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID)

	require.NoError(t, requestSvc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller: authz.Caller{ZitadelUserID: executorZitadelID},
		Payload: requestsvc.UpdateServiceRequestStatusPayload{
			ServiceRequestID: reqID.String(),
			NewStatus:        model.ServiceRequestStatusInWork,
		},
	}))

	history, err := requestReader.GetServiceRequestHistory(ctx, sysadminCaller, reqID)
	require.NoError(t, err)
	// Create seeds NULL→created; UpdateStatus adds created→in_work.
	assert.Len(t, history.StatusHistory, 2)
	assert.Equal(t, string(model.ServiceRequestStatusInWork), history.StatusHistory[1].NewStatus)
	assert.Len(t, history.ExecutorHistory, 1)
	assert.Equal(t, empID, history.ExecutorHistory[0].EmployeeID)
}
