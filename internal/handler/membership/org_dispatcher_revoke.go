package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// RevokeOrganizationDispatcher translates a gRPC request into a service
// command and removes the employee's organization dispatcher role.
func (h *MembershipHandler) RevokeOrganizationDispatcher(ctx context.Context, req *membershipv1.RevokeOrganizationDispatcherRequest) (*membershipv1.RevokeOrganizationDispatcherResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
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
