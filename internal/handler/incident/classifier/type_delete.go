package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/services/incident/classifier"
)

func (h *IncidentClassifierHandler) DeleteIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.DeleteIncidentTypeRequest,
) (*incidentclassifierv1.DeleteIncidentTypeResponse, error) {
	typeID, err := parseIncidentTypeID(req.GetTypeId())
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Delete(ctx, classifiersvc.DeleteIncidentTypeCommand{TypeID: typeID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeleteIncidentTypeResponse{}, nil
}
