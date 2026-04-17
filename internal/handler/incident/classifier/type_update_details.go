package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) UpdateIncidentTypeDetails(
	ctx context.Context,
	req *incidentclassifierv1.UpdateIncidentTypeDetailsRequest,
) (*incidentclassifierv1.UpdateIncidentTypeDetailsResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaType(ctx, middleware.CallerID(ctx), typeID); err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.UpdateDetails(ctx, classifiersvc.UpdateIncidentTypeDetailsCommand{
		TypeID:      typeID,
		Name:        req.GetName(),
		Description: req.Description,
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.UpdateIncidentTypeDetailsResponse{}, nil
}
