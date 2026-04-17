package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateOrganizationDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateOrganizationDetailsRequest,
) (*orgstructurev1.UpdateOrganizationDetailsResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdmin(ctx, callerID, id); err != nil {
		return nil, err
	}
	if err := h.orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		ID:          id,
		Name:        req.GetName(),
		Description: req.Description,
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateOrganizationDetailsResponse{}, nil
}
