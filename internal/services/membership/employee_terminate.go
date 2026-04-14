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

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
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
			Errorf("employee id is empty")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var emp model.Employee
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", cmd.ID).
			First(&emp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.ID).
					Errorf("employee not found")
			}
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Cascade: revoke any DR/CH/OrgAdmin roles where this employee is
		// the holder, and clear any DR/CH/OrgAdmin roles where this employee
		// is a deputy. Both must run BEFORE the employee row is deleted
		// (FK constraints).
		if err := cascadeRevokeDepartmentResponsibleAll(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeClearDepartmentResponsibleDeputy(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeRevokeClinicHeadAll(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeClearClinicHeadDeputy(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeRevokeOrgAdminAll(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeClearOrgAdminDeputy(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeRevokeOrgHeadAll(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeClearOrgHeadDeputy(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeRevokeOrgDispatcherAll(tx, emp.ID, now); err != nil {
			return err
		}
		if err := cascadeClearOrgDispatcherDeputy(tx, emp.ID, now); err != nil {
			return err
		}

		res := tx.Delete(&model.Employee{}, "id = ?", cmd.ID)
		if res.Error != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeDeleteFailed).Wrap(res.Error)
		}
		if res.RowsAffected == 0 {
			return oops.In(scopeEmployee).
				Code(ErrCodeEmployeeNotFound).
				Public("Employee not found.").
				With("employee_id", cmd.ID).
				Errorf("employee not found")
		}

		ev := &employeev1.EmployeeTerminated{}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(now),
			AggregateType: AggregateTypeEmployee,
			AggregateId:   emp.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectEmployeeTerminated, envelope, nil)
	})
}
