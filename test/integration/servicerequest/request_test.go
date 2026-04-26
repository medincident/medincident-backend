//go:build integration

package servicerequest_integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
)

func createRequest(t *testing.T, ctx context.Context, deptID, typeID, empID uuid.UUID) uuid.UUID {
	t.Helper()
	res, err := requestSvc.Create(ctx, &requestsvc.CreateServiceRequestCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.CreateServiceRequestPayload{
			DepartmentID:        deptID.String(),
			TypeID:              typeID.String(),
			Description:         "Описание заявки",
			ExecutorEmployeeIDs: []string{empID.String()},
		},
	})
	require.NoError(t, err)
	return res.ID
}

func TestServiceRequest_Create_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)

	res, err := requestSvc.Create(ctx, &requestsvc.CreateServiceRequestCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.CreateServiceRequestPayload{
			DepartmentID:        deptID.String(),
			TypeID:              typeID.String(),
			Description:         "Тестовое описание заявки",
			ExecutorEmployeeIDs: []string{empID.String()},
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	sr := loadRequest(t, res.ID)
	assert.Equal(t, model.ServiceRequestStatusCreated, sr.Status)
	assert.Equal(t, deptID, sr.DepartmentID)
	assert.Equal(t, typeID, sr.TypeID)

	var projStatus string
	require.NoError(t, testDB.Raw(
		`SELECT status FROM projections.service_requests WHERE id = ?`, res.ID,
	).Scan(&projStatus).Error)
	assert.Equal(t, "created", projStatus)

	var execCount int64
	require.NoError(t, testDB.Raw(
		`SELECT COUNT(*) FROM projections.service_request_executor_history WHERE request_id = ?`, res.ID,
	).Scan(&execCount).Error)
	assert.Equal(t, int64(1), execCount)
}

func TestServiceRequest_Create_DeptNotFound(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	typeID := seedRequestType(t, orgID)

	_, err := requestSvc.Create(ctx, &requestsvc.CreateServiceRequestCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.CreateServiceRequestPayload{
			DepartmentID:        uuid.New().String(),
			TypeID:              typeID.String(),
			Description:         "Тест",
			ExecutorEmployeeIDs: []string{uuid.New().String()},
		},
	})
	require.Error(t, err)
	assert.Equal(t, requestsvc.ErrCodeServiceRequestDeptNotFound, codeOf(t, err))
}

func TestServiceRequest_Create_TypeInactive(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)

	_, err := requestTypeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.DeactivateRequestTypePayload{TypeID: typeID.String()},
	})
	require.NoError(t, err)

	_, err = requestSvc.Create(ctx, &requestsvc.CreateServiceRequestCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.CreateServiceRequestPayload{
			DepartmentID:        deptID.String(),
			TypeID:              typeID.String(),
			Description:         "Тест",
			ExecutorEmployeeIDs: []string{empID.String()},
		},
	})
	require.Error(t, err)
	assert.Equal(t, requestsvc.ErrCodeServiceRequestTypeInactive, codeOf(t, err))
}

func TestServiceRequest_UpdateDescription_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID)

	require.NoError(t, requestSvc.UpdateDescription(ctx, requestsvc.UpdateServiceRequestDescriptionCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.UpdateServiceRequestDescriptionPayload{
			ServiceRequestID: reqID.String(),
			Description:      "Новое описание заявки",
		},
	}))
	assert.Equal(t, "Новое описание заявки", loadRequest(t, reqID).Description)
}

func TestServiceRequest_UpdateDescription_FrozenRequest(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID)

	require.NoError(t, requestSvc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{
			ServiceRequestID: reqID.String(),
			NewStatus:        string(model.ServiceRequestStatusCancelled),
		},
	}))

	err := requestSvc.UpdateDescription(ctx, requestsvc.UpdateServiceRequestDescriptionCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.UpdateServiceRequestDescriptionPayload{
			ServiceRequestID: reqID.String(),
			Description:      "Нельзя изменить",
		},
	})
	require.Error(t, err)
	assert.Equal(t, requestsvc.ErrCodeServiceRequestFrozen, codeOf(t, err))
}

func TestServiceRequest_UpdateStatus_ExecutorTransition(t *testing.T) {
	resetDB(t)
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
			NewStatus:        string(model.ServiceRequestStatusInWork),
		},
	}))
	assert.Equal(t, model.ServiceRequestStatusInWork, loadRequest(t, reqID).Status)

	var histCount int64
	require.NoError(t, testDB.Raw(
		`SELECT COUNT(*) FROM projections.service_request_status_history WHERE request_id = ?`, reqID,
	).Scan(&histCount).Error)
	// Create seeds one NULL→created row; UpdateStatus adds a second created→in_work row.
	assert.Equal(t, int64(2), histCount)
}

func TestServiceRequest_UpdateStatus_ResponsibleCompletes(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID)

	executorCaller := authz.Caller{ZitadelUserID: executorZitadelID}
	require.NoError(t, requestSvc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller:  executorCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{ServiceRequestID: reqID.String(), NewStatus: string(model.ServiceRequestStatusInWork)},
	}))
	require.NoError(t, requestSvc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller:  executorCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{ServiceRequestID: reqID.String(), NewStatus: string(model.ServiceRequestStatusPendingReview)},
	}))
	require.NoError(t, requestSvc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller:  sysadminCaller,
		Payload: requestsvc.UpdateServiceRequestStatusPayload{ServiceRequestID: reqID.String(), NewStatus: string(model.ServiceRequestStatusCompleted)},
	}))
	assert.Equal(t, model.ServiceRequestStatusCompleted, loadRequest(t, reqID).Status)
}

func TestServiceRequest_UpdateStatus_InvalidFlow(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	empID := seedEmployee(t, executorZitadelID, orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID)

	err := requestSvc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller: authz.Caller{ZitadelUserID: executorZitadelID},
		Payload: requestsvc.UpdateServiceRequestStatusPayload{
			ServiceRequestID: reqID.String(),
			NewStatus:        string(model.ServiceRequestStatusCompleted),
		},
	})
	require.Error(t, err)
	assert.Equal(t, requestsvc.ErrCodeServiceRequestInvalidStatusFlow, codeOf(t, err))
}

func TestServiceRequest_AssignExecutors_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := seedOrg(t)
	clinicID := seedClinic(t, orgID)
	deptID := seedDept(t, clinicID)
	typeID := seedRequestType(t, orgID)
	seedUser(t, executorZitadelID, "Executor")
	seedUser(t, "executor2", "Executor2")
	empID1 := seedEmployee(t, executorZitadelID, orgID, deptID)
	empID2 := seedEmployee(t, "executor2", orgID, deptID)
	reqID := createRequest(t, ctx, deptID, typeID, empID1)

	require.NoError(t, requestSvc.AssignExecutors(ctx, requestsvc.AssignExecutorsCommand{
		Caller: sysadminCaller,
		Payload: requestsvc.AssignExecutorsPayload{
			ServiceRequestID:    reqID.String(),
			ExecutorEmployeeIDs: []string{empID2.String()},
		},
	}))

	var count int64
	require.NoError(t, testDB.Raw(
		`SELECT COUNT(*) FROM domain.service_request_executors WHERE request_id = ? AND employee_id = ?`,
		reqID, empID2,
	).Scan(&count).Error)
	assert.Equal(t, int64(1), count)

	require.NoError(t, testDB.Raw(
		`SELECT COUNT(*) FROM domain.service_request_executors WHERE request_id = ? AND employee_id = ?`,
		reqID, empID1,
	).Scan(&count).Error)
	assert.Equal(t, int64(0), count)
}
