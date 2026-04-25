// Package incident is the gRPC query-side transport.
package incident

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/model"
	queryincident "github.com/medincident/medincident-backend/internal/service/query/incident"
	incidentqueryv1 "github.com/medincident/medincident-backend/pkg/query/incident/v1"
)

const (
	errCodeHandlerInvalidIncidentID = "handler_invalid_incident_id"
	errCodeHandlerInvalidOrgID      = "handler_invalid_organization_id"
	errCodeHandlerInvalidUUID       = "handler_invalid_uuid"
	errCodeHandlerInvalidTimestamp  = "handler_invalid_timestamp"
)

// IncidentQueryHandler implements the incident read methods.
type IncidentQueryHandler struct {
	incidentqueryv1.UnimplementedIncidentQueryServiceServer
	reader *queryincident.Reader
}

// NewIncidentQueryHandler wires the handler with a reader.
func NewIncidentQueryHandler(r *queryincident.Reader) *IncidentQueryHandler {
	return &IncidentQueryHandler{reader: r}
}

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.incident").
			Code(code).
			Public(field+" is not a valid UUID.").
			With(field, raw).Wrap(err)
	}
	return id, nil
}

func (h *IncidentQueryHandler) GetIncident(
	ctx context.Context, req *incidentqueryv1.GetIncidentRequest,
) (*incidentqueryv1.GetIncidentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(req.GetId(), "id", errCodeHandlerInvalidIncidentID)
	if err != nil {
		return nil, err
	}
	v, err := h.reader.GetIncident(ctx, callerID, id)
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.GetIncidentResponse{Incident: incidentToProto(v)}, nil
}

func (h *IncidentQueryHandler) ListIncidents(
	ctx context.Context, req *incidentqueryv1.ListIncidentsRequest,
) (*incidentqueryv1.ListIncidentsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseUUID(req.GetOrganizationId(), "organization_id", errCodeHandlerInvalidOrgID)
	if err != nil {
		return nil, err
	}
	f, err := decodeListFilters(req)
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListIncidents(ctx, callerID, orgID, &f)
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.ListIncidentsResponse{Items: incidentsToProto(items)}, nil
}

