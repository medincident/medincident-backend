package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignOrganizationDispatcherDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing OrgDispatcher role.
func (h *MembershipHandler) AssignOrganizationDispatcherDeputy(ctx context.Context, req *membershipv1.AssignOrganizationDispatcherDeputyRequest) (*membershipv1.AssignOrganizationDispatcherDeputyResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	deputyID, err := parseDeputyEmployeeID(req.GetDeputyEmployeeId())
	if err != nil {
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
