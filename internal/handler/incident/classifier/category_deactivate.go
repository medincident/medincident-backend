package classifier

import (
	"context"

	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
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
