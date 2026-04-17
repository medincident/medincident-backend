package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// AssignDepartmentResponsibleDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing DR role.
func (h *MembershipHandler) AssignDepartmentResponsibleDeputy(ctx context.Context, req *membershipv1.AssignDepartmentResponsibleDeputyRequest) (*membershipv1.AssignDepartmentResponsibleDeputyResponse, error) {
	var ids idErrs
	depID := ids.parse(req.GetDepartmentId(), parseDepartmentID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	deputyID := ids.parse(req.GetDeputyEmployeeId(), parseDeputyEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaDepartment(ctx, callerID, depID); err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignDepartmentResponsibleDeputy(ctx, membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     depID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignDepartmentResponsibleDeputyResponse{}, nil
}
