package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) ReactivateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.ReactivateIncidentCategoryRequest,
) (*incidentclassifierv1.ReactivateIncidentCategoryResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Reactivate(ctx, classifiersvc.ReactivateIncidentCategoryCommand{CategoryID: categoryID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.ReactivateIncidentCategoryResponse{}, nil
}
