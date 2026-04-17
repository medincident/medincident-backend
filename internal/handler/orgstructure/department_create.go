package orgstructure

import (
	"context"

	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/service/orgstructure/v1"
)

func (h *OrgStructureHandler) CreateDepartment(
	ctx context.Context,
	req *orgstructurev1.CreateDepartmentRequest,
) (*orgstructurev1.CreateDepartmentResponse, error) {
	clinicID, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	result, err := h.deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		ClinicID:    clinicID,
		Name:        req.GetName(),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &orgstructurev1.CreateDepartmentResponse{DepartmentId: result.ID.String()}, nil
}
