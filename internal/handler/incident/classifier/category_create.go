package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) CreateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentCategoryRequest,
) (*incidentclassifierv1.CreateIncidentCategoryResponse, error) {
	organizationID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Organization(organizationID)); err != nil {
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
