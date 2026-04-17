package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// TerminateEmployee translates a gRPC TerminateEmployeeRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) TerminateEmployee(ctx context.Context, req *membershipv1.TerminateEmployeeRequest) (*membershipv1.TerminateEmployeeResponse, error) {
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaEmployee(ctx, callerID, id); err != nil {
		return nil, err
	}
	if err := h.empSvc.Terminate(ctx, membership.TerminateEmployeeCommand{ID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.TerminateEmployeeResponse{}, nil
}
