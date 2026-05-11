package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

// TerminateEmployeePayload carries the ID of the employee to remove.
type TerminateEmployeePayload struct {
	ID string `validate:"required,uuid"`
}

// TerminateEmployeeCommand = caller + payload.
type TerminateEmployeeCommand struct {
	Caller  authz.Caller
	Payload TerminateEmployeePayload
}

// Terminate deletes the employee row. ON DELETE CASCADE on
// employee_vacations.employee_id drops any existing vacations in the
// same statement. Publishes EmployeeTerminated with an empty payload
// (aggregate_id in the envelope is sufficient for consumers).
//
// See: docs/services/Membership.md
func (s *EmployeeService) Terminate(ctx context.Context, cmd TerminateEmployeeCommand) error {
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
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", employeeID).
			First(&emp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeEmployee).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", employeeID).
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

		if err := tx.Delete(&model.Employee{}, "id = ?", employeeID).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeDeleteFailed).Wrap(err)
		}

		env, err := buildEmployeeTerminatedEnvelope(&emp, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.employee.v1.terminated", env)
	})
}

func buildEmployeeTerminatedEnvelope(emp *model.Employee, terminatedAt time.Time) (*eventv1.Envelope, error) {
	msg := &empv1.EmployeeTerminated{
		OrganizationId: emp.OrganizationID.String(),
		DepartmentId:   emp.DepartmentID.String(),
		TerminatedAt:   timestamppb.New(terminatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeEmployee).Code(ErrCodeEmployeeDeleteFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(terminatedAt),
		AggregateType: "employee",
		AggregateId:   emp.ID.String(),
		Payload:       payload,
	}, nil
}
