package request

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

func (h *ServiceRequestHandler) UpdateServiceRequestStatus(
	ctx context.Context,
	req *requestv1.UpdateServiceRequestStatusRequest,
) (*requestv1.UpdateServiceRequestStatusResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdateStatus(ctx, requestsvc.UpdateServiceRequestStatusCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: requestsvc.UpdateServiceRequestStatusPayload{
			ServiceRequestID: req.GetServiceRequestId(),
			NewStatus:        req.GetNewStatus(),
		},
	}); err != nil {
		return nil, err
	}
	return &requestv1.UpdateServiceRequestStatusResponse{}, nil
}
