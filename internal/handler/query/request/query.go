// Package request is the gRPC query-side transport for service requests.
package request

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	queryrequest "github.com/medincident/medincident-backend/internal/service/query/request"
	requestqueryv1 "github.com/medincident/medincident-backend/pkg/query/request/v1"
)

const (
	errCodeHandlerInvalidServiceRequestID = "handler_invalid_service_request_id"
	errCodeHandlerInvalidOrgID            = "handler_invalid_organization_id"
	errCodeHandlerInvalidIncidentID       = "handler_invalid_incident_id"
)

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.request").
			Code(code).
			Public(field+" is not a valid UUID.").
			With(field, raw).Wrap(err)
	}
	return id, nil
}

// ServiceRequestQueryHandler implements the service request read methods.
type ServiceRequestQueryHandler struct {
	requestqueryv1.UnimplementedServiceRequestQueryServiceServer
	reader *queryrequest.Reader
}

// NewServiceRequestQueryHandler wires the handler with a reader.
func NewServiceRequestQueryHandler(r *queryrequest.Reader) *ServiceRequestQueryHandler {
	return &ServiceRequestQueryHandler{reader: r}
}

func (h *ServiceRequestQueryHandler) GetServiceRequest(
	ctx context.Context, req *requestqueryv1.GetServiceRequestRequest,
) (*requestqueryv1.GetServiceRequestResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(req.GetId(), "id", errCodeHandlerInvalidServiceRequestID)
	if err != nil {
		return nil, err
	}
	v, err := h.reader.GetServiceRequest(ctx, authz.Caller{ZitadelUserID: callerID}, id)
	if err != nil {
		return nil, err
	}
	return &requestqueryv1.GetServiceRequestResponse{ServiceRequest: serviceRequestToProto(v)}, nil
}

func (h *ServiceRequestQueryHandler) ListServiceRequests(
	ctx context.Context, req *requestqueryv1.ListServiceRequestsRequest,
) (*requestqueryv1.ListServiceRequestsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseUUID(req.GetOrganizationId(), "organization_id", errCodeHandlerInvalidOrgID)
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListServiceRequests(ctx, authz.Caller{ZitadelUserID: callerID}, orgID, queryrequest.ListQuery{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	return &requestqueryv1.ListServiceRequestsResponse{Items: serviceRequestsToProto(items)}, nil
}

func (h *ServiceRequestQueryHandler) ListServiceRequestsByIncident(
	ctx context.Context, req *requestqueryv1.ListServiceRequestsByIncidentRequest,
) (*requestqueryv1.ListServiceRequestsByIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	incidentID, err := parseUUID(req.GetIncidentId(), "incident_id", errCodeHandlerInvalidIncidentID)
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListServiceRequestsByIncident(ctx, authz.Caller{ZitadelUserID: callerID}, incidentID, queryrequest.ListQuery{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	return &requestqueryv1.ListServiceRequestsByIncidentResponse{Items: serviceRequestsToProto(items)}, nil
}

func (h *ServiceRequestQueryHandler) GetServiceRequestHistory(
	ctx context.Context, req *requestqueryv1.GetServiceRequestHistoryRequest,
) (*requestqueryv1.GetServiceRequestHistoryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(req.GetServiceRequestId(), "service_request_id", errCodeHandlerInvalidServiceRequestID)
	if err != nil {
		return nil, err
	}
	hist, err := h.reader.GetServiceRequestHistory(ctx, authz.Caller{ZitadelUserID: callerID}, id)
	if err != nil {
		return nil, err
	}
	return &requestqueryv1.GetServiceRequestHistoryResponse{
		StatusHistory:   statusHistoryToProto(hist.StatusHistory),
		ExecutorHistory: executorHistoryToProto(hist.ExecutorHistory),
	}, nil
}

func serviceRequestToProto(v *queryrequest.ServiceRequestView) *requestqueryv1.ServiceRequest {
	out := &requestqueryv1.ServiceRequest{
		Id:                v.ID.String(),
		OrganizationId:    v.OrganizationID.String(),
		ClinicId:          v.ClinicID.String(),
		DepartmentId:      v.DepartmentID.String(),
		TypeId:            v.TypeID.String(),
		Description:       v.Description,
		Status:            string(v.Status),
		AuthorId:          v.AuthorID,
		AuthorDisplayName: v.AuthorDisplayName,
		CreatedAt:         v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:         v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if v.IncidentID != nil {
		s := v.IncidentID.String()
		out.IncidentId = &s
	}
	for _, e := range v.Executors {
		out.Executors = append(out.Executors, &requestqueryv1.Executor{
			EmployeeId:   e.EmployeeID.String(),
			AssignedAt:   e.AssignedAt.UTC().Format(time.RFC3339Nano),
			AssignedById: e.AssignedByID,
		})
	}
	return out
}

func serviceRequestsToProto(views []queryrequest.ServiceRequestView) []*requestqueryv1.ServiceRequest {
	out := make([]*requestqueryv1.ServiceRequest, 0, len(views))
	for i := range views {
		out = append(out, serviceRequestToProto(&views[i]))
	}
	return out
}

func statusHistoryToProto(items []queryrequest.StatusHistoryView) []*requestqueryv1.StatusHistoryEntry {
	out := make([]*requestqueryv1.StatusHistoryEntry, 0, len(items))
	for i := range items {
		e := &items[i]
		entry := &requestqueryv1.StatusHistoryEntry{
			Id:        e.ID.String(),
			NewStatus: e.NewStatus,
			ActorId:   e.ActorID,
			ActorName: e.ActorName,
			ChangedAt: e.ChangedAt.UTC().Format(time.RFC3339Nano),
		}
		if e.OldStatus != nil {
			entry.OldStatus = e.OldStatus
		}
		out = append(out, entry)
	}
	return out
}

func executorHistoryToProto(items []queryrequest.ExecutorHistoryView) []*requestqueryv1.ExecutorHistoryEntry {
	out := make([]*requestqueryv1.ExecutorHistoryEntry, 0, len(items))
	for i := range items {
		e := &items[i]
		out = append(out, &requestqueryv1.ExecutorHistoryEntry{
			Id:           e.ID.String(),
			Action:       e.Action,
			EmployeeId:   e.EmployeeID.String(),
			EmployeeName: e.EmployeeName,
			ActorId:      e.ActorID,
			ActorName:    e.ActorName,
			ChangedAt:    e.ChangedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	return out
}
