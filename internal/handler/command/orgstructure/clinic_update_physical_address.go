package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateClinicPhysicalAddress(
	ctx context.Context,
	req *orgstructurev1.UpdateClinicPhysicalAddressRequest,
) (*orgstructurev1.UpdateClinicPhysicalAddressResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.clinSvc.UpdatePhysicalAddress(ctx, orgsvc.UpdateClinicPhysicalAddressCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.UpdateClinicPhysicalAddressPayload{
			ID:      req.GetClinicId(),
			Address: addressInputFromProto(req.GetPhysicalAddress()),
		},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateClinicPhysicalAddressResponse{}, nil
}
