package membership

import (
	"context"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/middleware"
	"github.com/medincident/medincident-command-service/internal/service/membership"
)

// ForceEndVacation translates a gRPC ForceEndVacationRequest into a
// service command and returns an empty response on success.
func (h *MembershipHandler) ForceEndVacation(ctx context.Context, req *membershipv1.ForceEndVacationRequest) (*membershipv1.ForceEndVacationResponse, error) {
	id, err := parseVacationID(req.GetVacationId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdminViaVacation(ctx, callerID, id); err != nil {
		return nil, err
	}
	if err := h.empSvc.ForceEndVacation(ctx, membership.ForceEndVacationCommand{VacationID: id}); err != nil {
		return nil, err
	}
	return &membershipv1.ForceEndVacationResponse{}, nil
}
