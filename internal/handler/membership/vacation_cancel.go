package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// CancelScheduledVacation translates a gRPC CancelScheduledVacationRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) CancelScheduledVacation(ctx context.Context, req *membershipv1.CancelScheduledVacationRequest) (*membershipv1.CancelScheduledVacationResponse, error) {
	id, err := parseVacationID(req.GetVacationId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Vacation(id)); err != nil {
		return nil, err
	}
	if err := h.empSvc.CancelScheduledVacation(ctx, membership.CancelScheduledVacationCommand{VacationID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.CancelScheduledVacationResponse{}, nil
}
