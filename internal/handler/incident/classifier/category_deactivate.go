package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/services/incident/classifier"
)

func (h *IncidentClassifierHandler) DeactivateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.DeactivateIncidentCategoryRequest,
) (*incidentclassifierv1.DeactivateIncidentCategoryResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{CategoryID: categoryID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeactivateIncidentCategoryResponse{}, nil
}
