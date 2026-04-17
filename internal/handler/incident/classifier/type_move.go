package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) MoveIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.MoveIncidentTypeRequest,
) (*incidentclassifierv1.MoveIncidentTypeResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaType(ctx, middleware.CallerID(ctx), typeID); err != nil {
		return nil, err
	}
	newCategoryID, err := parseIncidentCategoryID(req.GetNewCategoryId())
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Move(ctx, classifiersvc.MoveIncidentTypeCommand{
		TypeID:        typeID,
		NewCategoryID: newCategoryID,
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.MoveIncidentTypeResponse{}, nil
}
