package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// UpdateEmployeePosition translates a gRPC UpdateEmployeePositionRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeePosition(ctx context.Context, req *membershipv1.UpdateEmployeePositionRequest) (*membershipv1.UpdateEmployeePositionResponse, error) {
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Employee(id)); err != nil {
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
