package request

import (
	"context"
	"strings"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

func (h *ServiceRequestHandler) UpdateServiceRequestDescription(
	ctx context.Context,
	req *requestv1.UpdateServiceRequestDescriptionRequest,
) (*requestv1.UpdateServiceRequestDescriptionResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdateDescription(ctx, requestsvc.UpdateServiceRequestDescriptionCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: requestsvc.UpdateServiceRequestDescriptionPayload{
			ServiceRequestID: req.GetServiceRequestId(),
			Description:      strings.TrimSpace(req.GetDescription()),
		},
	}); err != nil {
		return nil, err
	}
	return &requestv1.UpdateServiceRequestDescriptionResponse{}, nil
}
