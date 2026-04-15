package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/services/incident/classifier"
)

func (h *IncidentClassifierHandler) CreateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentCategoryRequest,
) (*incidentclassifierv1.CreateIncidentCategoryResponse, error) {
	organizationID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
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
