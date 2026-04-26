package request

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

func (h *ServiceRequestHandler) CreateServiceRequest(
	ctx context.Context,
	req *requestv1.CreateServiceRequestRequest,
) (*requestv1.CreateServiceRequestResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.svc.Create(ctx, &requestsvc.CreateServiceRequestCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: requestsvc.CreateServiceRequestPayload{
			DepartmentID:        req.GetDepartmentId(),
			TypeID:              req.GetTypeId(),
			IncidentID:          req.IncidentId,
			Description:         req.GetDescription(),
			ExecutorEmployeeIDs: req.GetExecutorEmployeeIds(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &requestv1.CreateServiceRequestResponse{ServiceRequestId: result.ID.String()}, nil
}
