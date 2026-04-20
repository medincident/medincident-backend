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

// AssignDepartmentResponsibleDeputyCommand carries the identifiers
// needed to set the deputy slot on an existing DR role.
type AssignDepartmentResponsibleDeputyCommand struct {
	DepartmentID     uuid.UUID `validate:"required"`
	EmployeeID       uuid.UUID `validate:"required"`
	DeputyEmployeeID uuid.UUID `validate:"required"`
}

// AssignDepartmentResponsibleDeputy sets the deputy slot on an
// existing DR role. See spec §8.5.
func (s *EmployeeService) AssignDepartmentResponsibleDeputy(ctx context.Context, cmd AssignDepartmentResponsibleDeputyCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.DepartmentResponsible
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("department_id = ? AND employee_id = ?", cmd.DepartmentID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleNotFound).
					Public("Department responsible not found.").
					With("department_id", cmd.DepartmentID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if deputy.DepartmentID != cmd.DepartmentID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyNotInDepartment).
				Public("Deputy employee does not belong to this department.").
				With("deputy_department_id", deputy.DepartmentID).
				With("target_department_id", cmd.DepartmentID).
				Errorf("deputy not in department")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		row.DeputyEmployeeID = null.ValueFrom(cmd.DeputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		return projector.DepartmentResponsibleDeputyAssigned(tx, &row)
	})
}
