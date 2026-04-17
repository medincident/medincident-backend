package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) CreateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentCategoryRequest,
) (*incidentclassifierv1.CreateIncidentCategoryResponse, error) {
	organizationID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdmin(ctx, callerID, organizationID); err != nil {
		return nil, err
	}
	cmd := classifiersvc.CreateIncidentCategoryCommand{
		OrganizationID: organizationID,
		Name:           req.GetName(),
		Description:    req.Description,
	}
	if req.ParentCategoryId != nil {
		parentID, err := parseIncidentCategoryID(*req.ParentCategoryId)
		if err != nil {
			return nil, err
		}
		cmd.ParentCategoryID = &parentID
	}
	result, err := h.categorySvc.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &incidentclassifierv1.CreateIncidentCategoryResponse{
		CategoryId: result.ID.String(),
	}, nil
}
