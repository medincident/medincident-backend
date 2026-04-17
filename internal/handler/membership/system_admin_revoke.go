package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

func (h *MembershipHandler) RevokeSystemAdmin(ctx context.Context, req *membershipv1.RevokeSystemAdminRequest) (*membershipv1.RevokeSystemAdminResponse, error) {
	if err := h.authz.RequireSystemAdmin(ctx, middleware.CallerID(ctx)); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeSystemAdmin(ctx, membership.RevokeSystemAdminCommand{
		ZitadelUserID: req.GetZitadelUserId(),
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeSystemAdminResponse{}, nil
}
