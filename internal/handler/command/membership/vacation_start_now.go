package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// StartVacationNow translates a gRPC StartVacationNowRequest into a
// service command and returns the new vacation ID on success.
func (h *MembershipHandler) StartVacationNow(ctx context.Context, req *membershipv1.StartVacationNowRequest) (*membershipv1.StartVacationNowResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	payload := membership.StartVacationNowPayload{
		EmployeeID: req.GetEmployeeId(),
	}
	if ts := req.GetEndsAt(); ts != nil {
		t := ts.AsTime()
		payload.EndsAt = &t
	}
	res, err := h.empSvc.StartVacationNow(ctx, membership.StartVacationNowCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}
	return &membershipv1.StartVacationNowResponse{VacationId: res.ID.String()}, nil
}
