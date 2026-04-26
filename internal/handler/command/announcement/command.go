// Package announcement is the gRPC transport for AnnouncementCommandService.
package announcement

import (
	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

// AnnouncementHandler implements announcementv1.AnnouncementCommandServiceServer.
type AnnouncementHandler struct {
	announcementv1.UnimplementedAnnouncementCommandServiceServer
	svc *announcementsvc.AnnouncementService
}

// NewAnnouncementHandler wires the handler with the service.
func NewAnnouncementHandler(svc *announcementsvc.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{svc: svc}
}
