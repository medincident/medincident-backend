package classifier

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	classifierread "github.com/medincident/medincident-backend/internal/service/query/request/classifier"
	classifierqueryv1 "github.com/medincident/medincident-backend/pkg/query/request/classifier/v1"
)

// Error codes emitted by query-handler ID parsing.
const (
	ErrCodeHandlerInvalidRequestTypeID  = "handler_invalid_request_type_id"
	ErrCodeHandlerInvalidOrganizationID = "handler_invalid_organization_id"
)

func parseRequestTypeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.request.classifier").
			Code(ErrCodeHandlerInvalidRequestTypeID).
			Public("id is not a valid UUID.").
			With("id", raw).
			Wrap(err)
	}
	return id, nil
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.request.classifier").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("Invalid organization id.").
			With("organization_id", raw).
			Wrap(err)
	}
	return id, nil
}

// afterPtr converts an empty proto string to nil, treating empty as
// "start from the beginning".
func afterPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// RequestClassifierQueryHandler implements
// classifierqueryv1.RequestClassifierQueryServiceServer.
type RequestClassifierQueryHandler struct {
	classifierqueryv1.UnimplementedRequestClassifierQueryServiceServer

	reader *classifierread.Reader
}

// NewRequestClassifierQueryHandler wires the handler with a reader.
func NewRequestClassifierQueryHandler(reader *classifierread.Reader) *RequestClassifierQueryHandler {
	return &RequestClassifierQueryHandler{reader: reader}
}

// GetRequestType returns one request type.
func (h *RequestClassifierQueryHandler) GetRequestType(
	ctx context.Context,
	req *classifierqueryv1.GetRequestTypeRequest,
) (*classifierqueryv1.GetRequestTypeResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseRequestTypeID(req.GetId())
	if err != nil {
		return nil, err
	}
	view, err := h.reader.GetRequestType(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.GetRequestTypeResponse{RequestType: requestTypeToProto(view)}, nil
}

// ListRequestTypesByOrganization paginates request types for one org.
func (h *RequestClassifierQueryHandler) ListRequestTypesByOrganization(
	ctx context.Context,
	req *classifierqueryv1.ListRequestTypesByOrganizationRequest,
) (*classifierqueryv1.ListRequestTypesByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.reader.ListRequestTypesByOrganization(ctx, caller, id, classifierread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListRequestTypesByOrganizationResponse{
		Items:      requestTypesToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// ListActiveRequestTypesByOrganization returns active request types
// for an org.
func (h *RequestClassifierQueryHandler) ListActiveRequestTypesByOrganization(
	ctx context.Context,
	req *classifierqueryv1.ListActiveRequestTypesByOrganizationRequest,
) (*classifierqueryv1.ListActiveRequestTypesByOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	result, err := h.reader.ListActiveRequestTypesByOrganization(ctx, caller, id, classifierread.ListQuery{
		Limit: int(req.GetLimit()),
		After: afterPtr(req.GetAfter()),
	})
	if err != nil {
		return nil, err
	}
	return &classifierqueryv1.ListActiveRequestTypesByOrganizationResponse{
		Items:      requestTypesToProto(result.Items),
		NextCursor: result.NextCursor,
	}, nil
}

// requestTypeToProto adapts one RequestTypeView.
func requestTypeToProto(v *classifierread.RequestTypeView) *classifierqueryv1.RequestType {
	out := &classifierqueryv1.RequestType{
		Id:             v.ID.String(),
		OrganizationId: v.OrganizationID.String(),
		Name:           v.Name,
		IsActive:       v.IsActive,
		CreatedAt:      v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if v.Description != nil {
		out.Description = v.Description
	}
	return out
}

// requestTypesToProto adapts a list of RequestTypeView.
func requestTypesToProto(views []classifierread.RequestTypeView) []*classifierqueryv1.RequestType {
	out := make([]*classifierqueryv1.RequestType, 0, len(views))
	for i := range views {
		out = append(out, requestTypeToProto(&views[i]))
	}
	return out
}
