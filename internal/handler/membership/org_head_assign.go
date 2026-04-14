package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignOrganizationHead translates a gRPC request into a service
// command and links the employee to the organization as its head.
func (h *MembershipHandler) AssignOrganizationHead(ctx context.Context, req *membershipv1.AssignOrganizationHeadRequest) (*membershipv1.AssignOrganizationHeadResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationHead(ctx, membership.AssignOrganizationHeadCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationHeadResponse{}, nil
}
