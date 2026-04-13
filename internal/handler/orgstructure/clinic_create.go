package orgstructure

import (
	"context"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/orgstructure/v1"
	orgsvc "github.com/medincident/medincident-command-service/internal/services/orgstructure"
)

func (h *OrgStructureHandler) CreateClinic(
	ctx context.Context,
	req *orgstructurev1.CreateClinicRequest,
) (*orgstructurev1.CreateClinicResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		OrganizationID:  orgID,
		Name:            req.GetName(),
		Description:     req.Description,
		PhysicalAddress: addressInputFromProto(req.GetPhysicalAddress()),
	})
	if err != nil {
		return nil, err
	}
	return &orgstructurev1.CreateClinicResponse{ClinicId: result.ID.String()}, nil
}
