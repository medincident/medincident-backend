package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) ReactivateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.ReactivateIncidentTypeRequest,
) (*incidentclassifierv1.ReactivateIncidentTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Reactivate(ctx, classifiersvc.ReactivateIncidentTypeCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.ReactivateIncidentTypePayload{TypeID: req.GetTypeId()},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.ReactivateIncidentTypeResponse{}, nil
}
