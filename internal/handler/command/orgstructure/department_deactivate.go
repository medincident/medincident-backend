package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) DeactivateDepartment(
	ctx context.Context,
	req *orgstructurev1.DeactivateDepartmentRequest,
) (*orgstructurev1.DeactivateDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.deptSvc.Deactivate(ctx, orgsvc.DeactivateDepartmentCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.DeactivateDepartmentPayload{ID: req.GetDepartmentId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.DeactivateDepartmentResponse{}, nil
}

func (h *OrgStructureHandler) ActivateDepartment(
	ctx context.Context,
	req *orgstructurev1.ActivateDepartmentRequest,
) (*orgstructurev1.ActivateDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.deptSvc.Activate(ctx, orgsvc.ActivateDepartmentCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.ActivateDepartmentPayload{ID: req.GetDepartmentId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.ActivateDepartmentResponse{}, nil
}

func (h *OrgStructureHandler) DeleteDepartment(
	ctx context.Context,
	req *orgstructurev1.DeleteDepartmentRequest,
) (*orgstructurev1.DeleteDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.deptSvc.Delete(ctx, orgsvc.DeleteDepartmentCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.DeleteDepartmentPayload{ID: req.GetDepartmentId()},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.DeleteDepartmentResponse{}, nil
}
