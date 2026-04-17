package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// UpdateEmployeePosition translates a gRPC UpdateEmployeePositionRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeePosition(ctx context.Context, req *membershipv1.UpdateEmployeePositionRequest) (*membershipv1.UpdateEmployeePositionResponse, error) {
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.UpdatePosition(ctx, membership.UpdateEmployeePositionCommand{
		ID:       id,
		Position: req.Position,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateEmployeePositionResponse{}, nil
}
