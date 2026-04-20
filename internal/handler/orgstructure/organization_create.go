package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) CreateOrganization(
	ctx context.Context,
	req *orgstructurev1.CreateOrganizationRequest,
) (*orgstructurev1.CreateOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.SystemAdmin); err != nil {
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
