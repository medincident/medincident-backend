package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

func (h *MembershipHandler) GrantSystemAdmin(ctx context.Context, req *membershipv1.GrantSystemAdminRequest) (*membershipv1.GrantSystemAdminResponse, error) {
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireSystemAdmin(ctx, callerID); err != nil {
		return nil, err
	}
	if err := h.empSvc.GrantSystemAdmin(ctx, membership.GrantSystemAdminCommand{
		ZitadelUserID: req.GetZitadelUserId(),
	}); err != nil {
		return nil, err
	}
	return &membershipv1.GrantSystemAdminResponse{}, nil
}
