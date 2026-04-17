package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// RevokeClinicHead translates a gRPC request into a service command
// and removes the employee's clinic head role.
func (h *MembershipHandler) RevokeClinicHead(ctx context.Context, req *membershipv1.RevokeClinicHeadRequest) (*membershipv1.RevokeClinicHeadResponse, error) {
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
	if err := h.empSvc.RevokeClinicHead(ctx, membership.RevokeClinicHeadCommand{
		ClinicID:   clinicID,
		EmployeeID: empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeClinicHeadResponse{}, nil
}
