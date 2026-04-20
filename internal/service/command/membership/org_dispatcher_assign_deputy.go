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

// AssignOrganizationDispatcherDeputyCommand carries the identifiers needed
// to set the deputy slot on an existing OrgDispatcher role.
type AssignOrganizationDispatcherDeputyCommand struct {
	OrganizationID   uuid.UUID `validate:"required"`
	EmployeeID       uuid.UUID `validate:"required"`
	DeputyEmployeeID uuid.UUID `validate:"required"`
}

// AssignOrganizationDispatcherDeputy sets the deputy slot on an existing
// OrgDispatcher role. See spec §8.5 (OrgDispatcher variant).
func (s *EmployeeService) AssignOrganizationDispatcherDeputy(ctx context.Context, cmd AssignOrganizationDispatcherDeputyCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeOrgDispatcher).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeOrgDispatcher).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: organization_id is denormalized on employees — direct compare.
		if deputy.OrganizationID != cmd.OrganizationID {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyNotInOrganization).
				Public("Deputy employee does not belong to this organization.").
				With("deputy_employee_id", cmd.DeputyEmployeeID).
				With("deputy_organization_id", deputy.OrganizationID).
				With("target_organization_id", cmd.OrganizationID).
				Errorf("deputy not in target organization")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeOrgDispatcher).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		deputyID := cmd.DeputyEmployeeID
		row.DeputyEmployeeID = null.ValueFrom(deputyID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherSaveFailed).Wrap(err)
		}

		return projector.OrgDispatcherDeputyAssigned(tx, &row)
	})
}
