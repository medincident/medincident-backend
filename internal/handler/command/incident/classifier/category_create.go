package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) CreateIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentCategoryRequest,
) (*incidentclassifierv1.CreateIncidentCategoryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := h.categorySvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID:   req.GetOrganizationId(),
			ParentCategoryID: req.ParentCategoryId,
			Name:             req.GetName(),
			Description:      req.Description,
		},
	})
	if err != nil {
		return nil, err
	}
	return &incidentclassifierv1.CreateIncidentCategoryResponse{
		CategoryId: result.ID.String(),
	}, nil
}
