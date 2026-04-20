package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// RevokeOrganizationAdmin translates a gRPC request into a service
// command and removes the employee's organization admin role.
func (h *MembershipHandler) RevokeOrganizationAdmin(ctx context.Context, req *membershipv1.RevokeOrganizationAdminRequest) (*membershipv1.RevokeOrganizationAdminResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationAdmin(ctx, membership.RevokeOrganizationAdminCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RevokeOrganizationAdminPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationAdminResponse{}, nil
}
