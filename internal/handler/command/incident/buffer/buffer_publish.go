package buffer

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	buffersvc "github.com/medincident/medincident-backend/internal/service/command/incident/buffer"
	bufferv1 "github.com/medincident/medincident-backend/pkg/command/incident/buffer/v1"
)

func (h *BufferHandler) PublishPatientIncident(
	ctx context.Context, req *bufferv1.PublishPatientIncidentRequest,
) (*bufferv1.PublishPatientIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.svc.Publish(ctx, &buffersvc.PublishCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: buffersvc.PublishPayload{
			BufferID:     req.GetBufferId(),
			DepartmentID: req.GetDepartmentId(),
			CategoryID:   req.GetCategoryId(),
			TypeID:       req.GetTypeId(),
			Description:  req.Description,
		},
	})
	if err != nil {
		return nil, err
	}
	return &bufferv1.PublishPatientIncidentResponse{IncidentId: res.IncidentID.String()}, nil
}
