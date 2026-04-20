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

// RemoveOrganizationHeadDeputyCommand carries the identifiers needed
// to clear the deputy slot on an OrgHead role.
type RemoveOrganizationHeadDeputyCommand struct {
	OrganizationID uuid.UUID `validate:"required"`
	EmployeeID     uuid.UUID `validate:"required"`
}

// RemoveOrganizationHeadDeputy clears the deputy slot. Fails if the
// slot is already empty (no idempotent no-op per spec §4.7).
func (s *EmployeeService) RemoveOrganizationHeadDeputy(ctx context.Context, cmd RemoveOrganizationHeadDeputyCommand) error {
	if err := validation.Struct(cmd); err != nil {
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

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("organization_id", cmd.OrganizationID).
				With("employee_id", cmd.EmployeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		return publishOrgHeadDeputyRemoved(tx, cmd.OrganizationID, cmd.EmployeeID, now)
	})
}
