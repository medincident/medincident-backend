package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// RemoveOrganizationHeadDeputy translates a gRPC request into a
// service command and clears the deputy slot on an OrgHead role.
func (h *MembershipHandler) RemoveOrganizationHeadDeputy(ctx context.Context, req *membershipv1.RemoveOrganizationHeadDeputyRequest) (*membershipv1.RemoveOrganizationHeadDeputyResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveOrganizationHeadDeputy(ctx, membership.RemoveOrganizationHeadDeputyCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RemoveOrganizationHeadDeputyPayload{
			OrganizationID: req.GetOrganizationId(),
			EmployeeID:     req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveOrganizationHeadDeputyResponse{}, nil
}
