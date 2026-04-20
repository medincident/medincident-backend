package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateOrganizationLegalAddress(
	ctx context.Context,
	req *orgstructurev1.UpdateOrganizationLegalAddressRequest,
) (*orgstructurev1.UpdateOrganizationLegalAddressResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Organization(id)); err != nil {
		return nil, err
	}
	if err := h.orgSvc.UpdateLegalAddress(ctx, orgsvc.UpdateOrganizationLegalAddressCommand{
		ID:      id,
		Address: addressInputFromProto(req.GetLegalAddress()),
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateOrganizationLegalAddressResponse{}, nil
}
