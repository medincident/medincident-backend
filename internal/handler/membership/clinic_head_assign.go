package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// AssignClinicHead translates a gRPC request into a service command
// and links the employee to the clinic as its head.
func (h *MembershipHandler) AssignClinicHead(ctx context.Context, req *membershipv1.AssignClinicHeadRequest) (*membershipv1.AssignClinicHeadResponse, error) {
	var ids idErrs
	clinicID := ids.parse(req.GetClinicId(), parseClinicID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaClinic(ctx, callerID, clinicID); err != nil {
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
