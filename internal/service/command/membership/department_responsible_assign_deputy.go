package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// AssignDepartmentResponsibleDeputyPayload carries the identifiers
// needed to set the deputy slot on an existing DR role.
type AssignDepartmentResponsibleDeputyPayload struct {
	DepartmentID     string `validate:"required,uuid"`
	EmployeeID       string `validate:"required,uuid"`
	DeputyEmployeeID string `validate:"required,uuid"`
}

// AssignDepartmentResponsibleDeputyCommand = caller + payload.
type AssignDepartmentResponsibleDeputyCommand struct {
	Caller  authz.Caller
	Payload AssignDepartmentResponsibleDeputyPayload
}

// AssignDepartmentResponsibleDeputy sets the deputy slot on an
// existing DR role. See spec §8.5.
func (s *EmployeeService) AssignDepartmentResponsibleDeputy(ctx context.Context, cmd AssignDepartmentResponsibleDeputyCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	deputyEmployeeID := uuid.MustParse(cmd.Payload.DeputyEmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(departmentID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.DepartmentResponsible
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("department_id = ? AND employee_id = ?", departmentID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleNotFound).
					Public("Department responsible not found.").
					With("department_id", departmentID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", deputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", deputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if deputy.DepartmentID != departmentID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyNotInDepartment).
				Public("Deputy employee does not belong to this department.").
				With("deputy_department_id", deputy.DepartmentID).
				With("target_department_id", departmentID).
				Errorf("deputy not in department")
		}

		if deputyEmployeeID == employeeID {
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

		row.DeputyEmployeeID = null.ValueFrom(deputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		return projector.DepartmentResponsibleDeputyAssigned(tx, &row)
	})
}
