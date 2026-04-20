package membership

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// AssignClinicHeadDeputy translates a gRPC request into a service
// command and sets the deputy slot on an existing CH role.
func (h *MembershipHandler) AssignClinicHeadDeputy(ctx context.Context, req *membershipv1.AssignClinicHeadDeputyRequest) (*membershipv1.AssignClinicHeadDeputyResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignClinicHeadDeputy(ctx, membership.AssignClinicHeadDeputyCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignClinicHeadDeputyPayload{
			ClinicID:         req.GetClinicId(),
			EmployeeID:       req.GetEmployeeId(),
			DeputyEmployeeID: req.GetDeputyEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignClinicHeadDeputyResponse{}, nil
}
