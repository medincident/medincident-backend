package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// AssignOrganizationDispatcher translates a gRPC request into a service
// command and links the employee to the organization as its dispatcher.
func (h *MembershipHandler) AssignOrganizationDispatcher(ctx context.Context, req *membershipv1.AssignOrganizationDispatcherRequest) (*membershipv1.AssignOrganizationDispatcherResponse, error) {
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
	if err := h.authz.RequireOrgAdmin(ctx, callerID, orgID); err != nil {
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
