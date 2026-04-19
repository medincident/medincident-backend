package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// ScheduleVacation translates a gRPC ScheduleVacationRequest into a
// service command and returns the new vacation ID on success.
func (h *MembershipHandler) ScheduleVacation(ctx context.Context, req *membershipv1.ScheduleVacationRequest) (*membershipv1.ScheduleVacationResponse, error) {
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Employee(id)); err != nil {
		return nil, err
	}
	cmd := membership.ScheduleVacationCommand{EmployeeID: id}
	if ts := req.GetStartsAt(); ts != nil {
		cmd.StartsAt = ts.AsTime()
	}
	if ts := req.GetEndsAt(); ts != nil {
		t := ts.AsTime()
		cmd.EndsAt = &t
	}
	res, err := h.empSvc.ScheduleVacation(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &membershipv1.ScheduleVacationResponse{VacationId: res.ID.String()}, nil
}
