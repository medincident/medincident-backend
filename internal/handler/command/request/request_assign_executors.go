package request

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

func (h *ServiceRequestHandler) AssignExecutors(
	ctx context.Context,
	req *requestv1.AssignExecutorsRequest,
) (*requestv1.AssignExecutorsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.AssignExecutors(ctx, requestsvc.AssignExecutorsCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: requestsvc.AssignExecutorsPayload{
			ServiceRequestID:    req.GetServiceRequestId(),
			ExecutorEmployeeIDs: req.GetExecutorEmployeeIds(),
		},
	}); err != nil {
		return nil, err
	}
	return &requestv1.AssignExecutorsResponse{}, nil
}
