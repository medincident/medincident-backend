package classifier

import (
	"context"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) CreateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentTypeRequest,
) (*incidentclassifierv1.CreateIncidentTypeResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	result, err := h.typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID:  categoryID,
		Name:        req.GetName(),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &incidentclassifierv1.CreateIncidentTypeResponse{TypeId: result.ID.String()}, nil
}
