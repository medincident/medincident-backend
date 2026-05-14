package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) DeactivateOrganization(
	ctx context.Context,
	req *orgstructurev1.DeactivateOrganizationRequest,
) (*orgstructurev1.DeactivateOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.orgSvc.Deactivate(ctx, orgsvc.DeactivateOrganizationCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.DeactivateOrganizationPayload{ID: req.GetOrganizationId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.DeactivateOrganizationResponse{}, nil
}

func (h *OrgStructureHandler) ActivateOrganization(
	ctx context.Context,
	req *orgstructurev1.ActivateOrganizationRequest,
) (*orgstructurev1.ActivateOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.orgSvc.Activate(ctx, orgsvc.ActivateOrganizationCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.ActivateOrganizationPayload{ID: req.GetOrganizationId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.ActivateOrganizationResponse{}, nil
}

func (h *OrgStructureHandler) DeleteOrganization(
	ctx context.Context,
	req *orgstructurev1.DeleteOrganizationRequest,
) (*orgstructurev1.DeleteOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.orgSvc.Delete(ctx, orgsvc.DeleteOrganizationCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.DeleteOrganizationPayload{ID: req.GetOrganizationId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.DeleteOrganizationResponse{}, nil
}
