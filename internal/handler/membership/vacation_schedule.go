package membership

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// ScheduleVacation translates a gRPC ScheduleVacationRequest into a
// service command and returns the new vacation ID on success.
func (h *MembershipHandler) ScheduleVacation(ctx context.Context, req *membershipv1.ScheduleVacationRequest) (*membershipv1.ScheduleVacationResponse, error) {
	id, err := uuid.Parse(req.GetEmployeeId())
	if err != nil {
		return nil, oops.In("handler.membership").
			Code("employee_id_invalid").
			Public("employee_id is not a valid UUID.").
			Wrap(err)
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
