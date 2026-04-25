// Package buffer is the gRPC transport for IncidentBufferCommandService.
package buffer

import (
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

// BufferHandler implements bufferv1.IncidentBufferCommandServiceServer.
type BufferHandler struct {
	bufferv1.UnimplementedIncidentBufferCommandServiceServer
	svc *buffersvc.BufferService
}

func NewBufferHandler(svc *buffersvc.BufferService) *BufferHandler {
	return &BufferHandler{svc: svc}
}
