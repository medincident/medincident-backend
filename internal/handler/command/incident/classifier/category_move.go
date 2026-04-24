package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) MoveIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.MoveIncidentCategoryRequest,
) (*incidentclassifierv1.MoveIncidentCategoryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Move(ctx, classifiersvc.MoveIncidentCategoryCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.MoveIncidentCategoryPayload{
			CategoryID:          req.GetCategoryId(),
			NewParentCategoryID: req.NewParentCategoryId,
		},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.MoveIncidentCategoryResponse{}, nil
}
