package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	requestclassifierv1 "github.com/medincident/medincident-backend/pkg/command/request/classifier/v1"
)

func (h *RequestClassifierHandler) ReactivateRequestType(
	ctx context.Context,
	req *requestclassifierv1.ReactivateRequestTypeRequest,
) (*requestclassifierv1.ReactivateRequestTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Reactivate(ctx, classifiersvc.ReactivateRequestTypeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.ReactivateRequestTypePayload{
			TypeID: req.GetTypeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &requestclassifierv1.ReactivateRequestTypeResponse{}, nil
}
