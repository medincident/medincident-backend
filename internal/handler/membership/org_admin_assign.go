package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignOrganizationAdmin translates a gRPC request into a service
// command and links the employee to the organization as its admin.
func (h *MembershipHandler) AssignOrganizationAdmin(ctx context.Context, req *membershipv1.AssignOrganizationAdminRequest) (*membershipv1.AssignOrganizationAdminResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
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
