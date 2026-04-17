package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignDepartmentResponsible translates a gRPC request into a service
// command and links the employee to the department as its responsible.
func (h *MembershipHandler) AssignDepartmentResponsible(ctx context.Context, req *membershipv1.AssignDepartmentResponsibleRequest) (*membershipv1.AssignDepartmentResponsibleResponse, error) {
	var ids idErrs
	depID := ids.parse(req.GetDepartmentId(), parseDepartmentID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
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
	if err := h.empSvc.AssignDepartmentResponsible(ctx, membership.AssignDepartmentResponsibleCommand{
		DepartmentID: depID,
		EmployeeID:   empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignDepartmentResponsibleResponse{}, nil
}
