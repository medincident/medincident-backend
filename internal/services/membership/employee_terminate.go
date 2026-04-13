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

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

// TerminateEmployeeCommand carries the ID of the employee to remove.
type TerminateEmployeeCommand struct {
	ID uuid.UUID
}

// Terminate deletes the employee row. ON DELETE CASCADE on
// employee_vacations.employee_id drops any existing vacations in the
// same statement. Publishes EmployeeTerminated with an empty payload
// (aggregate_id in the envelope is sufficient for consumers).
func (s *EmployeeService) Terminate(ctx context.Context, cmd TerminateEmployeeCommand) error {
	if cmd.ID == uuid.Nil {
		return oops.In(scopeEmployee).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf(ErrCodeEmployeeIDEmpty)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var emp model.Employee
		if err := tx.First(&emp, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.ID).
					Errorf(ErrCodeEmployeeNotFound)
			}
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if err := tx.Delete(&model.Employee{}, "id = ?", cmd.ID).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeDeleteFailed).Wrap(err)
		}

		ev := &employeev1.EmployeeTerminated{}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeDeleteFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(time.Now().UTC()),
			AggregateType: AggregateTypeEmployee,
			AggregateId:   emp.ID.String(),
			Payload:       payload,
		}
		return AppendMembershipOutboxEvent(tx, SubjectEmployeeTerminated, envelope, nil)
	})
}
