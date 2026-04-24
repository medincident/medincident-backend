package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) DeactivateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.DeactivateIncidentTypeRequest,
) (*incidentclassifierv1.DeactivateIncidentTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Deactivate(ctx, classifiersvc.DeactivateIncidentTypeCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.DeactivateIncidentTypePayload{TypeID: req.GetTypeId()},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeactivateIncidentTypeResponse{}, nil
}
