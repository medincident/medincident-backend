package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignDepartmentResponsibleDeputy translates a gRPC request into a
// service command and sets the deputy slot on an existing DR role.
func (h *MembershipHandler) AssignDepartmentResponsibleDeputy(ctx context.Context, req *membershipv1.AssignDepartmentResponsibleDeputyRequest) (*membershipv1.AssignDepartmentResponsibleDeputyResponse, error) {
	depID, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	deputyID, err := parseDeputyEmployeeID(req.GetDeputyEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignDepartmentResponsibleDeputy(ctx, membership.AssignDepartmentResponsibleDeputyCommand{
		DepartmentID:     depID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignDepartmentResponsibleDeputyResponse{}, nil
}
