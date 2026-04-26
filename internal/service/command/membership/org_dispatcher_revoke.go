package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// RevokeOrganizationDispatcherPayload carries the identifiers needed to
// remove an employee's organization dispatcher role.
type RevokeOrganizationDispatcherPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// RevokeOrganizationDispatcherCommand = caller + payload.
type RevokeOrganizationDispatcherCommand struct {
	Caller  authz.Caller
	Payload RevokeOrganizationDispatcherPayload
}

// RevokeOrganizationDispatcher removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
//
// See: docs/services/Membership.md
func (s *EmployeeService) RevokeOrganizationDispatcher(ctx context.Context, cmd RevokeOrganizationDispatcherCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	organizationID := uuid.MustParse(cmd.Payload.OrganizationID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(organizationID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.OrgDispatcher
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", organizationID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeOrganizationDispatcherNotFound).
					Public("Organization dispatcher not found.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherLoadFailed).Wrap(err)
		}

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishOrgDispatcherDeputyRemoved(tx, organizationID, employeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.OrgDispatcher{}, "organization_id = ? AND employee_id = ?",
			organizationID, employeeID).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherDeleteFailed).Wrap(err)
		}

		return publishOrgDispatcherRevoked(tx, organizationID, employeeID, now)
	})
}

// publishOrgDispatcherRevoked is a shared helper for Revoke and
// cascade-on-terminate. It does NOT delete the domain row; the caller
// owns that. The projector call removes the projection row so explicit
// and cascaded revokes stay in sync.
func publishOrgDispatcherRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID, _ time.Time) error {
	return projector.OrgDispatcherRevoked(tx, organizationID, employeeID)
}

// publishOrgDispatcherDeputyRemoved is a shared helper; it clears the
// projection's deputy slot. Caller owns the domain-row update.
func publishOrgDispatcherDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	return projector.OrgDispatcherDeputyRemoved(tx, organizationID, employeeID, now)
}
