package membership

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// cascadeRevokeDepartmentResponsible revokes every DR role where the
// employee is the holder on the given department_id. For each revoked
// row:
//   - if it carries a deputy, publish DepartmentResponsibleDeputyRemoved first
//   - delete the row
//   - publish DepartmentResponsibleRevoked
func cascadeRevokeDepartmentResponsible(tx *gorm.DB, employeeID, departmentID uuid.UUID, now time.Time) error {
	var rows []model.DepartmentResponsible
	err := tx.Where("employee_id = ? AND department_id = ?", employeeID, departmentID).
		Find(&rows).Error
	if err != nil {
		return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			row.DepartmentID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}
		if err := publishDepartmentResponsibleRevoked(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeDepartmentResponsibleAll is the same as
// cascadeRevokeDepartmentResponsible, but unconstrained by
// department_id. Used by TerminateEmployee.
func cascadeRevokeDepartmentResponsibleAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.DepartmentResponsible
	err := tx.Where("employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			row.DepartmentID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}
		if err := publishDepartmentResponsibleRevoked(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearDepartmentResponsibleDeputy clears the deputy slot on
// every DR role where the given employee is the deputy. Used by
// TerminateEmployee to handle rows where the terminating employee was
// the deputy of someone else's role.
func cascadeClearDepartmentResponsibleDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.DepartmentResponsible
	err := tx.Where("deputy_employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if err := tx.Exec(
			`UPDATE domain.department_responsibles SET deputy_employee_id = NULL, updated_at = now() WHERE department_id = ? AND employee_id = ?`,
			row.DepartmentID, row.EmployeeID,
		).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}
		if err := publishDepartmentResponsibleDeputyRemoved(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}
