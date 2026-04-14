package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignOrganizationDispatcher translates a gRPC request into a service
// command and links the employee to the organization as its dispatcher.
func (h *MembershipHandler) AssignOrganizationDispatcher(ctx context.Context, req *membershipv1.AssignOrganizationDispatcherRequest) (*membershipv1.AssignOrganizationDispatcherResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationDispatcher(ctx, membership.AssignOrganizationDispatcherCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationDispatcherResponse{}, nil
}
