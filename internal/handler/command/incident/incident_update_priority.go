package incident

import (
	"context"

	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

const errCodeHandlerInvalidPriority = "handler_invalid_priority"

func protoPriorityToModel(p incidentv1.IncidentPriority) (model.IncidentPriority, error) {
	switch p {
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_LOW:
		return model.IncidentPriorityLow, nil
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_NORMAL:
		return model.IncidentPriorityNormal, nil
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_HIGH:
		return model.IncidentPriorityHigh, nil
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_CRITICAL:
		return model.IncidentPriorityCritical, nil
	default:
		return "", oops.In("handler.command.incident").
			Code(errCodeHandlerInvalidPriority).
			Public("Invalid priority.").Errorf("unknown priority %v", p)
	}
}

func (h *IncidentHandler) UpdateIncidentPriority(
	ctx context.Context, req *incidentv1.UpdateIncidentPriorityRequest,
) (*incidentv1.UpdateIncidentPriorityResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	pr, err := protoPriorityToModel(req.GetPriority())
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdatePriority(ctx, &incidentsvc.UpdateIncidentPriorityCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: incidentsvc.UpdateIncidentPriorityPayload{
			IncidentID: req.GetIncidentId(),
			Priority:   pr,
		},
	}); err != nil {
		return nil, err
	}
	return &incidentv1.UpdateIncidentPriorityResponse{}, nil
}
