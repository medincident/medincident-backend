package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	organizationv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/organization/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
)

// RevokeOrganizationDispatcherCommand carries the identifiers needed to
// remove an employee's organization dispatcher role.
type RevokeOrganizationDispatcherCommand struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

// RevokeOrganizationDispatcher removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
func (s *EmployeeService) RevokeOrganizationDispatcher(ctx context.Context, cmd RevokeOrganizationDispatcherCommand) error {
	if err := validateOrgRoleKeys(scopeOrgDispatcher, cmd.OrganizationID, cmd.EmployeeID); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.OrgDispatcher
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", cmd.OrganizationID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeOrganizationDispatcherNotFound).
					Public("Organization dispatcher not found.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherLoadFailed).Wrap(err)
		}

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishOrgDispatcherDeputyRemoved(tx, cmd.OrganizationID, cmd.EmployeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.OrgDispatcher{}, "organization_id = ? AND employee_id = ?",
			cmd.OrganizationID, cmd.EmployeeID).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherDeleteFailed).Wrap(err)
		}

		return publishOrgDispatcherRevoked(tx, cmd.OrganizationID, cmd.EmployeeID, now)
	})
}

// publishOrgDispatcherRevoked is a shared helper for Revoke,
// cascade-on-terminate. It does NOT delete the row; the caller is
// responsible for that.
func publishOrgDispatcherRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	ev := &organizationv1.OrganizationDispatcherRevoked{
		EmployeeId: employeeID.String(),
	}
	return outbox.Publish(tx, SubjectOrganizationDispatcherRevoked, AggregateTypeOrganization, organizationID.String(), now, ev)
}

// publishOrgDispatcherDeputyRemoved is a shared helper; it publishes the
// event only and does NOT update the row.
func publishOrgDispatcherDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	ev := &organizationv1.OrganizationDispatcherDeputyRemoved{
		EmployeeId: employeeID.String(),
	}
	return outbox.Publish(tx, SubjectOrganizationDispatcherDeputyRemoved, AggregateTypeOrganization, organizationID.String(), now, ev)
}
