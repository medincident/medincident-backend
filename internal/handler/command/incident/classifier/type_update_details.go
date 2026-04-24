package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) UpdateIncidentTypeDetails(
	ctx context.Context,
	req *incidentclassifierv1.UpdateIncidentTypeDetailsRequest,
) (*incidentclassifierv1.UpdateIncidentTypeDetailsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.UpdateDetails(ctx, classifiersvc.UpdateIncidentTypeDetailsCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.UpdateIncidentTypeDetailsPayload{
			TypeID:      req.GetTypeId(),
			Name:        req.GetName(),
			Description: req.Description,
		},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.UpdateIncidentTypeDetailsResponse{}, nil
}
