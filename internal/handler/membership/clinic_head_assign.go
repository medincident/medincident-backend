package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignClinicHead translates a gRPC request into a service command
// and links the employee to the clinic as its head.
func (h *MembershipHandler) AssignClinicHead(ctx context.Context, req *membershipv1.AssignClinicHeadRequest) (*membershipv1.AssignClinicHeadResponse, error) {
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	empID, err := parseEmployeeID(req.GetEmployeeId())
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignClinicHead(ctx, membership.AssignClinicHeadCommand{
		ClinicID:   clinicID,
		EmployeeID: empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignClinicHeadResponse{}, nil
}
