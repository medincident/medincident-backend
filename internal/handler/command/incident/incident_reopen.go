package incident

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

func (h *IncidentHandler) ReopenIncident(
	ctx context.Context, req *incidentv1.ReopenIncidentRequest,
) (*incidentv1.ReopenIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.svc.Reopen(ctx, incidentsvc.ReopenIncidentCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: incidentsvc.ReopenIncidentPayload{IncidentID: req.GetIncidentId()},
	})
	if err != nil {
		return nil, err
	}
	return &incidentv1.ReopenIncidentResponse{ReopenedIncidentId: res.NewIncidentID.String()}, nil
}
