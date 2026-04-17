package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignOrganizationAdmin translates a gRPC request into a service
// command and links the employee to the organization as its admin.
func (h *MembershipHandler) AssignOrganizationAdmin(ctx context.Context, req *membershipv1.AssignOrganizationAdminRequest) (*membershipv1.AssignOrganizationAdminResponse, error) {
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
	if err := h.empSvc.AssignOrganizationAdmin(ctx, membership.AssignOrganizationAdminCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationAdminResponse{}, nil
}
