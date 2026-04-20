package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// RevokeOrganizationDispatcher translates a gRPC request into a service
// command and removes the employee's organization dispatcher role.
func (h *MembershipHandler) RevokeOrganizationDispatcher(ctx context.Context, req *membershipv1.RevokeOrganizationDispatcherRequest) (*membershipv1.RevokeOrganizationDispatcherResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Organization(orgID)); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationDispatcher(ctx, membership.RevokeOrganizationDispatcherCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationDispatcherResponse{}, nil
}