func (h *IncidentQueryHandler) ListMyIncidents(
	ctx context.Context, req *incidentqueryv1.ListMyIncidentsRequest,
) (*incidentqueryv1.ListMyIncidentsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListMyIncidents(ctx, callerID, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.ListMyIncidentsResponse{Items: incidentsToProto(items)}, nil
}

func (h *IncidentQueryHandler) GetIncidentHistory(
	ctx context.Context, req *incidentqueryv1.GetIncidentHistoryRequest,
) (*incidentqueryv1.GetIncidentHistoryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(req.GetIncidentId(), "incident_id", errCodeHandlerInvalidIncidentID)
	if err != nil {
		return nil, err
	}
	hist, err := h.reader.GetIncidentHistory(ctx, callerID, id)
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.GetIncidentHistoryResponse{
		StatusHistory:   statusHistoryToProto(hist.Status),
		PriorityHistory: priorityHistoryToProto(hist.Priority),
	}, nil
}

func decodeListFilters(req *incidentqueryv1.ListIncidentsRequest) (queryincident.ListFilters, error) {
	f := queryincident.ListFilters{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
	for _, s := range req.GetStatuses() {
		ms, err := protoStatusToModel(s)
		if err != nil {
			return f, err
		}
		f.Statuses = append(f.Statuses, ms)
	}
	for _, p := range req.GetPriorities() {
		mp, err := protoPriorityToModel(p)
		if err != nil {
			return f, err
		}
		f.Priorities = append(f.Priorities, mp)
	}
	if req.ClinicId != nil {
		id, err := parseUUID(*req.ClinicId, "clinic_id", errCodeHandlerInvalidUUID)
		if err != nil {
			return f, err
		}
		f.ClinicID = uuid.NullUUID{UUID: id, Valid: true}
	}
	if req.DepartmentId != nil {
		id, err := parseUUID(*req.DepartmentId, "department_id", errCodeHandlerInvalidUUID)
		if err != nil {
			return f, err
		}
		f.DepartmentID = uuid.NullUUID{UUID: id, Valid: true}
	}
	if req.CategoryId != nil {
		id, err := parseUUID(*req.CategoryId, "category_id", errCodeHandlerInvalidUUID)
		if err != nil {
			return f, err
		}
		f.CategoryID = uuid.NullUUID{UUID: id, Valid: true}
	}
	if req.TypeId != nil {
		id, err := parseUUID(*req.TypeId, "type_id", errCodeHandlerInvalidUUID)
		if err != nil {
			return f, err
		}
		f.TypeID = uuid.NullUUID{UUID: id, Valid: true}
	}
	if req.OccurredFrom != nil {
		t, err := time.Parse(time.RFC3339Nano, *req.OccurredFrom)
		if err != nil {
			return f, oops.In("handler.query.incident").
				Code(errCodeHandlerInvalidTimestamp).
				Public("occurred_from is invalid.").Wrap(err)
		}
		f.OccurredFrom = &t
	}
	if req.OccurredTo != nil {
		t, err := time.Parse(time.RFC3339Nano, *req.OccurredTo)
		if err != nil {
			return f, oops.In("handler.query.incident").
				Code(errCodeHandlerInvalidTimestamp).
				Public("occurred_to is invalid.").Wrap(err)
		}
		f.OccurredTo = &t
	}
	return f, nil
}

func protoStatusToModel(s incidentqueryv1.IncidentStatus) (model.IncidentStatus, error) {
	switch s {
	case incidentqueryv1.IncidentStatus_INCIDENT_STATUS_PENDING:
		return model.IncidentStatusPending, nil
	case incidentqueryv1.IncidentStatus_INCIDENT_STATUS_IN_PROGRESS:
		return model.IncidentStatusInProgress, nil
	case incidentqueryv1.IncidentStatus_INCIDENT_STATUS_DONE:
		return model.IncidentStatusDone, nil
	case incidentqueryv1.IncidentStatus_INCIDENT_STATUS_REJECTED:
		return model.IncidentStatusRejected, nil
	case incidentqueryv1.IncidentStatus_INCIDENT_STATUS_CANCELLED:
		return model.IncidentStatusCancelled, nil
	default:
		return "", oops.In("handler.query.incident").
			Code(errCodeHandlerInvalidUUID).Public("Invalid status.").Errorf("unknown")
	}
}

func modelStatusToProto(s model.IncidentStatus) incidentqueryv1.IncidentStatus {
	switch s {
	case model.IncidentStatusPending:
		return incidentqueryv1.IncidentStatus_INCIDENT_STATUS_PENDING
	case model.IncidentStatusInProgress:
		return incidentqueryv1.IncidentStatus_INCIDENT_STATUS_IN_PROGRESS
	case model.IncidentStatusDone:
		return incidentqueryv1.IncidentStatus_INCIDENT_STATUS_DONE
	case model.IncidentStatusRejected:
		return incidentqueryv1.IncidentStatus_INCIDENT_STATUS_REJECTED
	case model.IncidentStatusCancelled:
		return incidentqueryv1.IncidentStatus_INCIDENT_STATUS_CANCELLED
	default:
		return incidentqueryv1.IncidentStatus_INCIDENT_STATUS_UNSPECIFIED
	}
}

func protoPriorityToModel(p incidentqueryv1.IncidentPriority) (model.IncidentPriority, error) {
	switch p {
	case incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_LOW:
		return model.IncidentPriorityLow, nil
	case incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_NORMAL:
		return model.IncidentPriorityNormal, nil
	case incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_HIGH:
		return model.IncidentPriorityHigh, nil
	case incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_CRITICAL:
		return model.IncidentPriorityCritical, nil
	default:
		return "", oops.In("handler.query.incident").
			Code(errCodeHandlerInvalidUUID).Public("Invalid priority.").Errorf("unknown")
	}
}

func modelPriorityToProto(p model.IncidentPriority) incidentqueryv1.IncidentPriority {
	switch p {
	case model.IncidentPriorityLow:
		return incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_LOW
	case model.IncidentPriorityNormal:
		return incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_NORMAL
	case model.IncidentPriorityHigh:
		return incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_HIGH
	case model.IncidentPriorityCritical:
		return incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_CRITICAL
	default:
		return incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_UNSPECIFIED
	}
}

// patientStatusForIncident maps an internal status to the patient
// four-value vocabulary. Only used when PatientPerspective is true.
func patientStatusForIncident(s model.IncidentStatus) incidentqueryv1.PatientStatus {
	switch s {
	case model.IncidentStatusPending, model.IncidentStatusInProgress:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_ACCEPTED
	case model.IncidentStatusDone, model.IncidentStatusRejected:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_CLOSED
	case model.IncidentStatusCancelled:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_CANCELLED
	default:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_UNSPECIFIED
	}
}

func incidentToProto(v *queryincident.IncidentView) *incidentqueryv1.IncidentView {
	out := &incidentqueryv1.IncidentView{
		Id:         v.ID.String(),
		Status:     modelStatusToProto(v.Status),
		Priority:   modelPriorityToProto(v.Priority),
		OccurredAt: v.OccurredAt.UTC().Format(time.RFC3339Nano),
		CreatedAt:  v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:  v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if v.Description.Valid {
		s := v.Description.String
		out.Description = &s
	}

	if v.PatientPerspective {
		ps := patientStatusForIncident(v.Status)
		out.PatientStatus = &ps
		out.Priority = incidentqueryv1.IncidentPriority_INCIDENT_PRIORITY_UNSPECIFIED
		return out
	}

	orgID := v.OrganizationID.String()
	clinicID := v.ClinicID.String()
	deptID := v.DepartmentID.String()
	catID := v.CategoryID.String()
	typID := v.TypeID.String()
	out.OrganizationId = &orgID
	out.ClinicId = &clinicID
	out.DepartmentId = &deptID
	out.CategoryId = &catID
	out.TypeId = &typID
	if v.PatientOriginalDescription.Valid {
		s := v.PatientOriginalDescription.String
		out.PatientOriginalDescription = &s
	}
	out.Registrar = &incidentqueryv1.RegistrarView{
		EmployeeId:     v.RegistrarEmployeeID.String(),
		DisplayName:    v.RegistrarDisplayName,
		OrganizationId: v.RegistrarOrganizationID.String(),
		ClinicId:       v.RegistrarClinicID.String(),
		DepartmentId:   v.RegistrarDepartmentID.String(),
	}
	if v.RegistrarPosition.Valid {
		s := v.RegistrarPosition.String
		out.Registrar.Position = &s
	}
	if v.SourcePatientZitadelUserID.Valid {
		s := v.SourcePatientZitadelUserID.String
		out.SourcePatientZitadelUserId = &s
	}
	if v.SourceBufferID.Valid {
		s := v.SourceBufferID.UUID.String()
		out.SourceBufferId = &s
	}
	if v.ReopenedFromIncidentID.Valid {
		s := v.ReopenedFromIncidentID.UUID.String()
		out.ReopenedFromIncidentId = &s
	}
	return out
}

func incidentsToProto(views []queryincident.IncidentView) []*incidentqueryv1.IncidentView {
	out := make([]*incidentqueryv1.IncidentView, 0, len(views))
	for i := range views {
		out = append(out, incidentToProto(&views[i]))
	}
	return out
}

func actorToProto(empID uuid.NullUUID, name null.String) *incidentqueryv1.ActorView {
	v := &incidentqueryv1.ActorView{}
	if empID.Valid {
		s := empID.UUID.String()
		v.EmployeeId = &s
	}
	if name.Valid {
		s := name.String
		v.DisplayName = &s
	}
	return v
}

func statusHistoryToProto(items []queryincident.StatusHistoryEntry) []*incidentqueryv1.StatusHistoryEntry {
	out := make([]*incidentqueryv1.StatusHistoryEntry, 0, len(items))
	for i := range items {
		e := &items[i]
		old := incidentqueryv1.IncidentStatus_INCIDENT_STATUS_UNSPECIFIED
		if e.OldStatus.Valid {
			old = modelStatusToProto(model.IncidentStatus(e.OldStatus.String))
		}
		out = append(out, &incidentqueryv1.StatusHistoryEntry{
			Id:        e.ID.String(),
			OldStatus: old,
			NewStatus: modelStatusToProto(e.NewStatus),
			Actor:     actorToProto(e.ActorEmployeeID, e.ActorDisplayName),
			ChangedAt: e.ChangedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	return out
}

func priorityHistoryToProto(items []queryincident.PriorityHistoryEntry) []*incidentqueryv1.PriorityHistoryEntry {
	out := make([]*incidentqueryv1.PriorityHistoryEntry, 0, len(items))
	for _, e := range items {
		out = append(out, &incidentqueryv1.PriorityHistoryEntry{
			Id:          e.ID.String(),
			OldPriority: modelPriorityToProto(e.OldPriority),
			NewPriority: modelPriorityToProto(e.NewPriority),
			Actor:       actorToProto(e.ActorEmployeeID, e.ActorDisplayName),
			ChangedAt:   e.ChangedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	return out
}

// CombinedIncidentQueryHandler implements the full IncidentQueryService
// surface by delegating buffer methods to a separate handler. It exists
// because the proto bundles incident + buffer reads under one service.
type CombinedIncidentQueryHandler struct {
	*IncidentQueryHandler
	bufferDelegate IncidentBufferQueryDelegate
}

// IncidentBufferQueryDelegate is implemented by buffer.BufferQueryHandler.
type IncidentBufferQueryDelegate interface {
	GetBufferEntry(ctx context.Context, req *incidentqueryv1.GetBufferEntryRequest) (*incidentqueryv1.GetBufferEntryResponse, error)
	ListBufferEntries(ctx context.Context, req *incidentqueryv1.ListBufferEntriesRequest) (*incidentqueryv1.ListBufferEntriesResponse, error)
	ListMyBufferEntries(ctx context.Context, req *incidentqueryv1.ListMyBufferEntriesRequest) (*incidentqueryv1.ListMyBufferEntriesResponse, error)
}

// NewCombinedIncidentQueryHandler wires both readers behind a single gRPC server.
func NewCombinedIncidentQueryHandler(
	inc *IncidentQueryHandler,
	buf IncidentBufferQueryDelegate,
) *CombinedIncidentQueryHandler {
	return &CombinedIncidentQueryHandler{IncidentQueryHandler: inc, bufferDelegate: buf}
}

func (h *CombinedIncidentQueryHandler) GetBufferEntry(
	ctx context.Context, req *incidentqueryv1.GetBufferEntryRequest,
) (*incidentqueryv1.GetBufferEntryResponse, error) {
	return h.bufferDelegate.GetBufferEntry(ctx, req)
}

func (h *CombinedIncidentQueryHandler) ListBufferEntries(
	ctx context.Context, req *incidentqueryv1.ListBufferEntriesRequest,
) (*incidentqueryv1.ListBufferEntriesResponse, error) {
	return h.bufferDelegate.ListBufferEntries(ctx, req)
}

func (h *CombinedIncidentQueryHandler) ListMyBufferEntries(
	ctx context.Context, req *incidentqueryv1.ListMyBufferEntriesRequest,
) (*incidentqueryv1.ListMyBufferEntriesResponse, error) {
	return h.bufferDelegate.ListMyBufferEntries(ctx, req)
}
