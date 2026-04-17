package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// RemoveOrganizationDispatcherDeputy translates a gRPC request into a
// service command and clears the deputy slot on an OrgDispatcher role.
func (h *MembershipHandler) RemoveOrganizationDispatcherDeputy(ctx context.Context, req *membershipv1.RemoveOrganizationDispatcherDeputyRequest) (*membershipv1.RemoveOrganizationDispatcherDeputyResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdmin(ctx, middleware.CallerID(ctx), orgID); err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveOrganizationDispatcherDeputy(ctx, membership.RemoveOrganizationDispatcherDeputyCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveOrganizationDispatcherDeputyResponse{}, nil
}
