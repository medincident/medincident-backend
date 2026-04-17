package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
)

// AssignOrganizationDispatcherDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing OrgDispatcher role.
func (h *MembershipHandler) AssignOrganizationDispatcherDeputy(ctx context.Context, req *membershipv1.AssignOrganizationDispatcherDeputyRequest) (*membershipv1.AssignOrganizationDispatcherDeputyResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	deputyID := ids.parse(req.GetDeputyEmployeeId(), parseDeputyEmployeeID)
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
	if err := h.empSvc.AssignOrganizationDispatcherDeputy(ctx, membership.AssignOrganizationDispatcherDeputyCommand{
		OrganizationID:   orgID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationDispatcherDeputyResponse{}, nil
}
