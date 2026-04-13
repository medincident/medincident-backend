package membership

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// UpdateEmployeePosition translates a gRPC UpdateEmployeePositionRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeePosition(ctx context.Context, req *membershipv1.UpdateEmployeePositionRequest) (*membershipv1.UpdateEmployeePositionResponse, error) {
	id, err := uuid.Parse(req.GetEmployeeId())
	if err != nil {
		return nil, oops.In("handler.membership").
			Code("employee_id_invalid").
			Public("employee_id is not a valid UUID.").
			Wrap(err)
	}
	if err := h.empSvc.UpdatePosition(ctx, membership.UpdateEmployeePositionCommand{
		ID:       id,
		Position: req.Position,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateEmployeePositionResponse{}, nil
}
