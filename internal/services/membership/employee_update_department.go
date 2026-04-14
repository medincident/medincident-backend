package membership

import (
	"context"
	"database/sql"
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

// UpdateEmployeeDepartmentCommand carries the inputs required to move
// an employee to a different department.
type UpdateEmployeeDepartmentCommand struct {
	ID           uuid.UUID
	DepartmentID uuid.UUID
}

// UpdateDepartment moves an employee to a different department. The
// target department must belong to the same organisation as the
// current employee row; moving across organisations is forbidden.
func (s *EmployeeService) UpdateDepartment(ctx context.Context, cmd UpdateEmployeeDepartmentCommand) error {
	var errs []error
	if cmd.ID == uuid.Nil {
		errs = append(errs, oops.In(scopeEmployee).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id is empty"))
	}
	if cmd.DepartmentID == uuid.Nil {
		errs = append(errs, oops.In(scopeEmployee).
			Code(ErrCodeEmployeeDepartmentIDEmpty).
			Public("Department ID is required.").
			Errorf("department id is empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

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

		if cmd.DepartmentID == emp.DepartmentID {
			return nil
		}

		var targetOrgID uuid.UUID
		err := tx.Raw(`
			SELECT c.organization_id
			FROM domain.departments d
			JOIN domain.clinics c ON c.id = d.clinic_id
			WHERE d.id = ?`, cmd.DepartmentID).Row().Scan(&targetOrgID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return oops.In(scopeEmployee).
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", cmd.DepartmentID).
					Errorf("department not found")
			}
			return oops.In(scopeEmployee).
				Code(ErrCodeDepartmentLookupFailed).
				With("department_id", cmd.DepartmentID).
				Wrap(err)
		}
		if targetOrgID != emp.OrganizationID {
			return oops.In(scopeEmployee).
				Code(ErrCodeDepartmentNotInSameOrganization).
				Public("Target department belongs to a different organization.").
				With("employee_organization_id", emp.OrganizationID).
				With("target_organization_id", targetOrgID).
				Errorf("department is in a different organization")
		}

		emp.DepartmentID = cmd.DepartmentID
		if err := tx.Save(&emp).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
		}

		ev := &employeev1.EmployeeDepartmentChanged{DepartmentId: cmd.DepartmentID.String()}
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
		return outbox.AppendEvent(tx, SubjectEmployeeDepartmentChanged, envelope, nil)
	})
}
