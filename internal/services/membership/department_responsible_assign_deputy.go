package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	departmentv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/department/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// AssignDepartmentResponsibleDeputyCommand carries the identifiers
// needed to set the deputy slot on an existing DR role.
type AssignDepartmentResponsibleDeputyCommand struct {
	DepartmentID     uuid.UUID
	EmployeeID       uuid.UUID
	DeputyEmployeeID uuid.UUID
}

// AssignDepartmentResponsibleDeputy sets the deputy slot on an
// existing DR role. See spec §8.5.
func (s *EmployeeService) AssignDepartmentResponsibleDeputy(ctx context.Context, cmd AssignDepartmentResponsibleDeputyCommand) error {
	var errs []error
	if cmd.DepartmentID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentIDEmpty).
			Public("Department ID is required.").Errorf("department id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if cmd.DeputyEmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeDeputyEmployeeIDEmpty).
			Public("Deputy employee ID is required.").Errorf("deputy id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.DepartmentResponsible
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("department_id = ? AND employee_id = ?", cmd.DepartmentID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDepartmentResponsibleNotFound).
					Public("Department responsible not found.").
					With("department_id", cmd.DepartmentID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeDepartmentResponsible).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		if deputy.DepartmentID != cmd.DepartmentID {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyNotInDepartment).
				Public("Deputy employee does not belong to this department.").
				With("deputy_department_id", deputy.DepartmentID).
				With("target_department_id", cmd.DepartmentID).
				Errorf("deputy not in department")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
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

		row.DeputyEmployeeID = null.ValueFrom(cmd.DeputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		ev := &departmentv1.DepartmentResponsibleDeputyAssigned{
			EmployeeId:       cmd.EmployeeID.String(),
			DeputyEmployeeId: cmd.DeputyEmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(time.Now().UTC()),
			AggregateType: AggregateTypeDepartment,
			AggregateId:   cmd.DepartmentID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectDepartmentResponsibleDeputyAssigned, envelope, nil)
	})
}
