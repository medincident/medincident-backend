package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) CreateClinic(
	ctx context.Context,
	req *orgstructurev1.CreateClinicRequest,
) (*orgstructurev1.CreateClinicResponse, error) {
	orgID, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	callerID, err := middleware.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.RequireOrgAdmin(ctx, callerID, orgID); err != nil {
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
