package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateOrganizationLegalAddress(
	ctx context.Context,
	req *orgstructurev1.UpdateOrganizationLegalAddressRequest,
) (*orgstructurev1.UpdateOrganizationLegalAddressResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.orgSvc.UpdateLegalAddress(ctx, orgsvc.UpdateOrganizationLegalAddressCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.UpdateOrganizationLegalAddressPayload{
			ID:      req.GetOrganizationId(),
			Address: addressInputFromProto(req.GetLegalAddress()),
		},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateOrganizationLegalAddressResponse{}, nil
}
