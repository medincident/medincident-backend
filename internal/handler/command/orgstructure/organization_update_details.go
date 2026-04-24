package orgstructure

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	orgsvc "github.com/medincident/medincident-command-service/internal/service/command/orgstructure"
	orgstructurev1 "github.com/medincident/medincident-command-service/pkg/command/orgstructure/v1"
)

func (h *OrgStructureHandler) UpdateOrganizationDetails(
	ctx context.Context,
	req *orgstructurev1.UpdateOrganizationDetailsRequest,
) (*orgstructurev1.UpdateOrganizationDetailsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.orgSvc.UpdateDetails(ctx, orgsvc.UpdateOrganizationDetailsCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: orgsvc.UpdateOrganizationDetailsPayload{
			ID:          req.GetOrganizationId(),
			Name:        req.GetName(),
			Description: req.Description,
		},
	}); err != nil {
		return nil, err
	}
	return &orgstructurev1.UpdateOrganizationDetailsResponse{}, nil
}
