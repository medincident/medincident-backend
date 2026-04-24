package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// UpdateEmployeePosition translates a gRPC UpdateEmployeePositionRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeePosition(ctx context.Context, req *membershipv1.UpdateEmployeePositionRequest) (*membershipv1.UpdateEmployeePositionResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.UpdatePosition(ctx, membership.UpdateEmployeePositionCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.UpdateEmployeePositionPayload{
			ID:       req.GetEmployeeId(),
			Position: req.Position,
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateEmployeePositionResponse{}, nil
}
