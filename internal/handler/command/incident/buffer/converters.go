package buffer

import (
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

// protoBufferPriorityToString maps the proto enum to the domain string
// used by SubmitPayload and UpdatePayload (validate:"oneof=normal high").
// UNSPECIFIED and unknown values map to "" so that service-layer validation
// rejects them with validation_failed (HTTP 400) rather than a handler error.
func protoBufferPriorityToString(p bufferv1.BufferPriority) string {
	switch p {
	case bufferv1.BufferPriority_BUFFER_PRIORITY_NORMAL:
		return "normal"
	case bufferv1.BufferPriority_BUFFER_PRIORITY_HIGH:
		return "high"
	default:
		return ""
	}
}
