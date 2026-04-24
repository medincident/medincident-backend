package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

func (h *MembershipHandler) RevokeSystemAdmin(ctx context.Context, req *membershipv1.RevokeSystemAdminRequest) (*membershipv1.RevokeSystemAdminResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeSystemAdmin(ctx, membership.RevokeSystemAdminCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RevokeSystemAdminPayload{
			ZitadelUserID: req.GetZitadelUserId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeSystemAdminResponse{}, nil
}
