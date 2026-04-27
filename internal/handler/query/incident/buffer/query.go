// Package buffer is the gRPC query-side transport for buffer reads.
package buffer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/model"
	bufferread "github.com/medincident/medincident-backend/internal/service/query/incident/buffer"
	incidentqueryv1 "github.com/medincident/medincident-backend/pkg/query/incident/v1"
)

const (
	errCodeHandlerInvalidBufferID     = "handler_invalid_buffer_id"
	errCodeHandlerInvalidOrgID        = "handler_invalid_organization_id"
	errCodeHandlerInvalidBufferStatus = "handler_invalid_buffer_status"
)

// BufferQueryHandler implements the buffer read methods for IncidentQueryService.
type BufferQueryHandler struct {
	reader *bufferread.Reader
}

// NewBufferQueryHandler wires the handler with a reader.
func NewBufferQueryHandler(r *bufferread.Reader) *BufferQueryHandler {
	return &BufferQueryHandler{reader: r}
}

// GetBufferEntry implements the GetBufferEntry RPC of IncidentQueryService.
func (h *BufferQueryHandler) GetBufferEntry(
	ctx context.Context, req *incidentqueryv1.GetBufferEntryRequest,
) (*incidentqueryv1.GetBufferEntryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(req.GetId(), "id", errCodeHandlerInvalidBufferID)
	if err != nil {
		return nil, err
	}
	v, err := h.reader.GetBufferEntry(ctx, callerID, id)
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.GetBufferEntryResponse{Entry: bufferToProto(v)}, nil
}

func (h *BufferQueryHandler) ListBufferEntries(
	ctx context.Context, req *incidentqueryv1.ListBufferEntriesRequest,
) (*incidentqueryv1.ListBufferEntriesResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseUUID(req.GetOrganizationId(), "organization_id", errCodeHandlerInvalidOrgID)
	if err != nil {
		return nil, err
	}
	statuses := make([]model.BufferStatus, 0, len(req.GetStatuses()))
	for _, s := range req.GetStatuses() {
		ms, err := protoBufStatusToModel(s)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, ms)
	}
	items, err := h.reader.ListBufferEntries(ctx, callerID, orgID, &bufferread.ListBufferFilters{
		Statuses: statuses,
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.ListBufferEntriesResponse{Items: buffersToProto(items)}, nil
}

func (h *BufferQueryHandler) ListMyBufferEntries(
	ctx context.Context, req *incidentqueryv1.ListMyBufferEntriesRequest,
) (*incidentqueryv1.ListMyBufferEntriesResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	items, err := h.reader.ListMyBufferEntries(ctx, callerID, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, err
	}
	return &incidentqueryv1.ListMyBufferEntriesResponse{Items: buffersToProto(items)}, nil
}

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.incident.buffer").
			Code(code).
			Public(field+" is not a valid UUID.").
			With(field, raw).Wrap(err)
	}
	return id, nil
}

func protoBufStatusToModel(s incidentqueryv1.BufferStatus) (model.BufferStatus, error) {
	switch s {
	case incidentqueryv1.BufferStatus_BUFFER_STATUS_PENDING:
		return model.BufferStatusPending, nil
	case incidentqueryv1.BufferStatus_BUFFER_STATUS_PUBLISHED:
		return model.BufferStatusPublished, nil
	case incidentqueryv1.BufferStatus_BUFFER_STATUS_REJECTED:
		return model.BufferStatusRejected, nil
	case incidentqueryv1.BufferStatus_BUFFER_STATUS_CANCELLED:
		return model.BufferStatusCancelled, nil
	default:
		return "", oops.In("handler.query.incident.buffer").
			Code(errCodeHandlerInvalidBufferStatus).
			Public("Invalid buffer status.").Errorf("unknown")
	}
}

func modelBufStatusToProto(s model.BufferStatus) incidentqueryv1.BufferStatus {
	switch s {
	case model.BufferStatusPending:
		return incidentqueryv1.BufferStatus_BUFFER_STATUS_PENDING
	case model.BufferStatusPublished:
		return incidentqueryv1.BufferStatus_BUFFER_STATUS_PUBLISHED
	case model.BufferStatusRejected:
		return incidentqueryv1.BufferStatus_BUFFER_STATUS_REJECTED
	case model.BufferStatusCancelled:
		return incidentqueryv1.BufferStatus_BUFFER_STATUS_CANCELLED
	default:
		return incidentqueryv1.BufferStatus_BUFFER_STATUS_UNSPECIFIED
	}
}

// patientStatusForBuffer maps internal buffer state to the patient
// vocabulary used when the caller is treated as a patient.
func patientStatusForBuffer(s model.BufferStatus) incidentqueryv1.PatientStatus {
	switch s {
	case model.BufferStatusPending:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_PENDING
	case model.BufferStatusPublished:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_ACCEPTED
	case model.BufferStatusRejected:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_CLOSED
	case model.BufferStatusCancelled:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_CANCELLED
	default:
		return incidentqueryv1.PatientStatus_PATIENT_STATUS_UNSPECIFIED
	}
}

func bufferToProto(v *bufferread.BufferEntryView) *incidentqueryv1.BufferEntryView {
	out := &incidentqueryv1.BufferEntryView{
		Id:                   v.ID.String(),
		OrganizationId:       v.OrganizationID.String(),
		PatientZitadelUserId: v.PatientZitadelUserID,
		Status:               modelBufStatusToProto(v.Status),
		CreatedAt:            v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:            v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if v.CategoryID.Valid {
		s := v.CategoryID.UUID.String()
		out.CategoryId = &s
	}
	if v.TypeID.Valid {
		s := v.TypeID.UUID.String()
		out.TypeId = &s
	}
	if v.Description.Valid {
		s := v.Description.String
		out.Description = &s
	}
	if v.OccurredAt.Valid {
		s := v.OccurredAt.Time.UTC().Format(time.RFC3339Nano)
		out.OccurredAt = &s
	}
	if v.PublishedIncidentID.Valid {
		s := v.PublishedIncidentID.UUID.String()
		out.PublishedIncidentId = &s
	}
	if v.PatientPerspective {
		ps := patientStatusForBuffer(v.Status)
		out.PatientStatus = &ps
	}
	return out
}

func buffersToProto(views []bufferread.BufferEntryView) []*incidentqueryv1.BufferEntryView {
	out := make([]*incidentqueryv1.BufferEntryView, 0, len(views))
	for i := range views {
		out = append(out, bufferToProto(&views[i]))
	}
	return out
}
