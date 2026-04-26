package announcement

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

func (h *AnnouncementHandler) ArchiveAnnouncement(
	ctx context.Context, req *announcementv1.ArchiveAnnouncementRequest,
) (*announcementv1.ArchiveAnnouncementResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.Archive(ctx, &announcementsvc.ArchiveAnnouncementCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: announcementsvc.ArchiveAnnouncementPayload{ID: req.GetId()},
	}); err != nil {
		return nil, err
	}
	return &announcementv1.ArchiveAnnouncementResponse{}, nil
}
