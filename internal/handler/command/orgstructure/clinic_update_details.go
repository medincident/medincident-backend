package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateClinicDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateClinicDetailsRequest,
) (*orgstructurev1.UpdateClinicDetailsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.clinSvc.UpdateDetails(ctx, orgsvc.UpdateClinicDetailsCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.UpdateClinicDetailsPayload{
			ID:          req.GetClinicId(),
			Name:        req.GetName(),
			Description: req.Description,
		},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateClinicDetailsResponse{}, nil
}
