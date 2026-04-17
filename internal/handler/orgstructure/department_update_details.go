package orgstructure

import (
	"context"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
)

func (h *OrgStructureHandler) UpdateDepartmentDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateDepartmentDetailsRequest,
) (*orgstructurev1.UpdateDepartmentDetailsResponse, error) {
	id, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaDepartment(ctx, callerID, id); err != nil {
		return nil, err
	}
	if err := h.deptSvc.UpdateDetails(ctx, orgsvc.UpdateDepartmentDetailsCommand{
		ID:          id,
		Name:        req.GetName(),
		Description: req.Description,
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateDepartmentDetailsResponse{}, nil
}
