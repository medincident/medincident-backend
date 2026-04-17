package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
)

func (h *MembershipHandler) RevokeSystemAdmin(ctx context.Context, req *membershipv1.RevokeSystemAdminRequest) (*membershipv1.RevokeSystemAdminResponse, error) {
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireSystemAdmin(ctx, callerID); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeSystemAdmin(ctx, membership.RevokeSystemAdminCommand{
		ZitadelUserID: req.GetZitadelUserId(),
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeSystemAdminResponse{}, nil
}
