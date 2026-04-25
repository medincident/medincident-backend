package buffer

import (
	"context"
	"time"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

func (h *BufferHandler) UpdatePatientIncident(
	ctx context.Context, req *bufferv1.UpdatePatientIncidentRequest,
) (*bufferv1.UpdatePatientIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	occurredVal, occurredSet, err := parseOptionalTimestamp(req.OccurredAt, "occurred_at")
	if err != nil {
		return nil, err
	}
	var occurred *time.Time
	if occurredSet {
		occurred = &occurredVal
	}
	if err := h.svc.Update(ctx, &buffersvc.UpdateCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: buffersvc.UpdatePayload{
			BufferID:    req.GetBufferId(),
			CategoryID:  req.CategoryId,
			TypeID:      req.TypeId,
			Description: req.Description,
			OccurredAt:  occurred,
		},
	}); err != nil {
		return nil, err
	}
	return &bufferv1.UpdatePatientIncidentResponse{}, nil
}
