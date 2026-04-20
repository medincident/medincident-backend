package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// AssignOrganizationHead translates a gRPC request into a service
// command and links the employee to the organization as its head.
func (h *MembershipHandler) AssignOrganizationHead(ctx context.Context, req *membershipv1.AssignOrganizationHeadRequest) (*membershipv1.AssignOrganizationHeadResponse, error) {
	var ids idErrs
	orgID := ids.parse(req.GetOrganizationId(), parseOrganizationID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Organization(orgID)); err != nil {
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
