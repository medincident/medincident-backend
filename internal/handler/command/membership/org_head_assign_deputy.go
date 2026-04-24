package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignOrganizationHeadDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing OrgHead role.
func (h *MembershipHandler) AssignOrganizationHeadDeputy(ctx context.Context, req *membershipv1.AssignOrganizationHeadDeputyRequest) (*membershipv1.AssignOrganizationHeadDeputyResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationHeadDeputy(ctx, membership.AssignOrganizationHeadDeputyCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignOrganizationHeadDeputyPayload{
			OrganizationID:   req.GetOrganizationId(),
			EmployeeID:       req.GetEmployeeId(),
			DeputyEmployeeID: req.GetDeputyEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationHeadDeputyResponse{}, nil
}
