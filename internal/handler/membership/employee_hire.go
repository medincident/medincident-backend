package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// HireEmployee translates a gRPC HireEmployeeRequest into a service
// command and returns the new employee ID on success.
func (h *MembershipHandler) HireEmployee(ctx context.Context, req *membershipv1.HireEmployeeRequest) (*membershipv1.HireEmployeeResponse, error) {
	depID, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Department(depID)); err != nil {
		return nil, err
	}
	cmd := membership.HireEmployeeCommand{
		ZitadelUserID: req.GetZitadelUserId(),
		DepartmentID:  depID,
		Position:      req.Position,
	}
	result, err := h.empSvc.Hire(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &membershipv1.HireEmployeeResponse{EmployeeId: result.ID.String()}, nil
}
