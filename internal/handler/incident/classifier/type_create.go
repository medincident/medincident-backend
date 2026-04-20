package classifier

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-backend/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) CreateIncidentType(
	ctx context.Context,
	req *incidentclassifierv1.CreateIncidentTypeRequest,
) (*incidentclassifierv1.CreateIncidentTypeResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.authz.Require(ctx, callerID, authz.AdminOf.Category(categoryID)); err != nil {
		return nil, err
	}
	result, err := h.typeSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		CategoryID:  categoryID,
		Name:        req.GetName(),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &incidentclassifierv1.CreateIncidentTypeResponse{TypeId: result.ID.String()}, nil
}
