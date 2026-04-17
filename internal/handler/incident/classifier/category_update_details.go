package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) UpdateIncidentCategoryDetails(
	ctx context.Context,
	req *incidentclassifierv1.UpdateIncidentCategoryDetailsRequest,
) (*incidentclassifierv1.UpdateIncidentCategoryDetailsResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.UpdateDetails(ctx, classifiersvc.UpdateIncidentCategoryDetailsCommand{
		CategoryID:  categoryID,
		Name:        req.GetName(),
		Description: req.Description,
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.UpdateIncidentCategoryDetailsResponse{}, nil
}
