package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// RemoveDepartmentResponsibleDeputy translates a gRPC request into a
// service command and clears the deputy slot on a DR role.
func (h *MembershipHandler) RemoveDepartmentResponsibleDeputy(ctx context.Context, req *membershipv1.RemoveDepartmentResponsibleDeputyRequest) (*membershipv1.RemoveDepartmentResponsibleDeputyResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RemoveDepartmentResponsibleDeputy(ctx, membership.RemoveDepartmentResponsibleDeputyCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RemoveDepartmentResponsibleDeputyPayload{
			DepartmentID: req.GetDepartmentId(),
			EmployeeID:   req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RemoveDepartmentResponsibleDeputyResponse{}, nil
}
