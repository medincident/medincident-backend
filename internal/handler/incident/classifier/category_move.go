package classifier

import (
	"context"

	incidentclassifierv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/incident/classifier/v1"
	classifiersvc "github.com/medincident/medincident-command-service/internal/services/incident/classifier"
)

func (h *IncidentClassifierHandler) MoveIncidentCategory(
	ctx context.Context,
	req *incidentclassifierv1.MoveIncidentCategoryRequest,
) (*incidentclassifierv1.MoveIncidentCategoryResponse, error) {
	categoryID, err := parseIncidentCategoryID(req.GetCategoryId())
	if err != nil {
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
