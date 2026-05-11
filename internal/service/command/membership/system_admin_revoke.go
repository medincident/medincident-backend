package membership

import (
	"context"
	"strings"
	"time"

	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// RevokeSystemAdminPayload carries the Zitadel user ID to remove from system admin.
type RevokeSystemAdminPayload struct {
	ZitadelUserID string `validate:"required,no_extra_ws"`
}

// RevokeSystemAdminCommand = caller + payload.
type RevokeSystemAdminCommand struct {
	Caller  authz.Caller
	Payload RevokeSystemAdminPayload
}

// RevokeSystemAdmin removes a SystemAdmin grant. No Zitadel verify —
// a stale admin entry should be removable even if the user has been
// deleted from Zitadel. See spec §8.9.
//
// See: docs/services/Membership.md
func (s *EmployeeService) RevokeSystemAdmin(ctx context.Context, cmd RevokeSystemAdminCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return err
	}
	id := strings.TrimSpace(cmd.Payload.ZitadelUserID)

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

		now := time.Now().UTC()
		env, err := buildSystemAdminRevokedEnvelope(id, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.system_admin.v1.revoked", env)
	})
}

func buildSystemAdminRevokedEnvelope(zitadelID string, revokedAt time.Time) (*eventv1.Envelope, error) {
	msg := &sav1.SystemAdminRevoked{RevokedAt: timestamppb.New(revokedAt)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeSystemAdmin).Code(ErrCodeSystemAdminDeleteFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(revokedAt),
		AggregateType: "system_admin",
		AggregateId:   zitadelID,
		Payload:       payload,
	}, nil
}
