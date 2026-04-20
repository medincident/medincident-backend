package classifier

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	classifierread "github.com/medincident/medincident-command-service/internal/service/query/incident/classifier"
	classifierqueryv1 "github.com/medincident/medincident-command-service/pkg/query/incident/classifier/v1"
)

// Error codes emitted by query-handler ID parsing.
const (
	ErrCodeHandlerInvalidCategoryID     = "handler_invalid_category_id"
	ErrCodeHandlerInvalidTypeID         = "handler_invalid_type_id"
	ErrCodeHandlerInvalidOrganizationID = "handler_invalid_organization_id"
)

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.incident.classifier").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("Invalid organization id.").
			With("organization_id", raw).
			Wrap(err)
	}
	return id, nil
}

// IncidentClassifierQueryHandler implements
// classifierqueryv1.IncidentClassifierQueryServiceServer.
type IncidentClassifierQueryHandler struct {
	classifierqueryv1.UnimplementedIncidentClassifierQueryServiceServer

	reader *classifierread.Reader
}

// NewIncidentClassifierQueryHandler wires the handler with a reader.
func NewIncidentClassifierQueryHandler(reader *classifierread.Reader) *IncidentClassifierQueryHandler {
	return &IncidentClassifierQueryHandler{reader: reader}
}

// GetCategory returns one category.
func (h *IncidentClassifierQueryHandler) GetCategory(
	ctx context.Context,
	req *classifierqueryv1.GetCategoryRequest,
) (*classifierqueryv1.GetCategoryResponse, error) {
	id, err := parseCategoryID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.GetCategoryResponse{Category: categoryToProto(view)}, nil
}

// ListCategoriesByOrganization paginates categories for one org.
func (h *IncidentClassifierQueryHandler) ListCategoriesByOrganization(
	ctx context.Context,
	req *classifierqueryv1.ListCategoriesByOrganizationRequest,
) (*classifierqueryv1.ListCategoriesByOrganizationResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListCategoriesByOrganization(ctx, id, classifierread.ListQuery{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListCategoriesByOrganizationResponse{Items: categoriesToProto(items)}, nil
}

// ListActiveRootCategories returns active roots.
func (h *IncidentClassifierQueryHandler) ListActiveRootCategories(
	ctx context.Context,
	req *classifierqueryv1.ListActiveRootCategoriesRequest,
) (*classifierqueryv1.ListActiveRootCategoriesResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListActiveRootCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListActiveRootCategoriesResponse{Items: categoriesToProto(items)}, nil
}

// ListCategorySubtree returns a flattened subtree.
func (h *IncidentClassifierQueryHandler) ListCategorySubtree(
	ctx context.Context,
	req *classifierqueryv1.ListCategorySubtreeRequest,
) (*classifierqueryv1.ListCategorySubtreeResponse, error) {
	id, err := parseCategoryID(req.GetRootCategoryId())
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListCategorySubtree(ctx, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListCategorySubtreeResponse{Items: categoriesToProto(items)}, nil
}

// GetType returns one type.
func (h *IncidentClassifierQueryHandler) GetType(
	ctx context.Context,
	req *classifierqueryv1.GetTypeRequest,
) (*classifierqueryv1.GetTypeResponse, error) {
	id, err := parseTypeID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetType(ctx, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.GetTypeResponse{Type: typeToProto(view)}, nil
}

// ListTypesByCategory returns every type under one category.
func (h *IncidentClassifierQueryHandler) ListTypesByCategory(
	ctx context.Context,
	req *classifierqueryv1.ListTypesByCategoryRequest,
) (*classifierqueryv1.ListTypesByCategoryResponse, error) {
	id, err := parseCategoryID(req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListTypesByCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListTypesByCategoryResponse{Items: typesToProto(items)}, nil
}

// ListActiveTypesByOrganization returns active types for an org.
func (h *IncidentClassifierQueryHandler) ListActiveTypesByOrganization(
	ctx context.Context,
	req *classifierqueryv1.ListActiveTypesByOrganizationRequest,
) (*classifierqueryv1.ListActiveTypesByOrganizationResponse, error) {
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListActiveTypesByOrganization(ctx, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListActiveTypesByOrganizationResponse{Items: typesToProto(items)}, nil
}

// parseCategoryID parses a category UUID.
func parseCategoryID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.incident.classifier.query").
			Code(ErrCodeHandlerInvalidCategoryID).
			Public("category_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

// parseTypeID parses a type UUID.
func parseTypeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.incident.classifier.query").
			Code(ErrCodeHandlerInvalidTypeID).
			Public("type_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

// categoryToProto adapts one CategoryView.
func categoryToProto(v *classifierread.CategoryView) *classifierqueryv1.Category {
	out := &classifierqueryv1.Category{
		Id:             v.ID.String(),
		OrganizationId: v.OrganizationID.String(),
		Name:           v.Name,
		Description:    v.Description,
		IsActive:       v.IsActive,
		CreatedAt:      v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if v.ParentCategoryID != nil {
		s := v.ParentCategoryID.String()
		out.ParentCategoryId = &s
	}
	return out
}

// categoriesToProto adapts a list of CategoryView.
func categoriesToProto(views []classifierread.CategoryView) []*classifierqueryv1.Category {
	out := make([]*classifierqueryv1.Category, 0, len(views))
	for i := range views {
		out = append(out, categoryToProto(&views[i]))
	}
	return out
}

// typeToProto adapts one TypeView.
func typeToProto(v *classifierread.TypeView) *classifierqueryv1.Type {
	return &classifierqueryv1.Type{
		Id:             v.ID.String(),
		OrganizationId: v.OrganizationID.String(),
		CategoryId:     v.CategoryID.String(),
		Name:           v.Name,
		Description:    v.Description,
		IsActive:       v.IsActive,
		CreatedAt:      v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

// typesToProto adapts a list of TypeView.
func typesToProto(views []classifierread.TypeView) []*classifierqueryv1.Type {
	out := make([]*classifierqueryv1.Type, 0, len(views))
	for i := range views {
		out = append(out, typeToProto(&views[i]))
	}
	return out
}
