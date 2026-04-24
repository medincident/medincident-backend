package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// UpdateEmployeeDepartment translates a gRPC UpdateEmployeeDepartmentRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeeDepartment(ctx context.Context, req *membershipv1.UpdateEmployeeDepartmentRequest) (*membershipv1.UpdateEmployeeDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.UpdateDepartment(ctx, membership.UpdateEmployeeDepartmentCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.UpdateEmployeeDepartmentPayload{
			ID:           req.GetEmployeeId(),
			DepartmentID: req.GetDepartmentId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateEmployeeDepartmentResponse{}, nil
}
