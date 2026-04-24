package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// RevokeClinicHead translates a gRPC request into a service command
// and removes the employee's clinic head role.
func (h *MembershipHandler) RevokeClinicHead(ctx context.Context, req *membershipv1.RevokeClinicHeadRequest) (*membershipv1.RevokeClinicHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeClinicHead(ctx, membership.RevokeClinicHeadCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RevokeClinicHeadPayload{
			ClinicID:   req.GetClinicId(),
			EmployeeID: req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeClinicHeadResponse{}, nil
}
