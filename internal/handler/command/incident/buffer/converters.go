package buffer

import (
	"github.com/samber/oops"

	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

const errCodeHandlerInvalidBufferPriority = "handler_invalid_buffer_priority"

// protoBufferPriorityToString maps the proto enum to the domain string
// used by SubmitPayload and UpdatePayload (validate:"oneof=normal high").
func protoBufferPriorityToString(p bufferv1.BufferPriority) (string, error) {
	switch p {
	case bufferv1.BufferPriority_BUFFER_PRIORITY_NORMAL:
		return "normal", nil
	case bufferv1.BufferPriority_BUFFER_PRIORITY_HIGH:
		return "high", nil
	default:
		return "", oops.In("handler.command.incident.buffer").
			Code(errCodeHandlerInvalidBufferPriority).
			Public("Invalid buffer priority.").Errorf("unknown: %v", p)
	}
}
