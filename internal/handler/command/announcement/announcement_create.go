package announcement

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

func (h *AnnouncementHandler) CreateAnnouncement(
	ctx context.Context, req *announcementv1.CreateAnnouncementRequest,
) (*announcementv1.CreateAnnouncementResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}

	payload := announcementsvc.CreateAnnouncementPayload{
		OrganizationID: req.GetOrganizationId(),
		ClinicID:       req.ClinicId,
		DepartmentID:   req.DepartmentId,
		Title:          req.GetTitle(),
		Content:        req.GetContent(),
		Priority:       protoToModelPriority(req.GetPriority()),
	}
	if req.StartsAt != nil {
		t, err := parseTimestamp(*req.StartsAt, "starts_at")
		if err != nil {
			return nil, err
		}
		payload.StartsAt = &t
	}
	if req.EndsAt != nil {
		t, err := parseTimestamp(*req.EndsAt, "ends_at")
		if err != nil {
			return nil, err
		}
		payload.EndsAt = &t
	}

	result, err := h.svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}
	return &announcementv1.CreateAnnouncementResponse{Id: result.ID.String()}, nil
}
