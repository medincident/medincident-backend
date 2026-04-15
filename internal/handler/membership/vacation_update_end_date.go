package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// UpdateVacationEndDate translates a gRPC UpdateVacationEndDateRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateVacationEndDate(ctx context.Context, req *membershipv1.UpdateVacationEndDateRequest) (*membershipv1.UpdateVacationEndDateResponse, error) {
	id, err := parseVacationID(req.GetVacationId())
	if err != nil {
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
