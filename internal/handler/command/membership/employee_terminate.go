package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// TerminateEmployee translates a gRPC TerminateEmployeeRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) TerminateEmployee(ctx context.Context, req *membershipv1.TerminateEmployeeRequest) (*membershipv1.TerminateEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.Terminate(ctx, membership.TerminateEmployeeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.TerminateEmployeePayload{
			ID: req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.TerminateEmployeeResponse{}, nil
}
