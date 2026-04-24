package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) MoveIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.MoveIncidentTypeRequest,
) (*incidentclassifierv1.MoveIncidentTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.Move(ctx, classifiersvc.MoveIncidentTypeCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.MoveIncidentTypePayload{
			TypeID:        req.GetTypeId(),
			NewCategoryID: req.GetNewCategoryId(),
		},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.MoveIncidentTypeResponse{}, nil
}
