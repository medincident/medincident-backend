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
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// RevokeOrganizationHeadPayload carries the identifiers needed to
// remove an employee's organization head role.
type RevokeOrganizationHeadPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// RevokeOrganizationHeadCommand = caller + payload.
type RevokeOrganizationHeadCommand struct {
	Caller  authz.Caller
	Payload RevokeOrganizationHeadPayload
}

// RevokeOrganizationHead removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
func (s *EmployeeService) RevokeOrganizationHead(ctx context.Context, cmd RevokeOrganizationHeadCommand) error {
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

		var row model.OrgHead
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", organizationID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeOrganizationHeadNotFound).
					Public("Organization head not found.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadLoadFailed).Wrap(err)
		}

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishOrgHeadDeputyRemoved(tx, organizationID, employeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.OrgHead{}, "organization_id = ? AND employee_id = ?",
			organizationID, employeeID).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadDeleteFailed).Wrap(err)
		}

		return publishOrgHeadRevoked(tx, organizationID, employeeID, now)
	})
}

// publishOrgHeadRevoked is a shared helper for Revoke and
// cascade-on-terminate. It does NOT delete the domain row; the caller
// owns that. The projector call removes the projection row so
// explicit and cascaded revokes stay in sync.
func publishOrgHeadRevoked(tx *gorm.DB, organizationID, employeeID uuid.UUID, _ time.Time) error {
	return projector.OrgHeadRevoked(tx, organizationID, employeeID)
}

// publishOrgHeadDeputyRemoved is a shared helper; it clears the
// projection's deputy slot. Caller owns the domain-row update.
func publishOrgHeadDeputyRemoved(tx *gorm.DB, organizationID, employeeID uuid.UUID, now time.Time) error {
	return projector.OrgHeadDeputyRemoved(tx, organizationID, employeeID, now)
}
