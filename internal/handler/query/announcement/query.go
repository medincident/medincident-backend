// Package announcement is the gRPC query-side transport for AnnouncementQueryService.
package announcement

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	announcementread "github.com/medincident/medincident-backend/internal/service/query/announcement"
	announcementqueryv1 "github.com/medincident/medincident-backend/pkg/query/announcement/v1"
)

const handlerScope = "handler.query.announcement"

const (
	errCodeBadID      = "handler_invalid_announcement_id"
	errCodeBadScopeID = "handler_invalid_scope_id"
)

// AnnouncementQueryHandler implements AnnouncementQueryServiceServer.
type AnnouncementQueryHandler struct {
	announcementqueryv1.UnimplementedAnnouncementQueryServiceServer
	reader *announcementread.Reader
}

// NewAnnouncementQueryHandler wires the handler with a reader.
func NewAnnouncementQueryHandler(r *announcementread.Reader) *AnnouncementQueryHandler {
	return &AnnouncementQueryHandler{reader: r}
}

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In(handlerScope).
			Code(code).
			Public(field+" is not a valid UUID.").
			With(field, raw).Wrap(err)
	}
	return id, nil
}

func rowToProto(r *announcementread.AnnouncementRow) *announcementqueryv1.AnnouncementView {
	v := &announcementqueryv1.AnnouncementView{
		Id:             r.ID.String(),
		OrganizationId: r.OrganizationID.String(),
		AuthorId:       r.AuthorID,
		Title:          r.Title,
		Content:        r.Content,
		Priority:       modelToProtoPriority(r.Priority),
		IsArchived:     r.IsArchived,
		CreatedAt:      r.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      r.UpdatedAt.UTC().Format(time.RFC3339Nano),
		ViewCount:      r.ViewCount,
	}
	if r.ClinicID.Valid {
		s := r.ClinicID.UUID.String()
		v.ClinicId = &s
	}
	if r.DepartmentID.Valid {
		s := r.DepartmentID.UUID.String()
		v.DepartmentId = &s
	}
	v.StartsAt = announcementread.NullTimePtr(r.StartsAt)
	v.EndsAt = announcementread.NullTimePtr(r.EndsAt)
	return v
}

func modelToProtoPriority(p string) announcementqueryv1.AnnouncementPriority {
	switch p {
	case "high":
		return announcementqueryv1.AnnouncementPriority_ANNOUNCEMENT_PRIORITY_HIGH
	default:
		return announcementqueryv1.AnnouncementPriority_ANNOUNCEMENT_PRIORITY_NORMAL
	}
}

func protoToModelPriority(p announcementqueryv1.AnnouncementPriority) string {
	switch p {
	case announcementqueryv1.AnnouncementPriority_ANNOUNCEMENT_PRIORITY_HIGH:
		return "high"
	case announcementqueryv1.AnnouncementPriority_ANNOUNCEMENT_PRIORITY_NORMAL:
		return "normal"
	default:
		return ""
	}
}

func (h *AnnouncementQueryHandler) GetAnnouncement(
	ctx context.Context, req *announcementqueryv1.GetAnnouncementRequest,
) (*announcementqueryv1.GetAnnouncementResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(req.GetId(), "id", errCodeBadID)
	if err != nil {
		return nil, err
	}
	row, err := h.reader.GetAnnouncement(ctx, callerID, id)
	if err != nil {
		return nil, err
	}
	return &announcementqueryv1.GetAnnouncementResponse{Announcement: rowToProto(row)}, nil
}

func (h *AnnouncementQueryHandler) ListAnnouncementsForOrganization(
	ctx context.Context, req *announcementqueryv1.ListAnnouncementsForOrganizationRequest,
) (*announcementqueryv1.ListAnnouncementsForOrganizationResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	orgID, err := parseUUID(req.GetOrganizationId(), "organization_id", errCodeBadScopeID)
	if err != nil {
		return nil, err
	}
	result, err := h.reader.ListForOrganization(ctx, callerID, orgID, announcementread.ListFilter{
		Priority:        protoToModelPriority(req.GetPriority()),
		IncludeArchived: req.GetIncludeArchived(),
		Limit:           int(req.GetLimit()),
		Cursor:          req.Cursor,
	})
	if err != nil {
		return nil, err
	}
	items := toAnnouncementViews(result.Items)
	return &announcementqueryv1.ListAnnouncementsForOrganizationResponse{
		Items:      items,
		NextCursor: result.NextCursor,
	}, nil
}

func (h *AnnouncementQueryHandler) ListAnnouncementsForClinic(
	ctx context.Context, req *announcementqueryv1.ListAnnouncementsForClinicRequest,
) (*announcementqueryv1.ListAnnouncementsForClinicResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	clinicID, err := parseUUID(req.GetClinicId(), "clinic_id", errCodeBadScopeID)
	if err != nil {
		return nil, err
	}
	result, err := h.reader.ListForClinic(ctx, callerID, clinicID, announcementread.ListFilter{
		Priority:        protoToModelPriority(req.GetPriority()),
		IncludeArchived: req.GetIncludeArchived(),
		Limit:           int(req.GetLimit()),
		Cursor:          req.Cursor,
	})
	if err != nil {
		return nil, err
	}
	items := toAnnouncementViews(result.Items)
	return &announcementqueryv1.ListAnnouncementsForClinicResponse{
		Items:      items,
		NextCursor: result.NextCursor,
	}, nil
}

func (h *AnnouncementQueryHandler) ListAnnouncementsForDepartment(
	ctx context.Context, req *announcementqueryv1.ListAnnouncementsForDepartmentRequest,
) (*announcementqueryv1.ListAnnouncementsForDepartmentResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	deptID, err := parseUUID(req.GetDepartmentId(), "department_id", errCodeBadScopeID)
	if err != nil {
		return nil, err
	}
	result, err := h.reader.ListForDepartment(ctx, callerID, deptID, announcementread.ListFilter{
		Priority:        protoToModelPriority(req.GetPriority()),
		IncludeArchived: req.GetIncludeArchived(),
		Limit:           int(req.GetLimit()),
		Cursor:          req.Cursor,
	})
	if err != nil {
		return nil, err
	}
	items := toAnnouncementViews(result.Items)
	return &announcementqueryv1.ListAnnouncementsForDepartmentResponse{
		Items:      items,
		NextCursor: result.NextCursor,
	}, nil
}

func toAnnouncementViews(rows []announcementread.AnnouncementRow) []*announcementqueryv1.AnnouncementView {
	items := make([]*announcementqueryv1.AnnouncementView, len(rows))
	for i := range rows {
		items[i] = rowToProto(&rows[i])
	}
	return items
}
