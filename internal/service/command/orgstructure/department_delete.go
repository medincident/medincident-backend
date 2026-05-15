package orgstructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeleteDepartmentPayload identifies the department to delete.
type DeleteDepartmentPayload struct {
	ID string `validate:"required,uuid"`
}

// DeleteDepartmentCommand = caller + payload.
type DeleteDepartmentCommand struct {
	Caller  authz.Caller
	Payload DeleteDepartmentPayload
}

// Delete hard-deletes a department and all its employees in the correct FK
// order: roles → employees → department.
//
// See: docs/services/OrgStructure.md
func (s *DepartmentService) Delete(
	ctx context.Context,
	cmd DeleteDepartmentCommand,
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

		// Collect and delete all employees in this department.
		var empIDs []uuid.UUID
		if err := tx.Raw(
			`SELECT id FROM domain.employees WHERE department_id = ? FOR UPDATE`,
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
		if len(empIDs) > 0 {
			if err := tx.Exec(
				`DELETE FROM domain.employees WHERE id IN (?)`, empIDs,
			).Error; err != nil {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentDeleteFailed).
					With("department_id", deptID).
					Wrap(err)
			}
		}

		// Delete the department itself.
		if err := tx.Delete(&model.Department{}, "id = ?", deptID).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentDeleteHasDependents).
					Public("Department has dependents and cannot be deleted.").
					With("department_id", deptID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentDeleteFailed).
				With("department_id", deptID).
				Wrap(err)
		}

		return appendDepartmentDeletedEvent(tx, deptID, now)
	})
}

func appendDepartmentDeletedEvent(tx *gorm.DB, deptID uuid.UUID, now time.Time) error {
	msg := &deptv1.DepartmentDeleted{DepartmentId: deptID.String(), DeletedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.department").Code(ErrCodeDepartmentDeleteFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "department",
		AggregateId:   deptID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.department.v1.deleted", env)
}
