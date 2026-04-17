package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
)

func (h *MembershipHandler) GrantSystemAdmin(ctx context.Context, req *membershipv1.GrantSystemAdminRequest) (*membershipv1.GrantSystemAdminResponse, error) {
	if err := h.empSvc.GrantSystemAdmin(ctx, membership.GrantSystemAdminCommand{
		ZitadelUserID: req.GetZitadelUserId(),
	}); err != nil {
		return nil, err
	}
	return &membershipv1.GrantSystemAdminResponse{}, nil
}
