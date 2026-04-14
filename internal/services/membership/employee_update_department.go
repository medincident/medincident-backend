package membership

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
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

		oldDepartmentID := emp.DepartmentID

		// Resolve both the old and new clinic IDs (plus the new department's
		// organization) in a single round-trip. We need all three values to
		// decide whether the move is cross-clinic and to validate that the
		// target department belongs to the employee's organisation.
		var lookup struct {
			OldClinicID uuid.UUID
			NewClinicID uuid.UUID
			NewOrgID    uuid.UUID
		}
		err := tx.Raw(`
			SELECT old_c.id           AS old_clinic_id,
			       new_c.id           AS new_clinic_id,
			       new_c.organization_id AS new_org_id
			FROM domain.departments old_d
			JOIN domain.clinics     old_c ON old_c.id = old_d.clinic_id
			JOIN domain.departments new_d ON new_d.id = ?
			JOIN domain.clinics     new_c ON new_c.id = new_d.clinic_id
			WHERE old_d.id = ?`, cmd.DepartmentID, oldDepartmentID,
		).Row().Scan(&lookup.OldClinicID, &lookup.NewClinicID, &lookup.NewOrgID)
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
		if lookup.NewOrgID != emp.OrganizationID {
			return oops.In(scopeEmployee).
				Code(ErrCodeDepartmentNotInSameOrganization).
				Public("Target department belongs to a different organization.").
				With("employee_organization_id", emp.OrganizationID).
				With("target_organization_id", lookup.NewOrgID).
				Errorf("department is in a different organization")
		}

		emp.DepartmentID = cmd.DepartmentID
		if err := tx.Save(&emp).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
		}

		ev := &employeev1.EmployeeDepartmentChanged{DepartmentId: cmd.DepartmentID.String()}
		if err := outbox.Publish(tx, SubjectEmployeeDepartmentChanged, AggregateTypeEmployee, emp.ID.String(), emp.UpdatedAt, ev); err != nil {
			return err
		}

		// Rule 1: cause first — DepartmentChanged is already in the outbox.
		// Cascade-revoke any DR roles the employee held in the old
		// department, and clear any DR deputy slots in the old department
		// that still reference this employee. The DR-deputy invariant is
		// same-department, so leaving the old department always invalidates
		// those slots.
		if err := cascadeRevokeDepartmentResponsible(tx, emp.ID, oldDepartmentID, emp.UpdatedAt); err != nil {
			return err
		}
		if err := cascadeClearDepartmentResponsibleDeputyInDepartment(tx, emp.ID, oldDepartmentID, emp.UpdatedAt); err != nil {
			return err
		}

		// Cross-clinic moves additionally invalidate CH roles and CH
		// deputy slots in the old clinic. Same-clinic moves leave CH
		// alone because the holder/deputy invariant is per-clinic.
		if lookup.NewClinicID != lookup.OldClinicID {
			if err := cascadeRevokeClinicHead(tx, emp.ID, lookup.OldClinicID, emp.UpdatedAt); err != nil {
				return err
			}
			if err := cascadeClearClinicHeadDeputyInClinic(tx, emp.ID, lookup.OldClinicID, emp.UpdatedAt); err != nil {
				return err
			}
		}
		return nil
	})
}
