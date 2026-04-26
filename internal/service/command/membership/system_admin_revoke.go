package membership

import (
	"context"
	"strings"

	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
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

		return projector.SystemAdminRevoked(tx, id)
	})
}
