package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) DeactivateClinic(
	ctx context.Context,
	req *orgstructurev1.DeactivateClinicRequest,
) (*orgstructurev1.DeactivateClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.clinSvc.Deactivate(ctx, orgsvc.DeactivateClinicCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.DeactivateClinicPayload{ID: req.GetClinicId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.DeactivateClinicResponse{}, nil
}

func (h *OrgStructureHandler) ActivateClinic(
	ctx context.Context,
	req *orgstructurev1.ActivateClinicRequest,
) (*orgstructurev1.ActivateClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.clinSvc.Activate(ctx, orgsvc.ActivateClinicCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.ActivateClinicPayload{ID: req.GetClinicId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.ActivateClinicResponse{}, nil
}

func (h *OrgStructureHandler) DeleteClinic(
	ctx context.Context,
	req *orgstructurev1.DeleteClinicRequest,
) (*orgstructurev1.DeleteClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.clinSvc.Delete(ctx, orgsvc.DeleteClinicCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.DeleteClinicPayload{ID: req.GetClinicId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.DeleteClinicResponse{}, nil
}
