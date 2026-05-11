package membership

import (
	"context"
	"database/sql"
	"errors"

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

// UpdateEmployeeDepartmentPayload carries the inputs required to move
// an employee to a different department.
type UpdateEmployeeDepartmentPayload struct {
	ID           string `validate:"required,uuid"`
	DepartmentID string `validate:"required,uuid"`
}

// UpdateEmployeeDepartmentCommand = caller + payload.
type UpdateEmployeeDepartmentCommand struct {
	Caller  authz.Caller
	Payload UpdateEmployeeDepartmentPayload
}

// UpdateDepartment moves an employee to a different department. The
// target department must belong to the same organisation as the
// current employee row; moving across organisations is forbidden.
//
// See: docs/services/Membership.md
func (s *EmployeeService) UpdateDepartment(ctx context.Context, cmd UpdateEmployeeDepartmentCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	employeeID := uuid.MustParse(cmd.Payload.ID)
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Employee(employeeID)); err != nil {
		return err
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

		if departmentID == emp.DepartmentID {
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
			WHERE old_d.id = ?`, departmentID, oldDepartmentID,
		).Row().Scan(&lookup.OldClinicID, &lookup.NewClinicID, &lookup.NewOrgID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return oops.In(scopeEmployee).
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", departmentID).
					Errorf("department not found")
			}
			return oops.In(scopeEmployee).
				Code(ErrCodeDepartmentLookupFailed).
				With("department_id", departmentID).
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

		emp.DepartmentID = departmentID
		if err := tx.Save(&emp).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
		}

		env, err := buildEmployeeDeptChangedEnvelope(&emp)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.employee.v1.department_changed", env); err != nil {
			return err
		}

		// Rule 1: cause first — the department change is already
		// reflected in both the domain row and projections. Cascade-
		// revoke any DR roles the employee held in the old department,
		// and clear any DR deputy slots in the old department that
		// still reference this employee. The DR-deputy invariant is
		// same-department, so leaving the old department always
		// invalidates those slots.
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

func buildEmployeeDeptChangedEnvelope(emp *model.Employee) (*eventv1.Envelope, error) {
	msg := &empv1.EmployeeDepartmentChanged{
		NewDepartmentId: emp.DepartmentID.String(),
		UpdatedAt:       timestamppb.New(emp.UpdatedAt),
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
