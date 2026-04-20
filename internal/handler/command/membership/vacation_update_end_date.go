package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// UpdateVacationEndDate translates a gRPC UpdateVacationEndDateRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateVacationEndDate(ctx context.Context, req *membershipv1.UpdateVacationEndDateRequest) (*membershipv1.UpdateVacationEndDateResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	payload := membership.UpdateVacationEndDatePayload{
		VacationID: req.GetVacationId(),
	}
	if ts := req.GetEndsAt(); ts != nil {
		payload.EndsAt = ts.AsTime()
	}
	if err := h.empSvc.UpdateVacationEndDate(ctx, membership.UpdateVacationEndDateCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: payload,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateVacationEndDateResponse{}, nil
}
