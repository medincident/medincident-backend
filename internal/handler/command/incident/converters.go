package incident

import (
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/model"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

const (
	errCodeHandlerInvalidPriority = "handler_invalid_priority"
	errCodeHandlerInvalidStatus   = "handler_invalid_status"
)

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
