package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// RemoveOrganizationAdminDeputy translates a gRPC request into a
// service command and clears the deputy slot on an OrgAdmin role.
func (h *MembershipHandler) RemoveOrganizationAdminDeputy(ctx context.Context, req *membershipv1.RemoveOrganizationAdminDeputyRequest) (*membershipv1.RemoveOrganizationAdminDeputyResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveOrganizationAdminDeputy(ctx, membership.RemoveOrganizationAdminDeputyCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveOrganizationAdminDeputyResponse{}, nil
}
