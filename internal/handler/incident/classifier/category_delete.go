package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) DeleteIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.DeleteIncidentCategoryRequest,
) (*incidentclassifierv1.DeleteIncidentCategoryResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaCategory(ctx, callerID, categoryID); err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Delete(ctx, classifiersvc.DeleteIncidentCategoryCommand{CategoryID: categoryID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeleteIncidentCategoryResponse{}, nil
}
