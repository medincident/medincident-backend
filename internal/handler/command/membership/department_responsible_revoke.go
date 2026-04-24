package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// RevokeDepartmentResponsible translates a gRPC request into a service
// command and removes the employee's department responsible role.
func (h *MembershipHandler) RevokeDepartmentResponsible(ctx context.Context, req *membershipv1.RevokeDepartmentResponsibleRequest) (*membershipv1.RevokeDepartmentResponsibleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.RevokeDepartmentResponsible(ctx, membership.RevokeDepartmentResponsibleCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.RevokeDepartmentResponsiblePayload{
			DepartmentID: req.GetDepartmentId(),
			EmployeeID:   req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.RevokeDepartmentResponsibleResponse{}, nil
}
