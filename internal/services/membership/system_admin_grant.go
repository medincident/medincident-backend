package membership

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"gorm.io/gorm"

	systemadminv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/system_admin/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/pgerr"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
	"github.com/medincident/medincident-command-service/internal/services/zitadel"
)

// GrantSystemAdminCommand carries the Zitadel user ID to promote to system admin.
type GrantSystemAdminCommand struct {
	ZitadelUserID string
}

// GrantSystemAdmin creates a SystemAdmin grant for a Zitadel user.
// Verifies the user exists in Zitadel before writing. Not coupled to
// the employees table in any way. See spec §8.8.
func (s *EmployeeService) GrantSystemAdmin(ctx context.Context, cmd GrantSystemAdminCommand) error {
	id := strings.TrimSpace(cmd.ZitadelUserID)
	if id == "" {
		return oops.In(scopeSystemAdmin).
			Code(ErrCodeSystemAdminZitadelUserIDEmpty).
			Public("Zitadel user ID is required.").
			Errorf("zitadel user id empty")
	}

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
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerr.CodeUniqueViolation {
				return oops.In(scopeSystemAdmin).
					Code(ErrCodeSystemAdminAlreadyGranted).
					Public("User is already a system admin.").
					With("zitadel_user_id", id).
					Wrap(err)
			}
			return oops.In(scopeSystemAdmin).Code(ErrCodeSystemAdminSaveFailed).Wrap(err)
		}

		ev := &systemadminv1.SystemAdminGranted{}
		return outbox.Publish(tx, SubjectSystemAdminGranted, AggregateTypeSystemAdmin, id, time.Now().UTC(), ev)
	})
}
