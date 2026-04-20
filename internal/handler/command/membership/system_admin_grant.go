package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

func (h *MembershipHandler) GrantSystemAdmin(ctx context.Context, req *membershipv1.GrantSystemAdminRequest) (*membershipv1.GrantSystemAdminResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.GrantSystemAdmin(ctx, membership.GrantSystemAdminCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.GrantSystemAdminPayload{
			ZitadelUserID: req.GetZitadelUserId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.GrantSystemAdminResponse{}, nil
}
