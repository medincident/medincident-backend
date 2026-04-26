package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// RevokeDepartmentResponsiblePayload carries the identifiers needed to
// remove an employee's department responsible role.
type RevokeDepartmentResponsiblePayload struct {
	DepartmentID string `validate:"required,uuid"`
	EmployeeID   string `validate:"required,uuid"`
}

// RevokeDepartmentResponsibleCommand = caller + payload.
type RevokeDepartmentResponsibleCommand struct {
	Caller  authz.Caller
	Payload RevokeDepartmentResponsiblePayload
}

// RevokeDepartmentResponsible removes the role row and publishes the
// Revoked event. If a deputy was assigned, a DeputyRemoved event is
// published FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
//
// See: docs/services/Membership.md
func (s *EmployeeService) RevokeDepartmentResponsible(ctx context.Context, cmd RevokeDepartmentResponsibleCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	departmentID := uuid.MustParse(cmd.Payload.DepartmentID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(departmentID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, departmentID, employeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			departmentID, employeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}

		return publishDepartmentResponsibleRevoked(tx, departmentID, employeeID, now)
	})
}

// publishDepartmentResponsibleRevoked is a shared helper for Revoke,
// cascade-on-transfer, and cascade-on-terminate. It does NOT delete
// the domain row; the caller owns that. The projector call removes
// the projection row so explicit and cascaded revokes stay in sync.
func publishDepartmentResponsibleRevoked(tx *gorm.DB, departmentID, employeeID uuid.UUID, _ time.Time) error {
	return projector.DepartmentResponsibleRevoked(tx, departmentID, employeeID)
}

// publishDepartmentResponsibleDeputyRemoved is a shared helper; it
// clears the projection's deputy slot. The caller owns the domain-
// row update.
func publishDepartmentResponsibleDeputyRemoved(tx *gorm.DB, departmentID, employeeID uuid.UUID, now time.Time) error {
	return projector.DepartmentResponsibleDeputyRemoved(tx, departmentID, employeeID, now)
}
