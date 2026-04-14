package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignDepartmentResponsible translates a gRPC request into a service
// command and links the employee to the department as its responsible.
func (h *MembershipHandler) AssignDepartmentResponsible(ctx context.Context, req *membershipv1.AssignDepartmentResponsibleRequest) (*membershipv1.AssignDepartmentResponsibleResponse, error) {
	depID, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignDepartmentResponsible(ctx, membership.AssignDepartmentResponsibleCommand{
		DepartmentID: depID,
		EmployeeID:   empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignDepartmentResponsibleResponse{}, nil
}
