package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
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
	if err := h.empSvc.RevokeOrganizationDispatcher(ctx, membership.RevokeOrganizationDispatcherCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationDispatcherResponse{}, nil
}
