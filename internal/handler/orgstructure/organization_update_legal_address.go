package orgstructure

import (
	"context"

	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/service/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateOrganizationLegalAddress(
	ctx context.Context,
	req *orgstructurev1.UpdateOrganizationLegalAddressRequest,
) (*orgstructurev1.UpdateOrganizationLegalAddressResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
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
