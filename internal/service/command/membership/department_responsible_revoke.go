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
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// RevokeDepartmentResponsiblePayload carries the identifiers needed to
// remove an employee's department responsible role.
type RevokeDepartmentResponsiblePayload struct {
	DepartmentID string `validate:"required,uuid"`
	EmployeeID   string `validate:"required,uuid"`
}

// RevokeDepartmentResponsibleCommand = caller + payload.
type RevokeDepartmentResponsibleCommand struct {
	Caller  authz.Caller
	Payload RevokeDepartmentResponsiblePayload
}

// RevokeDepartmentResponsible removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
//
// See: docs/services/Membership.md
func (s *EmployeeService) RevokeDepartmentResponsible(ctx context.Context, cmd RevokeDepartmentResponsibleCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(departmentID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var row model.DepartmentResponsible
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("department_id = ? AND employee_id = ?", departmentID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleNotFound).
					Public("Department responsible not found.").
					With("department_id", departmentID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
		}

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, departmentID, employeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			departmentID, employeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}

		return publishDepartmentResponsibleRevoked(tx, departmentID, employeeID, now)
	})
}

// publishDepartmentResponsibleRevoked is a shared helper for Revoke,
// cascade-on-transfer, and cascade-on-terminate. It does NOT delete
// the domain row; the caller owns that.
func publishDepartmentResponsibleRevoked(tx *gorm.DB, departmentID, employeeID uuid.UUID, now time.Time) error {
	env, err := buildDeptResponsibleRevokedEnvelope(departmentID, employeeID, now)
	if err != nil {
		return err
	}
	return outbox.Append(tx, "medincident.event.department.v1.dept_responsible_revoked", env)
}

// publishDepartmentResponsibleDeputyRemoved is a shared helper; it
// clears the projection's deputy slot. The caller owns the domain-
// row update.
func publishDepartmentResponsibleDeputyRemoved(tx *gorm.DB, departmentID, employeeID uuid.UUID, now time.Time) error {
	env, err := buildDeptResponsibleDeputyRemovedEnvelope(departmentID, employeeID, now)
	if err != nil {
		return err
	}
	return outbox.Append(tx, "medincident.event.department.v1.dept_responsible_deputy_removed", env)
}

func buildDeptResponsibleRevokedEnvelope(departmentID, employeeID uuid.UUID, now time.Time) (*eventv1.Envelope, error) {
	msg := &deptv1.DeptResponsibleRevoked{EmployeeId: employeeID.String()}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "department",
		AggregateId:   departmentID.String(),
		Payload:       payload,
	}, nil
}

func buildDeptResponsibleDeputyRemovedEnvelope(departmentID, employeeID uuid.UUID, now time.Time) (*eventv1.Envelope, error) {
	msg := &deptv1.DeptResponsibleDeputyRemoved{EmployeeId: employeeID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "department",
		AggregateId:   departmentID.String(),
		Payload:       payload,
	}, nil
}
