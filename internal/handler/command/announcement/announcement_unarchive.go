package announcement

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

func (h *AnnouncementHandler) UnarchiveAnnouncement(
	ctx context.Context, req *announcementv1.UnarchiveAnnouncementRequest,
) (*announcementv1.UnarchiveAnnouncementResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.Unarchive(ctx, &announcementsvc.UnarchiveAnnouncementCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: announcementsvc.UnarchiveAnnouncementPayload{ID: req.GetId()},
	}); err != nil {
		return nil, err
	}
	return &announcementv1.UnarchiveAnnouncementResponse{}, nil
}
