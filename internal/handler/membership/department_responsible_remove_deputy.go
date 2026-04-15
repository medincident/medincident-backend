package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// RemoveDepartmentResponsibleDeputy translates a gRPC request into a
// service command and clears the deputy slot on a DR role.
func (h *MembershipHandler) RemoveDepartmentResponsibleDeputy(ctx context.Context, req *membershipv1.RemoveDepartmentResponsibleDeputyRequest) (*membershipv1.RemoveDepartmentResponsibleDeputyResponse, error) {
	var ids idErrs
	depID := ids.parse(req.GetDepartmentId(), parseDepartmentID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveDepartmentResponsibleDeputy(ctx, membership.RemoveDepartmentResponsibleDeputyCommand{
		DepartmentID: depID,
		EmployeeID:   empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveDepartmentResponsibleDeputyResponse{}, nil
}
