package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// RevokeOrganizationDispatcher translates a gRPC request into a service
// command and removes the employee's organization dispatcher role.
func (h *MembershipHandler) RevokeOrganizationDispatcher(ctx context.Context, req *membershipv1.RevokeOrganizationDispatcherRequest) (*membershipv1.RevokeOrganizationDispatcherResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationDispatcher(ctx, membership.RevokeOrganizationDispatcherCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RevokeOrganizationDispatcherPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationDispatcherResponse{}, nil
}
