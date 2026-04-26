package classifier

import (
	"context"
	"strings"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	requestclassifierv1 "github.com/medincident/medincident-backend/pkg/command/request/classifier/v1"
)

func (h *RequestClassifierHandler) CreateRequestType(
	ctx context.Context,
	req *requestclassifierv1.CreateRequestTypeRequest,
) (*requestclassifierv1.CreateRequestTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.typeSvc.Create(ctx, classifiersvc.CreateRequestTypeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.CreateRequestTypePayload{
			OrganizationID: req.GetOrganizationId(),
			Name:           strings.TrimSpace(req.GetName()),
			Description:    trimOptionalString(req.Description),
		},
	})
	if err != nil {
		return nil, err
	}
	return &requestclassifierv1.CreateRequestTypeResponse{TypeId: result.ID.String()}, nil
}
