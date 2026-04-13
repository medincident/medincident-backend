package membership

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// UpdateEmployeeDepartment translates a gRPC UpdateEmployeeDepartmentRequest
// into a service command and returns an empty response on success.
func (h *MembershipHandler) UpdateEmployeeDepartment(ctx context.Context, req *membershipv1.UpdateEmployeeDepartmentRequest) (*membershipv1.UpdateEmployeeDepartmentResponse, error) {
	empID, err := uuid.Parse(req.GetEmployeeId())
	if err != nil {
		return nil, oops.In("handler.membership").Code("employee_id_invalid").Public("employee_id is not a valid UUID.").Wrap(err)
	}
	depID, err := uuid.Parse(req.GetDepartmentId())
	if err != nil {
		return nil, oops.In("handler.membership").Code("department_id_invalid").Public("department_id is not a valid UUID.").Wrap(err)
	}
	if err := h.empSvc.UpdateDepartment(ctx, membership.UpdateEmployeeDepartmentCommand{
		ID:           empID,
		DepartmentID: depID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.UpdateEmployeeDepartmentResponse{}, nil
}
