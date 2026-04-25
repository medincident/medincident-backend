package incident

import (
	"context"
	"time"

	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

const errCodeHandlerInvalidOccurredAt = "handler_invalid_occurred_at"

func parseTimestamp(raw, field string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return time.Time{}, oops.In("handler.command.incident").
			Code(errCodeHandlerInvalidOccurredAt).
			Public(field+" is not a valid RFC3339 timestamp.").
			With(field, raw).Wrap(err)
	}
	return t, nil
}

func (h *IncidentHandler) CreateIncident(
	ctx context.Context, req *incidentv1.CreateIncidentRequest,
) (*incidentv1.CreateIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	occurred, err := parseTimestamp(req.GetOccurredAt(), "occurred_at")
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
			OccurredAt:   occurred,
		},
	})
	if err != nil {
		return nil, err
	}
	return &incidentv1.CreateIncidentResponse{IncidentId: res.ID.String()}, nil
}
