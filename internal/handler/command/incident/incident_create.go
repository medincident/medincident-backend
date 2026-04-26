package incident

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

func (h *IncidentHandler) CreateIncident(
	ctx context.Context, req *incidentv1.CreateIncidentRequest,
) (*incidentv1.CreateIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.svc.Create(ctx, &incidentsvc.CreateIncidentCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: incidentsvc.CreateIncidentPayload{
			DepartmentID: req.GetDepartmentId(),
			CategoryID:   req.GetCategoryId(),
			TypeID:       req.GetTypeId(),
			Description:  req.Description,
			OccurredAt:   req.GetOccurredAt(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &incidentv1.CreateIncidentResponse{IncidentId: res.ID.String()}, nil
}
