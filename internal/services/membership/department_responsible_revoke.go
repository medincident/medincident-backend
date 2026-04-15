package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	departmentv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/department/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// RevokeDepartmentResponsibleCommand carries the identifiers needed to
// remove an employee's department responsible role.
type RevokeDepartmentResponsibleCommand struct {
	DepartmentID uuid.UUID
	EmployeeID   uuid.UUID
}

// RevokeDepartmentResponsible removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
func (s *EmployeeService) RevokeDepartmentResponsible(ctx context.Context, cmd RevokeDepartmentResponsibleCommand) error {
	var errs []error
	if cmd.DepartmentID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentIDEmpty).
			Public("Department ID is required.").Errorf("department id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, cmd.DepartmentID, cmd.EmployeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			cmd.DepartmentID, cmd.EmployeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}

		return publishDepartmentResponsibleRevoked(tx, cmd.DepartmentID, cmd.EmployeeID, now)
	})
}

// publishDepartmentResponsibleRevoked is a shared helper for Revoke,
// cascade-on-transfer, and cascade-on-terminate. It does NOT delete
// the row; the caller is responsible for that.
func publishDepartmentResponsibleRevoked(tx *gorm.DB, departmentID, employeeID uuid.UUID, now time.Time) error {
	ev := &departmentv1.DepartmentResponsibleRevoked{
		EmployeeId: employeeID.String(),
	}
	return outbox.Publish(tx, SubjectDepartmentResponsibleRevoked, AggregateTypeDepartment, departmentID.String(), now, ev)
}

// publishDepartmentResponsibleDeputyRemoved is a shared helper; it
// publishes the event only and does NOT update the row.
func publishDepartmentResponsibleDeputyRemoved(tx *gorm.DB, departmentID, employeeID uuid.UUID, now time.Time) error {
	ev := &departmentv1.DepartmentResponsibleDeputyRemoved{
		EmployeeId: employeeID.String(),
	}
	return outbox.Publish(tx, SubjectDepartmentResponsibleDeputyRemoved, AggregateTypeDepartment, departmentID.String(), now, ev)
}
