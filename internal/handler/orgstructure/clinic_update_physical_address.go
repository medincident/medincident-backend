package orgstructure

import (
	"context"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
)

func (h *OrgStructureHandler) UpdateClinicPhysicalAddress(
	ctx context.Context,
	req *orgstructurev1.UpdateClinicPhysicalAddressRequest,
) (*orgstructurev1.UpdateClinicPhysicalAddressResponse, error) {
	id, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaClinic(ctx, callerID, id); err != nil {
		return nil, err
	}
	if err := h.clinSvc.UpdatePhysicalAddress(ctx, orgsvc.UpdateClinicPhysicalAddressCommand{
		ID:      id,
		Address: addressInputFromProto(req.GetPhysicalAddress()),
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateClinicPhysicalAddressResponse{}, nil
}
