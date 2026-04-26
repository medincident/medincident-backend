package buffer

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

func (h *BufferHandler) RejectPatientIncident(
	ctx context.Context, req *bufferv1.RejectPatientIncidentRequest,
) (*bufferv1.RejectPatientIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.Reject(ctx, buffersvc.RejectCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: buffersvc.RejectPayload{BufferID: req.GetBufferId()},
	}); err != nil {
		return nil, err
	}
	return &bufferv1.RejectPatientIncidentResponse{}, nil
}
