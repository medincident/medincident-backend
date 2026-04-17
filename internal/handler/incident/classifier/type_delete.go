package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) DeleteIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.DeleteIncidentTypeRequest,
) (*incidentclassifierv1.DeleteIncidentTypeResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaType(ctx, callerID, typeID); err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Delete(ctx, classifiersvc.DeleteIncidentTypeCommand{TypeID: typeID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeleteIncidentTypeResponse{}, nil
}
