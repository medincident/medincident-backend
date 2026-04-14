package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// RevokeOrganizationAdmin translates a gRPC request into a service
// command and removes the employee's organization admin role.
func (h *MembershipHandler) RevokeOrganizationAdmin(ctx context.Context, req *membershipv1.RevokeOrganizationAdminRequest) (*membershipv1.RevokeOrganizationAdminResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationAdmin(ctx, membership.RevokeOrganizationAdminCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationAdminResponse{}, nil
}
