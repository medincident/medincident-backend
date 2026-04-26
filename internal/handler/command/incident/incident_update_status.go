package incident

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

func (h *IncidentHandler) UpdateIncidentStatus(
	ctx context.Context, req *incidentv1.UpdateIncidentStatusRequest,
) (*incidentv1.UpdateIncidentStatusResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	st, err := protoStatusToString(req.GetNewStatus())
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdateStatus(ctx, incidentsvc.UpdateIncidentStatusCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: incidentsvc.UpdateIncidentStatusPayload{
			IncidentID: req.GetIncidentId(),
			NewStatus:  st,
		},
	}); err != nil {
		return nil, err
	}
	return &incidentv1.UpdateIncidentStatusResponse{}, nil
}
