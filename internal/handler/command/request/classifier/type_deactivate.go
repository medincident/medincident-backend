package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	requestclassifierv1 "github.com/medincident/medincident-backend/pkg/command/request/classifier/v1"
)

func (h *RequestClassifierHandler) DeactivateRequestType(
	ctx context.Context,
	req *requestclassifierv1.DeactivateRequestTypeRequest,
) (*requestclassifierv1.DeactivateRequestTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Deactivate(ctx, classifiersvc.DeactivateRequestTypeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.DeactivateRequestTypePayload{
			TypeID: req.GetTypeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &requestclassifierv1.DeactivateRequestTypeResponse{}, nil
}
