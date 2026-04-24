package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignOrganizationHead translates a gRPC request into a service
// command and links the employee to the organization as its head.
func (h *MembershipHandler) AssignOrganizationHead(ctx context.Context, req *membershipv1.AssignOrganizationHeadRequest) (*membershipv1.AssignOrganizationHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationHead(ctx, membership.AssignOrganizationHeadCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignOrganizationHeadPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationHeadResponse{}, nil
}
