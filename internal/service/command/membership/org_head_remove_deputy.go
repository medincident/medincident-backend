package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// RemoveOrganizationHeadDeputyPayload carries the identifiers needed
// to clear the deputy slot on an OrgHead role.
type RemoveOrganizationHeadDeputyPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// RemoveOrganizationHeadDeputyCommand = caller + payload.
type RemoveOrganizationHeadDeputyCommand struct {
	Caller  authz.Caller
	Payload RemoveOrganizationHeadDeputyPayload
}

// RemoveOrganizationHeadDeputy clears the deputy slot. Fails if the
// slot is already empty (no idempotent no-op per spec §4.7).
//
// See: docs/services/Membership.md
func (s *EmployeeService) RemoveOrganizationHeadDeputy(ctx context.Context, cmd RemoveOrganizationHeadDeputyCommand) error {
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

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("organization_id", organizationID).
				With("employee_id", employeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		return publishOrgHeadDeputyRemoved(tx, organizationID, employeeID, now)
	})
}
