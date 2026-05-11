package membership

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/samber/oops"
	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	"github.com/medincident/medincident-backend/internal/service/zitadel"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// GrantSystemAdminPayload carries the Zitadel user ID to promote to system admin.
type GrantSystemAdminPayload struct {
	ZitadelUserID string `validate:"required,no_extra_ws"`
}

// GrantSystemAdminCommand = caller + payload.
type GrantSystemAdminCommand struct {
	Caller  authz.Caller
	Payload GrantSystemAdminPayload
}

// GrantSystemAdmin creates a SystemAdmin grant for a Zitadel user.
// Verifies the user exists in Zitadel before writing. Not coupled to
// the employees table in any way. See spec §8.8.
//
// See: docs/services/Membership.md
func (s *EmployeeService) GrantSystemAdmin(ctx context.Context, cmd GrantSystemAdminCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return err
	}
	id := strings.TrimSpace(cmd.Payload.ZitadelUserID)

	if err := s.verifier.Verify(ctx, id); err != nil {
		if errors.Is(err, zitadel.ErrUserNotFound) {
			return oops.In(scopeSystemAdmin).
				Code(ErrCodeZitadelUserNotFound).
				Public("User does not exist in identity provider.").
				With("zitadel_user_id", id).
				Wrap(err)
		}
		return oops.In(scopeSystemAdmin).
			Code(zitadel.ErrCodeZitadelVerifyFailed).
			With("zitadel_user_id", id).
			Wrap(err)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := model.SystemAdmin{ZitadelUserID: id}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeSystemAdmin).
					Code(ErrCodeSystemAdminAlreadyGranted).
					Public("User is already a system admin.").
					With("zitadel_user_id", id).
					Wrap(err)
			}
			return oops.In(scopeSystemAdmin).Code(ErrCodeSystemAdminSaveFailed).Wrap(err)
		}

		now := time.Now().UTC()
		env, err := buildSystemAdminGrantedEnvelope(id, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.system_admin.v1.granted", env)
	})
}

func buildSystemAdminGrantedEnvelope(zitadelID string, grantedAt time.Time) (*eventv1.Envelope, error) {
	msg := &sav1.SystemAdminGranted{GrantedAt: timestamppb.New(grantedAt)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeSystemAdmin).Code(ErrCodeSystemAdminSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(grantedAt),
		AggregateType: "system_admin",
		AggregateId:   zitadelID,
		Payload:       payload,
	}, nil
}
