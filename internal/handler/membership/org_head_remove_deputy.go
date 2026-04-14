package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// RemoveOrganizationHeadDeputy translates a gRPC request into a
// service command and clears the deputy slot on an OrgHead role.
func (h *MembershipHandler) RemoveOrganizationHeadDeputy(ctx context.Context, req *membershipv1.RemoveOrganizationHeadDeputyRequest) (*membershipv1.RemoveOrganizationHeadDeputyResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveOrganizationHeadDeputy(ctx, membership.RemoveOrganizationHeadDeputyCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveOrganizationHeadDeputyResponse{}, nil
}
