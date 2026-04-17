package orgstructure

import (
	"context"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
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
