package membership

import (
	"context"
	"errors"

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
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// AssignDepartmentResponsibleDeputyPayload carries the identifiers
// needed to set the deputy slot on an existing DR role.
type AssignDepartmentResponsibleDeputyPayload struct {
	DepartmentID     string `validate:"required,uuid"`
	EmployeeID       string `validate:"required,uuid"`
	DeputyEmployeeID string `validate:"required,uuid"`
}

// AssignDepartmentResponsibleDeputyCommand = caller + payload.
type AssignDepartmentResponsibleDeputyCommand struct {
	Caller  authz.Caller
	Payload AssignDepartmentResponsibleDeputyPayload
}

// AssignDepartmentResponsibleDeputy sets the deputy slot on an
// existing DR role. See spec §8.5.
//
// See: docs/services/Membership.md
func (s *EmployeeService) AssignDepartmentResponsibleDeputy(ctx context.Context, cmd AssignDepartmentResponsibleDeputyCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	deputyEmployeeID := uuid.MustParse(cmd.Payload.DeputyEmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(departmentID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		var deputy model.Employee
		err = tx.Where("id = ?", deputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", deputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if deputy.DepartmentID != departmentID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyNotInDepartment).
				Public("Deputy employee does not belong to this department.").
				With("deputy_department_id", deputy.DepartmentID).
				With("target_department_id", departmentID).
				Errorf("deputy not in department")
		}

		if deputyEmployeeID == employeeID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		row.DeputyEmployeeID = null.ValueFrom(deputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		env, err := buildDeptResponsibleDeputyAssignedEnvelope(&row)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.department.v1.dept_responsible_deputy_assigned", env)
	})
}

func buildDeptResponsibleDeputyAssignedEnvelope(row *model.DepartmentResponsible) (*eventv1.Envelope, error) {
	msg := &deptv1.DeptResponsibleDeputyAssigned{
		EmployeeId:       row.EmployeeID.String(),
		DeputyEmployeeId: row.DeputyEmployeeID.V.String(),
		UpdatedAt:        timestamppb.New(row.UpdatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(row.UpdatedAt),
		AggregateType: "department",
		AggregateId:   row.DepartmentID.String(),
		Payload:       payload,
	}, nil
}
