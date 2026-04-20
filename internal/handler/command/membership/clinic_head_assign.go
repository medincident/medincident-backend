package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignClinicHead translates a gRPC request into a service command
// and links the employee to the clinic as its head.
func (h *MembershipHandler) AssignClinicHead(ctx context.Context, req *membershipv1.AssignClinicHeadRequest) (*membershipv1.AssignClinicHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignClinicHead(ctx, membership.AssignClinicHeadCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignClinicHeadPayload{
			ClinicID:   req.GetClinicId(),
			EmployeeID: req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignClinicHeadResponse{}, nil
}
