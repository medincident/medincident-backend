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

const errCodeHandlerInvalidStatus = "handler_invalid_status"

func protoStatusToModel(s incidentv1.IncidentStatus) (model.IncidentStatus, error) {
	switch s {
	case incidentv1.IncidentStatus_INCIDENT_STATUS_PENDING:
		return model.IncidentStatusPending, nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_IN_PROGRESS:
		return model.IncidentStatusInProgress, nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_DONE:
		return model.IncidentStatusDone, nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_REJECTED:
		return model.IncidentStatusRejected, nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_CANCELLED:
		return model.IncidentStatusCancelled, nil
	default:
		return "", oops.In("handler.command.incident").
			Code(errCodeHandlerInvalidStatus).
			Public("Invalid status.").Errorf("unknown status %v", s)
	}
}

func (h *IncidentHandler) UpdateIncidentStatus(
	ctx context.Context, req *incidentv1.UpdateIncidentStatusRequest,
) (*incidentv1.UpdateIncidentStatusResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	st, err := protoStatusToModel(req.GetNewStatus())
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdateStatus(ctx, &incidentsvc.UpdateIncidentStatusCommand{
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
