package membership

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateEmployeePositionPayload carries the inputs required to change an
// employee's position. Position nil means "clear the position".
type UpdateEmployeePositionPayload struct {
	ID       string  `validate:"required,uuid"`
	Position *string `validate:"omitnil,no_extra_ws,min=2,max=256"`
}

// UpdateEmployeePositionCommand = caller + payload.
type UpdateEmployeePositionCommand struct {
	Caller  authz.Caller
	Payload UpdateEmployeePositionPayload
}

// UpdatePosition changes an employee's position. Idempotent: if the
// normalised new position equals the current one, no event is emitted.
// Takes a FOR UPDATE row lock on the employee to serialise concurrent
// updates; last-writer-wins is NOT acceptable because it would emit
// events describing overwritten intermediate states.
//
// See: docs/services/Membership.md
func (s *EmployeeService) UpdatePosition(ctx context.Context, cmd UpdateEmployeePositionCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	employeeID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Employee(employeeID)); err != nil {
		return err
	}

	var newPos null.String
	if cmd.Payload.Position != nil {
		newPos = null.StringFrom(strings.TrimSpace(*cmd.Payload.Position))
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var emp model.Employee
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", employeeID).
			First(&emp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", employeeID).
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

		env, err := buildEmployeePosChangedEnvelope(&emp)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.employee.v1.position_changed", env)
	})
}

func buildEmployeePosChangedEnvelope(emp *model.Employee) (*eventv1.Envelope, error) {
	var pos string
	if emp.Position.Valid {
		pos = emp.Position.String
	}
	msg := &empv1.EmployeePositionChanged{
		Position:  pos,
		UpdatedAt: timestamppb.New(emp.UpdatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(emp.UpdatedAt),
		AggregateType: "employee",
		AggregateId:   emp.ID.String(),
		Payload:       payload,
	}, nil
}
