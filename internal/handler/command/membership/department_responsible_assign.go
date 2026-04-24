package membership

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// AssignDepartmentResponsible translates a gRPC request into a service
// command and links the employee to the department as its responsible.
func (h *MembershipHandler) AssignDepartmentResponsible(ctx context.Context, req *membershipv1.AssignDepartmentResponsibleRequest) (*membershipv1.AssignDepartmentResponsibleResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.empSvc.AssignDepartmentResponsible(ctx, membership.AssignDepartmentResponsibleCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: membership.AssignDepartmentResponsiblePayload{
			DepartmentID: req.GetDepartmentId(),
			EmployeeID:   req.GetEmployeeId(),
		},
	}); err != nil {
		return nil, err
	}
	return &membershipv1.AssignDepartmentResponsibleResponse{}, nil
}
