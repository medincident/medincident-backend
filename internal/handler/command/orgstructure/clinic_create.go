package orgstructure

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-backend/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) CreateClinic(
	ctx context.Context,
	req *orgstructurev1.CreateClinicRequest,
) (*orgstructurev1.CreateClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.CreateClinicPayload{
			OrganizationID:  req.GetOrganizationId(),
			Name:            req.GetName(),
			Description:     req.Description,
			PhysicalAddress: addressInputFromProto(req.GetPhysicalAddress()),
		},
	})
	if err != nil {
		return nil, err
	}
	return &orgstructurev1.CreateClinicResponse{ClinicId: result.ID.String()}, nil
}
