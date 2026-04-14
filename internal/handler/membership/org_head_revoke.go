package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// RevokeOrganizationHead translates a gRPC request into a service
// command and removes the employee's organization head role.
func (h *MembershipHandler) RevokeOrganizationHead(ctx context.Context, req *membershipv1.RevokeOrganizationHeadRequest) (*membershipv1.RevokeOrganizationHeadResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeOrganizationHead(ctx, membership.RevokeOrganizationHeadCommand{
		OrganizationID: orgID,
		EmployeeID:     empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeOrganizationHeadResponse{}, nil
}
