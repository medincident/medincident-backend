package classifier

import (
	"context"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) ReactivateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.ReactivateIncidentTypeRequest,
) (*incidentclassifierv1.ReactivateIncidentTypeResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Reactivate(ctx, classifiersvc.ReactivateIncidentTypeCommand{TypeID: typeID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.ReactivateIncidentTypeResponse{}, nil
}
