package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// AssignClinicHeadDeputy translates a gRPC request into a service
// command and sets the deputy slot on an existing CH role.
func (h *MembershipHandler) AssignClinicHeadDeputy(ctx context.Context, req *membershipv1.AssignClinicHeadDeputyRequest) (*membershipv1.AssignClinicHeadDeputyResponse, error) {
	clinicID, err := parseClinicID(req.GetClinicId())
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
	if err := h.empSvc.AssignClinicHeadDeputy(ctx, membership.AssignClinicHeadDeputyCommand{
		ClinicID:         clinicID,
		EmployeeID:       empID,
		DeputyEmployeeID: deputyID,
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignClinicHeadDeputyResponse{}, nil
}
