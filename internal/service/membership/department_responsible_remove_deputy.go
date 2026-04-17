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

	"github.com/medincident/medincident-command-service/internal/model"
)

// RemoveDepartmentResponsibleDeputyCommand carries the identifiers
// needed to clear the deputy slot on a DR role.
type RemoveDepartmentResponsibleDeputyCommand struct {
	DepartmentID uuid.UUID
	EmployeeID   uuid.UUID
}

// RemoveDepartmentResponsibleDeputy clears the deputy slot. Fails if
// the slot is already empty (no idempotent no-op per spec §4.7).
func (s *EmployeeService) RemoveDepartmentResponsibleDeputy(ctx context.Context, cmd RemoveDepartmentResponsibleDeputyCommand) error {
	var errs []error
	if cmd.DepartmentID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentIDEmpty).
			Public("Department ID is required.").Errorf("department id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeDepartmentResponsible).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeDepartmentResponsible).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("department_id", cmd.DepartmentID).
				With("employee_id", cmd.EmployeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}

		return publishDepartmentResponsibleDeputyRemoved(tx, cmd.DepartmentID, cmd.EmployeeID, now)
	})
}
