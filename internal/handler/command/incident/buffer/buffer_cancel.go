package buffer

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

func (h *BufferHandler) CancelPatientIncident(
	ctx context.Context, req *bufferv1.CancelPatientIncidentRequest,
) (*bufferv1.CancelPatientIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.Cancel(ctx, &buffersvc.CancelCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: buffersvc.CancelPayload{BufferID: req.GetBufferId()},
	}); err != nil {
		return nil, err
	}
	return &bufferv1.CancelPatientIncidentResponse{}, nil
}
