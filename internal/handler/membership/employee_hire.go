package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
)

// HireEmployee translates a gRPC HireEmployeeRequest into a service
// command and returns the new employee ID on success.
func (h *MembershipHandler) HireEmployee(ctx context.Context, req *membershipv1.HireEmployeeRequest) (*membershipv1.HireEmployeeResponse, error) {
	depID, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
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
