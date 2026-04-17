package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) ReactivateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.ReactivateIncidentCategoryRequest,
) (*incidentclassifierv1.ReactivateIncidentCategoryResponse, error) {
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
	if _, err := h.categorySvc.Reactivate(ctx, classifiersvc.ReactivateIncidentCategoryCommand{CategoryID: categoryID}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.ReactivateIncidentCategoryResponse{}, nil
}
