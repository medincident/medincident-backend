package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// ScheduleVacation translates a gRPC ScheduleVacationRequest into a
// service command and returns the new vacation ID on success.
func (h *MembershipHandler) ScheduleVacation(ctx context.Context, req *membershipv1.ScheduleVacationRequest) (*membershipv1.ScheduleVacationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	payload := membership.ScheduleVacationPayload{
		EmployeeID: req.GetEmployeeId(),
	}
	if ts := req.GetStartsAt(); ts != nil {
		payload.StartsAt = ts.AsTime()
	}
	if ts := req.GetEndsAt(); ts != nil {
		t := ts.AsTime()
		payload.EndsAt = &t
	}
	res, err := h.empSvc.ScheduleVacation(ctx, membership.ScheduleVacationCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}
	return &membershipv1.ScheduleVacationResponse{VacationId: res.ID.String()}, nil
}
