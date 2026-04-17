package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

// AssignDepartmentResponsibleCommand carries the identifiers needed to
// link an employee to a department as its responsible.
type AssignDepartmentResponsibleCommand struct {
	DepartmentID uuid.UUID
	EmployeeID   uuid.UUID
}

// AssignDepartmentResponsible links the employee to the department as
// a responsible. The employee must currently work in the department.
// See spec §8.2.
func (s *EmployeeService) AssignDepartmentResponsible(ctx context.Context, cmd AssignDepartmentResponsibleCommand) error {
	var errs []error
	if cmd.DepartmentID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).
			Code(ErrCodeDepartmentIDEmpty).
			Public("Department ID is required.").
			Errorf("department id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireDepartmentExists(tx, scopeDepartmentResponsible, cmd.DepartmentID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeDepartmentResponsible, cmd.EmployeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInDepartment(scopeDepartmentResponsible, emp, cmd.DepartmentID); err != nil {
			return err
		}

		row := model.DepartmentResponsible{
			DepartmentID: cmd.DepartmentID,
			EmployeeID:   cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleAlreadyAssigned).
					Public("This employee is already a department responsible.").
					With("department_id", cmd.DepartmentID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		return projector.DepartmentResponsibleAssigned(tx, &row)
	})
}
