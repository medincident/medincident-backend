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
	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeactivateClinicPayload identifies the clinic to deactivate.
type DeactivateClinicPayload struct {
	ID string `validate:"required,uuid"`
}

// DeactivateClinicCommand = caller + payload.
type DeactivateClinicCommand struct {
	Caller  authz.Caller
	Payload DeactivateClinicPayload
}

// Deactivate performs cascading deactivation: the clinic, all its
// departments, and all employees in those departments (each loses all roles).
//
// See: docs/services/OrgStructure.md
func (s *ClinicService) Deactivate(
	ctx context.Context,
	cmd DeactivateClinicCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var clinic model.Clinic
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&clinic, "id = ?", clinicID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", clinicID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}

		// Collect active department IDs.
		var deptIDs []uuid.UUID
		if err := tx.Raw(
			`SELECT id FROM domain.departments WHERE clinic_id = ? AND is_active FOR UPDATE`,
			clinicID,
		).Scan(&deptIDs).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}

		if len(deptIDs) > 0 {
			// Collect active employee IDs in those departments.
			var empIDs []uuid.UUID
			if err := tx.Raw(
				`SELECT id FROM domain.employees WHERE department_id = ANY(?) AND is_active FOR UPDATE`,
				deptIDs,
			).Scan(&empIDs).Error; err != nil {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicLoadFailed).
					With("clinic_id", clinicID).
					Wrap(err)
			}

			for _, empID := range empIDs {
				if err := membership.RevokeAllEmployeeRoles(tx, empID, now); err != nil {
					return err
				}
			}

			// Deactivate employees.
			var deactivatedEmpIDs []uuid.UUID
			if err := tx.Raw(
				`UPDATE domain.employees SET is_active = FALSE, updated_at = now() WHERE id = ANY(?) AND is_active RETURNING id`,
				empIDs,
			).Scan(&deactivatedEmpIDs).Error; err != nil {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicSaveFailed).
					With("clinic_id", clinicID).
					Wrap(err)
			}
			for _, id := range deactivatedEmpIDs {
				if err := appendClinicCascadeEmployeeDeactivatedEvent(tx, id, now); err != nil {
					return err
				}
			}

			// Deactivate departments.
			var deactivatedDeptIDs []uuid.UUID
			if err := tx.Raw(
				`UPDATE domain.departments SET is_active = FALSE, updated_at = now() WHERE id = ANY(?) AND is_active RETURNING id`,
				deptIDs,
			).Scan(&deactivatedDeptIDs).Error; err != nil {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicSaveFailed).
					With("clinic_id", clinicID).
					Wrap(err)
			}
			for _, id := range deactivatedDeptIDs {
				if err := appendClinicCascadeDepartmentDeactivatedEvent(tx, id, now); err != nil {
					return err
				}
			}
		}

		// Deactivate the clinic itself (idempotent).
		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.clinics SET is_active = FALSE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			clinicID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}
		return appendClinicDeactivatedEvent(tx, clinicID, updatedAt)
	})
}

func appendClinicDeactivatedEvent(tx *gorm.DB, clinicID uuid.UUID, now time.Time) error {
	msg := &clinicv1.ClinicDeactivated{ClinicId: clinicID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.clinic").Code(ErrCodeClinicSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "clinic",
		AggregateId:   clinicID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.clinic.v1.deactivated", env)
}

func appendClinicCascadeDepartmentDeactivatedEvent(tx *gorm.DB, deptID uuid.UUID, now time.Time) error {
	msg := &deptv1.DepartmentDeactivated{DepartmentId: deptID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.clinic").Code(ErrCodeClinicSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "department",
		AggregateId:   deptID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.department.v1.deactivated", env)
}

func appendClinicCascadeEmployeeDeactivatedEvent(tx *gorm.DB, empID uuid.UUID, now time.Time) error {
	msg := &empv1.EmployeeDeactivated{EmployeeId: empID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.clinic").Code(ErrCodeClinicSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "employee",
		AggregateId:   empID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.employee.v1.deactivated", env)
}
