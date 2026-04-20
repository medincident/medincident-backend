package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) CreateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentTypeRequest,
) (*incidentclassifierv1.CreateIncidentTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID:  req.GetCategoryId(),
			Name:        req.GetName(),
			Description: req.Description,
		},
	})
	if err != nil {
		return nil, err
	}
	return &incidentclassifierv1.CreateIncidentTypeResponse{TypeId: result.ID.String()}, nil
}
