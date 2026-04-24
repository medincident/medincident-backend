package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// RemoveDepartmentResponsibleDeputyPayload carries the identifiers
// needed to clear the deputy slot on a DR role.
type RemoveDepartmentResponsibleDeputyPayload struct {
	DepartmentID string `validate:"required,uuid"`
	EmployeeID   string `validate:"required,uuid"`
}

// RemoveDepartmentResponsibleDeputyCommand = caller + payload.
type RemoveDepartmentResponsibleDeputyCommand struct {
	Caller  authz.Caller
	Payload RemoveDepartmentResponsibleDeputyPayload
}

// RemoveDepartmentResponsibleDeputy clears the deputy slot. Fails if
// the slot is already empty (no idempotent no-op per spec §4.7).
func (s *EmployeeService) RemoveDepartmentResponsibleDeputy(ctx context.Context, cmd RemoveDepartmentResponsibleDeputyCommand) error {
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

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("department_id", departmentID).
				With("employee_id", employeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		return publishDepartmentResponsibleDeputyRemoved(tx, departmentID, employeeID, now)
	})
}
