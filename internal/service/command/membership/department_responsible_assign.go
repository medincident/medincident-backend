package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// AssignDepartmentResponsiblePayload carries the identifiers needed to
// link an employee to a department as its responsible.
type AssignDepartmentResponsiblePayload struct {
	DepartmentID string `validate:"required,uuid"`
	EmployeeID   string `validate:"required,uuid"`
}

// AssignDepartmentResponsibleCommand = caller + payload.
type AssignDepartmentResponsibleCommand struct {
	Caller  authz.Caller
	Payload AssignDepartmentResponsiblePayload
}

// AssignDepartmentResponsible links the employee to the department as
// a responsible. The employee must currently work in the department.
// See spec §8.2.
//
// See: docs/services/Membership.md
func (s *EmployeeService) AssignDepartmentResponsible(ctx context.Context, cmd AssignDepartmentResponsibleCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(departmentID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := requireDepartmentExists(tx, scopeDepartmentResponsible, departmentID); err != nil {
			return err
		}
		emp, err := loadEmployeeForRoleAssignment(tx, scopeDepartmentResponsible, employeeID)
		if err != nil {
			return err
		}
		if err := requireEmployeeInDepartment(scopeDepartmentResponsible, emp, departmentID); err != nil {
			return err
		}

		row := model.DepartmentResponsible{
			DepartmentID: departmentID,
			EmployeeID:   employeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleAlreadyAssigned).
					Public("This employee is already a department responsible.").
					With("department_id", departmentID).
					With("employee_id", employeeID).
					Wrap(err)
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		env, err := buildDeptResponsibleAssignedEnvelope(&row)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.department.v1.dept_responsible_assigned", env)
	})
}

func buildDeptResponsibleAssignedEnvelope(row *model.DepartmentResponsible) (*eventv1.Envelope, error) {
	msg := &deptv1.DeptResponsibleAssigned{
		EmployeeId: row.EmployeeID.String(),
		AssignedAt: timestamppb.New(row.CreatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(row.CreatedAt),
		AggregateType: "department",
		AggregateId:   row.DepartmentID.String(),
		Payload:       payload,
	}, nil
}
