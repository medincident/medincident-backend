package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) ReactivateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.ReactivateIncidentTypeRequest,
) (*incidentclassifierv1.ReactivateIncidentTypeResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Reactivate(ctx, classifiersvc.ReactivateIncidentTypeCommand{TypeID: typeID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.ReactivateIncidentTypeResponse{}, nil
}
