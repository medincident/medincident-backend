package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

func (h *MembershipHandler) GrantSystemAdmin(ctx context.Context, req *membershipv1.GrantSystemAdminRequest) (*membershipv1.GrantSystemAdminResponse, error) {
	if err := h.empSvc.GrantSystemAdmin(ctx, membership.GrantSystemAdminCommand{
		ZitadelUserID: req.GetZitadelUserId(),
	}); err != nil {
		return nil, err
	}
	return &membershipv1.GrantSystemAdminResponse{}, nil
}
