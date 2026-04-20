package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// AssignOrganizationHeadDeputyCommand carries the identifiers needed
// to set the deputy slot on an existing OrgHead role.
type AssignOrganizationHeadDeputyCommand struct {
	OrganizationID   uuid.UUID `validate:"required"`
	EmployeeID       uuid.UUID `validate:"required"`
	DeputyEmployeeID uuid.UUID `validate:"required"`
}

// AssignOrganizationHeadDeputy sets the deputy slot on an existing
// OrgHead role. See spec §8.5 (OrgHead variant).
func (s *EmployeeService) AssignOrganizationHeadDeputy(ctx context.Context, cmd AssignOrganizationHeadDeputyCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgHead).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgHead).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", cmd.DeputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("deputy not in target organization")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgHead).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		deputyID := cmd.DeputyEmployeeID
		row.DeputyEmployeeID = null.ValueFrom(deputyID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}

		return projector.OrgHeadDeputyAssigned(tx, &row)
	})
}
