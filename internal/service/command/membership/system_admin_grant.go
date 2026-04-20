package membership

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
	"github.com/medincident/medincident-command-service/internal/service/zitadel"
)

// GrantSystemAdminPayload carries the Zitadel user ID to promote to system admin.
type GrantSystemAdminPayload struct {
	ZitadelUserID string `validate:"required"`
}

// GrantSystemAdminCommand = caller + payload.
type GrantSystemAdminCommand struct {
	Caller  authz.Caller
	Payload GrantSystemAdminPayload
}

// GrantSystemAdmin creates a SystemAdmin grant for a Zitadel user.
// Verifies the user exists in Zitadel before writing. Not coupled to
// the employees table in any way. See spec §8.8.
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

		return projector.SystemAdminGranted(tx, id, time.Now().UTC())
	})
}
