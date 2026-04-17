package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// AssignClinicHeadDeputy translates a gRPC request into a service
// command and sets the deputy slot on an existing CH role.
func (h *MembershipHandler) AssignClinicHeadDeputy(ctx context.Context, req *membershipv1.AssignClinicHeadDeputyRequest) (*membershipv1.AssignClinicHeadDeputyResponse, error) {
	var ids idErrs
	clinicID := ids.parse(req.GetClinicId(), parseClinicID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	deputyID := ids.parse(req.GetDeputyEmployeeId(), parseDeputyEmployeeID)
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
	if err := h.empSvc.AssignClinicHeadDeputy(ctx, membership.AssignClinicHeadDeputyCommand{
		ClinicID:         clinicID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignClinicHeadDeputyResponse{}, nil
}
