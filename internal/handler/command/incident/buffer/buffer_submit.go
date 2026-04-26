package buffer

import (
	"context"
	"time"

	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

const errCodeHandlerInvalidOccurredAt = "handler_invalid_occurred_at"

// parseOptionalTimestamp parses a nullable RFC3339Nano timestamp string.
// Returns (zero, false, nil) when raw is nil; (t, true, nil) on success;
// (zero, false, err) on parse failure.
func parseOptionalTimestamp(raw *string, field string) (time.Time, bool, error) {
	if raw == nil {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(time.RFC3339Nano, *raw)
	if err != nil {
		return time.Time{}, false, oops.In("handler.command.incident.buffer").
			Code(errCodeHandlerInvalidOccurredAt).
			Public(field+" is not a valid RFC3339 timestamp.").
			With(field, *raw).Wrap(err)
	}
	return t, true, nil
}

func (h *BufferHandler) SubmitPatientIncident(
	ctx context.Context, req *bufferv1.SubmitPatientIncidentRequest,
) (*bufferv1.SubmitPatientIncidentResponse, error) {
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
	res, err := h.svc.Submit(ctx, buffersvc.SubmitCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: buffersvc.SubmitPayload{
			OrganizationID: req.GetOrganizationId(),
			CategoryID:     req.CategoryId,
			TypeID:         req.TypeId,
			Description:    req.Description,
			OccurredAt:     occurred,
		},
	})
	if err != nil {
		return nil, err
	}
	return &bufferv1.SubmitPatientIncidentResponse{BufferId: res.ID.String()}, nil
}
