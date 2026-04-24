package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) DeactivateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.DeactivateIncidentCategoryRequest,
) (*incidentclassifierv1.DeactivateIncidentCategoryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Deactivate(ctx, classifiersvc.DeactivateIncidentCategoryCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.DeactivateIncidentCategoryPayload{CategoryID: req.GetCategoryId()},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeactivateIncidentCategoryResponse{}, nil
}
