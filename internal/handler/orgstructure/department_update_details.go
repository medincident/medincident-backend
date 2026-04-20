package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateDepartmentDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateDepartmentDetailsRequest,
) (*orgstructurev1.UpdateDepartmentDetailsResponse, error) {
	id, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Department(id)); err != nil {
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
