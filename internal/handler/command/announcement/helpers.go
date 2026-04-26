package announcement

import (
	"time"

	"github.com/samber/oops"

	announcementv1 "github.com/medincident/medincident-backend/pkg/command/announcement/v1"
)

const handlerScope = "handler.command.announcement"

const errCodeBadTimestamp = "handler_invalid_timestamp"

func parseTimestamp(s, field string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, oops.In(handlerScope).
			Code(errCodeBadTimestamp).
			Public(field+" is not a valid RFC3339Nano timestamp.").
			With(field, s).
			Wrap(err)
	}
	return t, nil
}

func protoToModelPriority(p announcementv1.AnnouncementPriority) string {
	switch p {
	case announcementv1.AnnouncementPriority_ANNOUNCEMENT_PRIORITY_HIGH:
		return "high"
	default:
		return "normal"
	}
}
