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

// DeactivateEmployeePayload identifies the employee to deactivate.
type DeactivateEmployeePayload struct {
	ID string `validate:"required,uuid"`
}

// DeactivateEmployeeCommand = caller + payload.
type DeactivateEmployeeCommand struct {
	Caller  authz.Caller
	Payload DeactivateEmployeePayload
}

// Deactivate marks the employee as inactive and revokes all their roles.
//
// See: docs/services/Membership.md
func (s *EmployeeService) Deactivate(
	ctx context.Context,
	cmd DeactivateEmployeeCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	employeeID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Employee(employeeID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		if err := RevokeAllEmployeeRoles(tx, emp.ID, now); err != nil {
			return err
		}

		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.employees SET is_active = FALSE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			employeeID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeSaveFailed).
				With("employee_id", employeeID).
				Wrap(err)
		}
		return appendEmployeeDeactivatedEvent(tx, &emp, updatedAt)
	})
}

func appendEmployeeDeactivatedEvent(tx *gorm.DB, emp *model.Employee, updatedAt time.Time) error {
	msg := &empv1.EmployeeDeactivated{EmployeeId: emp.ID.String(), UpdatedAt: timestamppb.New(updatedAt)}
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
	return outbox.Append(tx, "medincident.event.employee.v1.deactivated", env)
}
