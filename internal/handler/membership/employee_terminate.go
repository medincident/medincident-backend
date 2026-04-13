package membership

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// TerminateEmployee translates a gRPC TerminateEmployeeRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) TerminateEmployee(ctx context.Context, req *membershipv1.TerminateEmployeeRequest) (*membershipv1.TerminateEmployeeResponse, error) {
	id, err := uuid.Parse(req.GetEmployeeId())
	if err != nil {
		return nil, oops.In("handler.membership").
			Code("employee_id_invalid").
			Public("employee_id is not a valid UUID.").
			Wrap(err)
	}
	if err := h.empSvc.Terminate(ctx, membership.TerminateEmployeeCommand{ID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.TerminateEmployeeResponse{}, nil
}
