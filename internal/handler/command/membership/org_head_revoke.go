package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// RevokeOrganizationHead translates a gRPC request into a service
// command and removes the employee's organization head role.
func (h *MembershipHandler) RevokeOrganizationHead(ctx context.Context, req *membershipv1.RevokeOrganizationHeadRequest) (*membershipv1.RevokeOrganizationHeadResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationHead(ctx, membership.RevokeOrganizationHeadCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RevokeOrganizationHeadPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationHeadResponse{}, nil
}
