package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/outbox"
	organizationv1 "github.com/medincident/medincident-command-service/pkg/event/organization/v1"
)

// RevokeOrganizationHeadCommand carries the identifiers needed to
// remove an employee's organization head role.
type RevokeOrganizationHeadCommand struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

// RevokeOrganizationHead removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
func (s *EmployeeService) RevokeOrganizationHead(ctx context.Context, cmd RevokeOrganizationHeadCommand) error {
	if err := validateOrgRoleKeys(scopeOrgHead, cmd.OrganizationID, cmd.EmployeeID); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.OrgHead
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", cmd.OrganizationID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeOrganizationHeadNotFound).
					Public("Organization head not found.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadLoadFailed).Wrap(err)
		}

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishOrgHeadDeputyRemoved(tx, cmd.OrganizationID, cmd.EmployeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.OrgHead{}, "organization_id = ? AND employee_id = ?",
			cmd.OrganizationID, cmd.EmployeeID).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadDeleteFailed).Wrap(err)
		}

		return publishOrgHeadRevoked(tx, cmd.OrganizationID, cmd.EmployeeID, now)
	})
}

// publishOrgHeadRevoked is a shared helper for Revoke,
// cascade-on-terminate. It does NOT delete the row; the caller is
// responsible for that.
func publishOrgHeadRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	ev := &organizationv1.OrganizationHeadRevoked{
		EmployeeId: employeeID.String(),
	}
	return outbox.Publish(tx, SubjectOrganizationHeadRevoked, AggregateTypeOrganization, organizationID.String(), now, ev)
}

// publishOrgHeadDeputyRemoved is a shared helper; it publishes the
// event only and does NOT update the row.
func publishOrgHeadDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	ev := &organizationv1.OrganizationHeadDeputyRemoved{
		EmployeeId: employeeID.String(),
	}
	return outbox.Publish(tx, SubjectOrganizationHeadDeputyRemoved, AggregateTypeOrganization, organizationID.String(), now, ev)
}
