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

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// RemoveOrganizationAdminDeputyCommand carries the identifiers needed
// to clear the deputy slot on an OrgAdmin role.
type RemoveOrganizationAdminDeputyCommand struct {
	OrganizationID uuid.UUID `validate:"required"`
	EmployeeID     uuid.UUID `validate:"required"`
}

// RemoveOrganizationAdminDeputy clears the deputy slot. Fails if the
// slot is already empty (no idempotent no-op per spec §4.7).
func (s *EmployeeService) RemoveOrganizationAdminDeputy(ctx context.Context, cmd RemoveOrganizationAdminDeputyCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.OrgAdmin
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("organization_id = ? AND employee_id = ?", cmd.OrganizationID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgAdmin).
					Code(ErrCodeOrganizationAdminNotFound).
					Public("Organization admin not found.").
					With("organization_id", cmd.OrganizationID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminLoadFailed).Wrap(err)
		}

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgAdmin).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("organization_id", cmd.OrganizationID).
				With("employee_id", cmd.EmployeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}

		return publishOrgAdminDeputyRemoved(tx, cmd.OrganizationID, cmd.EmployeeID, now)
	})
}
