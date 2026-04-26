package classifier

import (
	"context"
	"strings"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	requestclassifierv1 "github.com/medincident/medincident-backend/pkg/command/request/classifier/v1"
)

func (h *RequestClassifierHandler) UpdateRequestTypeDetails(
	ctx context.Context,
	req *requestclassifierv1.UpdateRequestTypeDetailsRequest,
) (*requestclassifierv1.UpdateRequestTypeDetailsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.UpdateDetails(ctx, classifiersvc.UpdateRequestTypeDetailsCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.UpdateRequestTypeDetailsPayload{
			TypeID:      req.GetTypeId(),
			Name:        strings.TrimSpace(req.GetName()),
			Description: trimOptionalString(req.Description),
		},
	}); err != nil {
		return nil, err
	}
	return &requestclassifierv1.UpdateRequestTypeDetailsResponse{}, nil
}
