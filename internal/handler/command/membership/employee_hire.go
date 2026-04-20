package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// HireEmployee translates a gRPC HireEmployeeRequest into a service
// command and returns the new employee ID on success.
func (h *MembershipHandler) HireEmployee(ctx context.Context, req *membershipv1.HireEmployeeRequest) (*membershipv1.HireEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.empSvc.Hire(ctx, membership.HireEmployeeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.HireEmployeePayload{
			ZitadelUserID: req.GetZitadelUserId(),
			DepartmentID:  req.GetDepartmentId(),
			Position:      req.Position,
		},
	})
	if err != nil {
		return nil, err
	}
	return &membershipv1.HireEmployeeResponse{EmployeeId: result.ID.String()}, nil
}
