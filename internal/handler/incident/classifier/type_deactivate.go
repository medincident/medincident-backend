package classifier

import (
	"context"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) DeactivateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.DeactivateIncidentTypeRequest,
) (*incidentclassifierv1.DeactivateIncidentTypeResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Deactivate(ctx, classifiersvc.DeactivateIncidentTypeCommand{TypeID: typeID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeactivateIncidentTypeResponse{}, nil
}
