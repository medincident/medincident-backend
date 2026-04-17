package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) CreateOrganization(
	ctx context.Context,
	req *orgstructurev1.CreateOrganizationRequest,
) (*orgstructurev1.CreateOrganizationResponse, error) {
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireSystemAdmin(ctx, callerID); err != nil {
		return nil, err
	}
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
