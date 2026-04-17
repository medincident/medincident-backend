package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// TerminateEmployee translates a gRPC TerminateEmployeeRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) TerminateEmployee(ctx context.Context, req *membershipv1.TerminateEmployeeRequest) (*membershipv1.TerminateEmployeeResponse, error) {
	id, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.Terminate(ctx, membership.TerminateEmployeeCommand{ID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.TerminateEmployeeResponse{}, nil
}
