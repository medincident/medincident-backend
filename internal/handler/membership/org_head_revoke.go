package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// RevokeOrganizationHead translates a gRPC request into a service
// command and removes the employee's organization head role.
func (h *MembershipHandler) RevokeOrganizationHead(ctx context.Context, req *membershipv1.RevokeOrganizationHeadRequest) (*membershipv1.RevokeOrganizationHeadResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Organization(orgID)); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationHead(ctx, membership.RevokeOrganizationHeadCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationHeadResponse{}, nil
}
