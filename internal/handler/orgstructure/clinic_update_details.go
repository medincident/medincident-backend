package orgstructure

import (
	"context"

	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/service/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateClinicDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateClinicDetailsRequest,
) (*orgstructurev1.UpdateClinicDetailsResponse, error) {
	id, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	if err := h.clinSvc.UpdateDetails(ctx, orgsvc.UpdateClinicDetailsCommand{
		ID:          id,
		Name:        req.GetName(),
		Description: req.Description,
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateClinicDetailsResponse{}, nil
}
