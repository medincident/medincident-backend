package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// UpdateEmployeePositionCommand carries the inputs required to change an
// employee's position. Position nil means "clear the position".
type UpdateEmployeePositionCommand struct {
	ID       uuid.UUID
	Position *string
}

// UpdatePosition changes an employee's position. Idempotent: if the
// normalised new position equals the current one, no event is emitted.
// Takes a FOR UPDATE row lock on the employee to serialise concurrent
// updates; last-writer-wins is NOT acceptable because it would emit
// events describing overwritten intermediate states.
func (s *EmployeeService) UpdatePosition(ctx context.Context, cmd UpdateEmployeePositionCommand) error {
	var errs []error
	if cmd.ID == uuid.Nil {
		errs = append(errs, oops.In(scopeEmployee).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id is empty"))
	}
	if err := validatePosition(cmd.Position); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	newPos := normalisePosition(cmd.Position)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var emp model.Employee
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", cmd.ID).
			First(&emp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.ID).
					Errorf("employee not found")
			}
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if emp.Position == newPos {
			return nil
		}
		emp.Position = newPos
		if err := tx.Save(&emp).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
		}

		ev := &employeev1.EmployeePositionChanged{}
		if newPos.Valid {
			p := newPos.String
			ev.Position = &p
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(emp.UpdatedAt),
			AggregateType: AggregateTypeEmployee,
			AggregateId:   emp.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectEmployeePositionChanged, envelope, nil)
	})
}
