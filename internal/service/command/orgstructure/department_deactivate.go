package orgstructure

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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	"github.com/medincident/medincident-backend/internal/service/validation"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeactivateDepartmentPayload identifies the department to deactivate.
type DeactivateDepartmentPayload struct {
	ID string `validate:"required,uuid"`
}

// DeactivateDepartmentCommand = caller + payload.
type DeactivateDepartmentCommand struct {
	Caller  authz.Caller
	Payload DeactivateDepartmentPayload
}

// Deactivate performs cascading deactivation: the department and all its
// employees (each loses all roles).
//
// See: docs/services/OrgStructure.md
func (s *DepartmentService) Deactivate(
	ctx context.Context,
	cmd DeactivateDepartmentCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	deptID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(deptID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var dept model.Department
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&dept, "id = ?", deptID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", deptID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", deptID).
				Wrap(err)
		}

		// Collect active employee IDs.
		var empIDs []uuid.UUID
		if err := tx.Raw(
			`SELECT id FROM domain.employees WHERE department_id = ? AND is_active FOR UPDATE`,
			deptID,
		).Scan(&empIDs).Error; err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", deptID).
				Wrap(err)
		}

		for _, empID := range empIDs {
			if err := membership.RevokeAllEmployeeRoles(tx, empID, now); err != nil {
				return err
			}
		}

		// Deactivate employees.
		var deactivatedEmpIDs []uuid.UUID
		if len(empIDs) > 0 {
			if err := tx.Raw(
				`UPDATE domain.employees SET is_active = FALSE, updated_at = now() WHERE id = ANY(?) AND is_active RETURNING id`,
				empIDs,
			).Scan(&deactivatedEmpIDs).Error; err != nil {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentSaveFailed).
					With("department_id", deptID).
					Wrap(err)
			}
		}
		for _, id := range deactivatedEmpIDs {
			if err := appendDeptCascadeEmployeeDeactivatedEvent(tx, id, now); err != nil {
				return err
			}
		}

		// Deactivate the department itself (idempotent).
		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.departments SET is_active = FALSE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			deptID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", deptID).
				Wrap(err)
		}
		return appendDepartmentDeactivatedEvent(tx, deptID, updatedAt)
	})
}

func appendDepartmentDeactivatedEvent(tx *gorm.DB, deptID uuid.UUID, now time.Time) error {
	msg := &deptv1.DepartmentDeactivated{DepartmentId: deptID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.department").Code(ErrCodeDepartmentSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "department",
		AggregateId:   deptID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.department.v1.deactivated", env)
}

func appendDeptCascadeEmployeeDeactivatedEvent(tx *gorm.DB, empID uuid.UUID, now time.Time) error {
	msg := &empv1.EmployeeDeactivated{EmployeeId: empID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.department").Code(ErrCodeDepartmentSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "employee",
		AggregateId:   empID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.employee.v1.deactivated", env)
}
