package incident

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

func (h *IncidentHandler) UpdateIncidentPriority(
	ctx context.Context, req *incidentv1.UpdateIncidentPriorityRequest,
) (*incidentv1.UpdateIncidentPriorityResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdatePriority(ctx, incidentsvc.UpdateIncidentPriorityCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: incidentsvc.UpdateIncidentPriorityPayload{
			IncidentID: req.GetIncidentId(),
			Priority:   req.GetPriority().String(),
		},
	}); err != nil {
		return nil, err
	}
	return &incidentv1.UpdateIncidentPriorityResponse{}, nil
}
