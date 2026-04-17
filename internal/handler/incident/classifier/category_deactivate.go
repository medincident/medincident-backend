package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
)

func (h *IncidentClassifierHandler) DeactivateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.DeactivateIncidentCategoryRequest,
) (*incidentclassifierv1.DeactivateIncidentCategoryResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaCategory(ctx, middleware.CallerID(ctx), categoryID); err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{CategoryID: categoryID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeactivateIncidentCategoryResponse{}, nil
}
