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

// RemoveOrganizationAdminDeputyPayload carries the identifiers needed
// to clear the deputy slot on an OrgAdmin role.
type RemoveOrganizationAdminDeputyPayload struct {
	OrganizationID string `validate:"required,uuid"`
	EmployeeID     string `validate:"required,uuid"`
}

// RemoveOrganizationAdminDeputyCommand = caller + payload.
type RemoveOrganizationAdminDeputyCommand struct {
	Caller  authz.Caller
	Payload RemoveOrganizationAdminDeputyPayload
}

// RemoveOrganizationAdminDeputy clears the deputy slot. Fails if the
// slot is already empty (no idempotent no-op per spec §4.7).
//
// See: docs/services/Membership.md
func (s *EmployeeService) RemoveOrganizationAdminDeputy(ctx context.Context, cmd RemoveOrganizationAdminDeputyCommand) error {
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

		var row model.OrgAdmin
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", organizationID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeOrganizationAdminNotFound).
					Public("Organization admin not found.").
					With("organization_id", organizationID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminLoadFailed).Wrap(err)
		}

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("organization_id", organizationID).
				With("employee_id", employeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}

		return publishOrgAdminDeputyRemoved(tx, organizationID, employeeID, now)
	})
}
