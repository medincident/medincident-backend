package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) DeleteIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.DeleteIncidentCategoryRequest,
) (*incidentclassifierv1.DeleteIncidentCategoryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.categorySvc.Delete(ctx, classifiersvc.DeleteIncidentCategoryCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.DeleteIncidentCategoryPayload{CategoryID: req.GetCategoryId()},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DeleteIncidentCategoryResponse{}, nil
}
