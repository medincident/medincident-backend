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
	"github.com/medincident/medincident-backend/internal/service/validation"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// ActivateDepartmentPayload identifies the department to activate.
type ActivateDepartmentPayload struct {
	ID string `validate:"required,uuid"`
}

// ActivateDepartmentCommand = caller + payload.
type ActivateDepartmentCommand struct {
	Caller  authz.Caller
	Payload ActivateDepartmentPayload
}

// Activate marks a deactivated department as active. The parent clinic
// must be active; employees keep their individual is_active state.
//
// See: docs/services/OrgStructure.md
func (s *DepartmentService) Activate(
	ctx context.Context,
	cmd ActivateDepartmentCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	deptID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(deptID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		if dept.IsActive {
			return nil
		}

		// Check that the parent clinic is active.
		var clinicActive bool
		if err := tx.Raw(
			`SELECT is_active FROM domain.clinics WHERE id = ?`,
			dept.ClinicID,
		).Row().Scan(&clinicActive); err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", deptID).
				Wrap(err)
		}
		if !clinicActive {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentActivateParentInactive).
				Public("Cannot activate department: parent clinic is inactive.").
				With("department_id", deptID).
				With("clinic_id", dept.ClinicID).
				Errorf("inactive parent clinic")
		}

		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.departments SET is_active = TRUE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			deptID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", deptID).
				Wrap(err)
		}
		return appendDepartmentActivatedEvent(tx, deptID, updatedAt)
	})
}

func appendDepartmentActivatedEvent(tx *gorm.DB, deptID uuid.UUID, now time.Time) error {
	msg := &deptv1.DepartmentActivated{DepartmentId: deptID.String(), UpdatedAt: timestamppb.New(now)}
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
	return outbox.Append(tx, "medincident.event.department.v1.activated", env)
}
