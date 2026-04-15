package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// CancelScheduledVacation translates a gRPC CancelScheduledVacationRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) CancelScheduledVacation(ctx context.Context, req *membershipv1.CancelScheduledVacationRequest) (*membershipv1.CancelScheduledVacationResponse, error) {
	id, err := parseVacationID(req.GetVacationId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.CancelScheduledVacation(ctx, membership.CancelScheduledVacationCommand{VacationID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.CancelScheduledVacationResponse{}, nil
}
