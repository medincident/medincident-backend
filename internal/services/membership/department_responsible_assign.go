package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	departmentv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/department/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/pgerr"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
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
		var deptExists int64
		if err := tx.Raw(`SELECT count(*) FROM domain.departments WHERE id = ?`, cmd.DepartmentID).
			Scan(&deptExists).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentLookupFailed).Wrap(err)
		}
		if deptExists == 0 {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDepartmentNotFound).
				Public("Department not found.").
				With("department_id", cmd.DepartmentID).
				Errorf("department not found")
		}

		var emp model.Employee
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", cmd.EmployeeID).
			First(&emp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.EmployeeID).
					Errorf("employee not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if emp.DepartmentID != cmd.DepartmentID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeEmployeeNotInDepartment).
				Public("Employee does not belong to this department.").
				With("employee_id", cmd.EmployeeID).
				With("employee_department_id", emp.DepartmentID).
				With("target_department_id", cmd.DepartmentID).
				Errorf("employee not in target department")
		}

		row := model.DepartmentResponsible{
			DepartmentID: cmd.DepartmentID,
			EmployeeID:   cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerr.CodeUniqueViolation {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleAlreadyAssigned).
					Public("This employee is already a department responsible.").
					With("department_id", cmd.DepartmentID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		ev := &departmentv1.DepartmentResponsibleAssigned{
			EmployeeId: cmd.EmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(row.CreatedAt),
			AggregateType: AggregateTypeDepartment,
			AggregateId:   cmd.DepartmentID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectDepartmentResponsibleAssigned, envelope)
	})
}
