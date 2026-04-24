package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignOrganizationDispatcherDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing OrgDispatcher role.
func (h *MembershipHandler) AssignOrganizationDispatcherDeputy(ctx context.Context, req *membershipv1.AssignOrganizationDispatcherDeputyRequest) (*membershipv1.AssignOrganizationDispatcherDeputyResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignOrganizationDispatcherDeputy(ctx, membership.AssignOrganizationDispatcherDeputyCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignOrganizationDispatcherDeputyPayload{
			OrganizationID:   req.GetOrganizationId(),
			EmployeeID:       req.GetEmployeeId(),
			DeputyEmployeeID: req.GetDeputyEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignOrganizationDispatcherDeputyResponse{}, nil
}
