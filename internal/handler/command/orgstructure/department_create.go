package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) CreateDepartment(
	ctx context.Context,
	req *orgstructurev1.CreateDepartmentRequest,
) (*orgstructurev1.CreateDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.CreateDepartmentPayload{
			ClinicID:    req.GetClinicId(),
			Name:        req.GetName(),
			Description: req.Description,
		},
	})
	if err != nil {
		return nil, err
	}
	return &orgstructurev1.CreateDepartmentResponse{DepartmentId: result.ID.String()}, nil
}
