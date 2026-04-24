package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// RemoveOrganizationAdminDeputy translates a gRPC request into a
// service command and clears the deputy slot on an OrgAdmin role.
func (h *MembershipHandler) RemoveOrganizationAdminDeputy(ctx context.Context, req *membershipv1.RemoveOrganizationAdminDeputyRequest) (*membershipv1.RemoveOrganizationAdminDeputyResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveOrganizationAdminDeputy(ctx, membership.RemoveOrganizationAdminDeputyCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RemoveOrganizationAdminDeputyPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveOrganizationAdminDeputyResponse{}, nil
}
