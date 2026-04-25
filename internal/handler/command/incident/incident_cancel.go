package incident

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

func (h *IncidentHandler) CancelIncident(
	ctx context.Context, req *incidentv1.CancelIncidentRequest,
) (*incidentv1.CancelIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.Cancel(ctx, incidentsvc.CancelIncidentCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: incidentsvc.CancelIncidentPayload{IncidentID: req.GetIncidentId()},
	}); err != nil {
		return nil, err
	}
	return &incidentv1.CancelIncidentResponse{}, nil
}
