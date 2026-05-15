package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

func (h *MembershipHandler) DeactivateEmployee(
	ctx context.Context,
	req *membershipv1.DeactivateEmployeeRequest,
) (*membershipv1.DeactivateEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.Deactivate(ctx, membership.DeactivateEmployeeCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: membership.DeactivateEmployeePayload{ID: req.GetEmployeeId()},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.DeactivateEmployeeResponse{}, nil
}

func (h *MembershipHandler) ActivateEmployee(
	ctx context.Context,
	req *membershipv1.ActivateEmployeeRequest,
) (*membershipv1.ActivateEmployeeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.Activate(ctx, membership.ActivateEmployeeCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: membership.ActivateEmployeePayload{ID: req.GetEmployeeId()},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.ActivateEmployeeResponse{}, nil
}
