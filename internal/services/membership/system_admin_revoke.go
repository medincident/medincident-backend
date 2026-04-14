package membership

import (
	"context"
	"strings"
	"time"

	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	systemadminv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/system_admin/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// RevokeSystemAdminCommand carries the Zitadel user ID to remove from system admin.
type RevokeSystemAdminCommand struct {
	ZitadelUserID string
}

// RevokeSystemAdmin removes a SystemAdmin grant. No Zitadel verify —
// a stale admin entry should be removable even if the user has been
// deleted from Zitadel. See spec §8.9.
func (s *EmployeeService) RevokeSystemAdmin(ctx context.Context, cmd RevokeSystemAdminCommand) error {
	id := strings.TrimSpace(cmd.ZitadelUserID)
	if id == "" {
		return oops.In(scopeSystemAdmin).
			Code(ErrCodeSystemAdminZitadelUserIDEmpty).
			Public("Zitadel user ID is required.").
			Errorf("zitadel user id empty")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&model.SystemAdmin{}, "zitadel_user_id = ?", id)
		if res.Error != nil {
			return oops.In(scopeSystemAdmin).Code(ErrCodeSystemAdminDeleteFailed).Wrap(res.Error)
		}
		if res.RowsAffected == 0 {
			return oops.In(scopeSystemAdmin).
				Code(ErrCodeSystemAdminNotFound).
				Public("System admin not found.").
				With("zitadel_user_id", id).
				Errorf("not found")
		}

		ev := &systemadminv1.SystemAdminRevoked{}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeSystemAdmin).Code(ErrCodeSystemAdminEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(time.Now().UTC()),
			AggregateType: AggregateTypeSystemAdmin,
			AggregateId:   id,
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectSystemAdminRevoked, envelope)
	})
}
