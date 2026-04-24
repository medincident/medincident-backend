package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignOrganizationDispatcher translates a gRPC request into a service
// command and links the employee to the organization as its dispatcher.
func (h *MembershipHandler) AssignOrganizationDispatcher(ctx context.Context, req *membershipv1.AssignOrganizationDispatcherRequest) (*membershipv1.AssignOrganizationDispatcherResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationDispatcher(ctx, membership.AssignOrganizationDispatcherCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignOrganizationDispatcherPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationDispatcherResponse{}, nil
}
