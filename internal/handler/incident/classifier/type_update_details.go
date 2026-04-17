package classifier

import (
	"context"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) UpdateIncidentTypeDetails(
	ctx context.Context,
	req *incidentclassifierv1.UpdateIncidentTypeDetailsRequest,
) (*incidentclassifierv1.UpdateIncidentTypeDetailsResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
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
