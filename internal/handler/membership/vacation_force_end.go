package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// ForceEndVacation translates a gRPC ForceEndVacationRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) ForceEndVacation(ctx context.Context, req *membershipv1.ForceEndVacationRequest) (*membershipv1.ForceEndVacationResponse, error) {
	id, err := parseVacationID(req.GetVacationId())
	if err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Vacation(id)); err != nil {
		return nil, err
	}
	if err := h.empSvc.ForceEndVacation(ctx, membership.ForceEndVacationCommand{VacationID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.ForceEndVacationResponse{}, nil
}
