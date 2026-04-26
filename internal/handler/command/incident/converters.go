package incident

import (
	"github.com/samber/oops"

	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

const (
	errCodeHandlerInvalidPriority = "handler_invalid_priority"
	errCodeHandlerInvalidStatus   = "handler_invalid_status"
)

func protoPriorityToString(p incidentv1.IncidentPriority) (string, error) {
	switch p {
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_LOW:
		return "low", nil
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_NORMAL:
		return "normal", nil
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_HIGH:
		return "high", nil
	case incidentv1.IncidentPriority_INCIDENT_PRIORITY_CRITICAL:
		return "critical", nil
	default:
		return "", oops.In("handler.command.incident").
			Code(errCodeHandlerInvalidPriority).
			Public("Invalid priority.").Errorf("unknown priority %v", p)
	}
}

func protoStatusToString(s incidentv1.IncidentStatus) (string, error) {
	switch s {
	case incidentv1.IncidentStatus_INCIDENT_STATUS_PENDING:
		return "pending", nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_IN_PROGRESS:
		return "in_progress", nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_DONE:
		return "done", nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_REJECTED:
		return "rejected", nil
	case incidentv1.IncidentStatus_INCIDENT_STATUS_CANCELLED:
		return "cancelled", nil
	default:
		return "", oops.In("handler.command.incident").
			Code(errCodeHandlerInvalidStatus).
			Public("Invalid status.").Errorf("unknown status %v", s)
	}
}
