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
	cmd := classifiersvc.MoveIncidentCategoryCommand{CategoryID: categoryID}
	if req.NewParentCategoryId != nil {
		newParentID, err := parseIncidentCategoryID(*req.NewParentCategoryId)
		if err != nil {
			return nil, err
		}
		cmd.NewParentCategoryID = &newParentID
	}
	if _, err := h.categorySvc.Move(ctx, cmd); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.MoveIncidentCategoryResponse{}, nil
}
