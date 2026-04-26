package announcement

import (
	"context"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

func (h *AnnouncementHandler) UpdateAnnouncementPriority(
	ctx context.Context, req *announcementv1.UpdateAnnouncementPriorityRequest,
) (*announcementv1.UpdateAnnouncementPriorityResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if err := h.svc.UpdatePriority(ctx, &announcementsvc.UpdateAnnouncementPriorityCommand{
		Caller: authz.Caller{ZitadelUserID: callerID},
		Payload: announcementsvc.UpdateAnnouncementPriorityPayload{
			ID:       req.GetId(),
			Priority: protoToModelPriority(req.GetPriority()),
		},
	}); err != nil {
		return nil, err
	}
	return &announcementv1.UpdateAnnouncementPriorityResponse{}, nil
}
