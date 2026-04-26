package announcement

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

func (h *AnnouncementHandler) UpdateAnnouncement(
	ctx context.Context, req *announcementv1.UpdateAnnouncementRequest,
) (*announcementv1.UpdateAnnouncementResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}

	payload := announcementsvc.UpdateAnnouncementPayload{
		ID:      req.GetId(),
		Title:   req.GetTitle(),
		Content: req.GetContent(),
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

	if err := h.svc.Update(ctx, &announcementsvc.UpdateAnnouncementCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: payload,
	}); err != nil {
		return nil, err
	}
	return &announcementv1.UpdateAnnouncementResponse{}, nil
}
