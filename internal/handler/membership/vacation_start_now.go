package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// StartVacationNow translates a gRPC StartVacationNowRequest into a
// service command and returns the new vacation ID on success.
func (h *MembershipHandler) StartVacationNow(ctx context.Context, req *membershipv1.StartVacationNowRequest) (*membershipv1.StartVacationNowResponse, error) {
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
	cmd := membership.StartVacationNowCommand{EmployeeID: id}
	if ts := req.GetEndsAt(); ts != nil {
		t := ts.AsTime()
		cmd.EndsAt = &t
	}
	res, err := h.empSvc.StartVacationNow(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &membershipv1.StartVacationNowResponse{VacationId: res.ID.String()}, nil
}
