package request

import (
	"github.com/samber/oops"

	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

const errCodeHandlerInvalidStatus = "handler_invalid_status"

func protoStatusToString(s requestv1.ServiceRequestStatus) (string, error) {
	switch s {
	case requestv1.ServiceRequestStatus_SERVICE_REQUEST_STATUS_IN_WORK:
		return "in_work", nil
	case requestv1.ServiceRequestStatus_SERVICE_REQUEST_STATUS_ON_HOLD:
		return "on_hold", nil
	case requestv1.ServiceRequestStatus_SERVICE_REQUEST_STATUS_PENDING_REVIEW:
		return "pending_review", nil
	case requestv1.ServiceRequestStatus_SERVICE_REQUEST_STATUS_COMPLETED:
		return "completed", nil
	case requestv1.ServiceRequestStatus_SERVICE_REQUEST_STATUS_CANCELLED:
		return "cancelled", nil
	default:
		return "", oops.In("handler.command.request").
			Code(errCodeHandlerInvalidStatus).
			Public("Invalid status.").Errorf("unsupported status %v for UpdateServiceRequestStatus", s)
	}
}
