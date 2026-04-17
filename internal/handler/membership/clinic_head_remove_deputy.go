package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// RemoveClinicHeadDeputy translates a gRPC request into a service
// command and clears the deputy slot on a CH role.
func (h *MembershipHandler) RemoveClinicHeadDeputy(ctx context.Context, req *membershipv1.RemoveClinicHeadDeputyRequest) (*membershipv1.RemoveClinicHeadDeputyResponse, error) {
	var ids idErrs
	clinicID := ids.parse(req.GetClinicId(), parseClinicID)
	empID := ids.parse(req.GetEmployeeId(), parseEmployeeID)
	if err := ids.err(); err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveClinicHeadDeputy(ctx, membership.RemoveClinicHeadDeputyCommand{
		ClinicID:   clinicID,
		EmployeeID: empID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveClinicHeadDeputyResponse{}, nil
}
