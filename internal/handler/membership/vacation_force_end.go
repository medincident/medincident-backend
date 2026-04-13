package membership

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// ForceEndVacation translates a gRPC ForceEndVacationRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) ForceEndVacation(ctx context.Context, req *membershipv1.ForceEndVacationRequest) (*membershipv1.ForceEndVacationResponse, error) {
	id, err := uuid.Parse(req.GetVacationId())
	if err != nil {
		return nil, oops.In("handler.membership").
			Code("vacation_id_invalid").
			Public("vacation_id is not a valid UUID.").
			Wrap(err)
	}
	if err := h.empSvc.ForceEndVacation(ctx, membership.ForceEndVacationCommand{VacationID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.ForceEndVacationResponse{}, nil
}
