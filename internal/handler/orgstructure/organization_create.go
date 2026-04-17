package orgstructure

import (
	"context"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
)

func (h *OrgStructureHandler) CreateOrganization(
	ctx context.Context,
	req *orgstructurev1.CreateOrganizationRequest,
) (*orgstructurev1.CreateOrganizationResponse, error) {
	result, err := h.orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Name:         req.GetName(),
		Description:  req.Description,
		LegalAddress: addressInputFromProto(req.GetLegalAddress()),
	})
	if err != nil {
		return nil, err
	}
	return &orgstructurev1.CreateOrganizationResponse{OrganizationId: result.ID.String()}, nil
}
