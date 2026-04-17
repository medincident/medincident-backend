package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// AssignOrganizationAdminDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing OrgAdmin role.
func (h *MembershipHandler) AssignOrganizationAdminDeputy(ctx context.Context, req *membershipv1.AssignOrganizationAdminDeputyRequest) (*membershipv1.AssignOrganizationAdminDeputyResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	deputyID := ids.parse(req.GetDeputyEmployeeId(), parseDeputyEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdmin(ctx, middleware.CallerID(ctx), orgID); err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationAdminDeputy(ctx, membership.AssignOrganizationAdminDeputyCommand{
		OrganizationID:   orgID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationAdminDeputyResponse{}, nil
}
