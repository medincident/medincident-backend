package membership

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// HireEmployee translates a gRPC HireEmployeeRequest into a service
// command and returns the new employee ID on success.
func (h *MembershipHandler) HireEmployee(ctx context.Context, req *membershipv1.HireEmployeeRequest) (*membershipv1.HireEmployeeResponse, error) {
	depID, err := uuid.Parse(req.GetDepartmentId())
	if err != nil {
		return nil, oops.In("handler.membership").
			Code("department_id_invalid").
			Public("department_id is not a valid UUID.").
			Wrap(err)
	}
	cmd := membership.HireEmployeeCommand{
		ZitadelUserID: req.GetZitadelUserId(),
		DepartmentID:  depID,
		Position:      req.Position,
	}
	result, err := h.empSvc.Hire(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &membershipv1.HireEmployeeResponse{EmployeeId: result.ID.String()}, nil
}
