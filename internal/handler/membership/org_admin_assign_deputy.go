package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
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
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Organization(orgID)); err != nil {
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
