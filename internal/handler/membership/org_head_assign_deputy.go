package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// AssignOrganizationHeadDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing OrgHead role.
func (h *MembershipHandler) AssignOrganizationHeadDeputy(ctx context.Context, req *membershipv1.AssignOrganizationHeadDeputyRequest) (*membershipv1.AssignOrganizationHeadDeputyResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	deputyID := ids.parse(req.GetDeputyEmployeeId(), parseDeputyEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationHeadDeputy(ctx, membership.AssignOrganizationHeadDeputyCommand{
		OrganizationID:   orgID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationHeadDeputyResponse{}, nil
}
