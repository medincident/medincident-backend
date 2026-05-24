package buffer

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

func (h *BufferHandler) SubmitPatientIncident(
	ctx context.Context, req *bufferv1.SubmitPatientIncidentRequest,
) (*bufferv1.SubmitPatientIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.svc.Submit(ctx, buffersvc.SubmitCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: buffersvc.SubmitPayload{
			OrganizationID: req.GetOrganizationId(),
			CategoryID:     req.CategoryId,
			TypeID:         req.TypeId,
			Description:    req.GetDescription(),
			Summary:        req.GetSummary(),
			Priority:       protoBufferPriorityToString(req.GetPriority()),
			OccurredAt:     req.OccurredAt,
		},
	})
	if err != nil {
		return nil, err
	}
	return &bufferv1.SubmitPatientIncidentResponse{BufferId: res.ID.String()}, nil
}
