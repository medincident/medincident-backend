package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
)

// UpdateEmployeeDepartment translates a gRPC UpdateEmployeeDepartmentRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeeDepartment(ctx context.Context, req *membershipv1.UpdateEmployeeDepartmentRequest) (*membershipv1.UpdateEmployeeDepartmentResponse, error) {
	var ids idErrs
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	depID := ids.parse(req.GetDepartmentId(), parseDepartmentID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaEmployee(ctx, callerID, empID); err != nil {
		return nil, err
	}
	if err := h.empSvc.UpdateDepartment(ctx, membership.UpdateEmployeeDepartmentCommand{
		ID:           empID,
		DepartmentID: depID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateEmployeeDepartmentResponse{}, nil
}
