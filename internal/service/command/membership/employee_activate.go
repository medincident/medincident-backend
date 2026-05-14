package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// ActivateEmployeePayload identifies the employee to activate.
type ActivateEmployeePayload struct {
	ID string `validate:"required,uuid"`
}

// ActivateEmployeeCommand = caller + payload.
type ActivateEmployeeCommand struct {
	Caller  authz.Caller
	Payload ActivateEmployeePayload
}

// Activate marks a deactivated employee as active. The parent department
// must be active; roles are NOT restored.
//
// See: docs/services/Membership.md
func (s *EmployeeService) Activate(
	ctx context.Context,
	cmd ActivateEmployeeCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	employeeID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Employee(employeeID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var emp model.Employee
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&emp, "id = ?", employeeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", employeeID).
					Wrap(err)
			}
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if emp.IsActive {
			return nil
		}

		// Check that the parent department is active.
		var deptActive bool
		if err := tx.Raw(
			`SELECT is_active FROM domain.departments WHERE id = ?`,
			emp.DepartmentID,
		).Row().Scan(&deptActive); err != nil {
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeLoadFailed).
				With("employee_id", employeeID).
				Wrap(err)
		}
		if !deptActive {
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeActivateParentInactive).
				Public("Cannot activate employee: parent department is inactive.").
				With("employee_id", employeeID).
				With("department_id", emp.DepartmentID).
				Errorf("inactive parent department")
		}

		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.employees SET is_active = TRUE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			employeeID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeSaveFailed).
				With("employee_id", employeeID).
				Wrap(err)
		}
		return appendEmployeeActivatedEvent(tx, &emp, updatedAt)
	})
}

func appendEmployeeActivatedEvent(tx *gorm.DB, emp *model.Employee, updatedAt time.Time) error {
	msg := &empv1.EmployeeActivated{EmployeeId: emp.ID.String(), UpdatedAt: timestamppb.New(updatedAt)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(updatedAt),
		AggregateType: "employee",
		AggregateId:   emp.ID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.employee.v1.activated", env)
}
