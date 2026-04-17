package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/service/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/service/membership/v1"
)

// RevokeDepartmentResponsible translates a gRPC request into a service
// command and removes the employee's department responsible role.
func (h *MembershipHandler) RevokeDepartmentResponsible(ctx context.Context, req *membershipv1.RevokeDepartmentResponsibleRequest) (*membershipv1.RevokeDepartmentResponsibleResponse, error) {
	var ids idErrs
	depID := ids.parse(req.GetDepartmentId(), parseDepartmentID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeDepartmentResponsible(ctx, membership.RevokeDepartmentResponsibleCommand{
		DepartmentID: depID,
		EmployeeID:   empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeDepartmentResponsibleResponse{}, nil
}
