package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

func (h *MembershipHandler) RevokeSystemAdmin(ctx context.Context, req *membershipv1.RevokeSystemAdminRequest) (*membershipv1.RevokeSystemAdminResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.SystemAdmin); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeSystemAdmin(ctx, membership.RevokeSystemAdminCommand{
		ZitadelUserID: req.GetZitadelUserId(),
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeSystemAdminResponse{}, nil
}
