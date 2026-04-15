package orgstructure

import (
	"context"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	orgsvc "github.com/medincident/medincident-command-service/internal/services/orgstructure"
)

func (h *OrgStructureHandler) UpdateOrganizationDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateOrganizationDetailsRequest,
) (*orgstructurev1.UpdateOrganizationDetailsResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
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
