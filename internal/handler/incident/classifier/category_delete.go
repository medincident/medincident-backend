package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) DeleteIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.DeleteIncidentCategoryRequest,
) (*incidentclassifierv1.DeleteIncidentCategoryResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Delete(ctx, classifiersvc.DeleteIncidentCategoryCommand{CategoryID: categoryID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeleteIncidentCategoryResponse{}, nil
}
