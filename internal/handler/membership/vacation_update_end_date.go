package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// UpdateVacationEndDate translates a gRPC UpdateVacationEndDateRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateVacationEndDate(ctx context.Context, req *membershipv1.UpdateVacationEndDateRequest) (*membershipv1.UpdateVacationEndDateResponse, error) {
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
	cmd := membership.UpdateVacationEndDateCommand{VacationID: id}
	if ts := req.GetEndsAt(); ts != nil {
		cmd.EndsAt = ts.AsTime()
	}
	if err := h.empSvc.UpdateVacationEndDate(ctx, cmd); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateVacationEndDateResponse{}, nil
}
