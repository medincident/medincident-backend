package membership

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
)

// UpdateEmployeePositionCommand carries the inputs required to change an
// employee's position. Position nil means "clear the position".
type UpdateEmployeePositionCommand struct {
	ID       uuid.UUID
	Position *string
}

// UpdatePosition changes an employee's position. Idempotent: if the
// normalised new position equals the current one, no event is emitted.
// Takes a FOR UPDATE row lock on the employee to serialise concurrent
// updates; last-writer-wins is NOT acceptable because it would emit
// events describing overwritten intermediate states.
func (s *EmployeeService) UpdatePosition(ctx context.Context, cmd UpdateEmployeePositionCommand) error {
	var errs []error
	if cmd.ID == uuid.Nil {
		errs = append(errs, oops.In(scopeEmployee).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id is empty"))
	}
	if err := validatePosition(cmd.Position); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	var newPos null.String
	if cmd.Position != nil {
		newPos = null.StringFrom(strings.TrimSpace(*cmd.Position))
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

		if emp.Position == newPos {
			return nil
		}
		emp.Position = newPos
		if err := tx.Save(&emp).Error; err != nil {
			return oops.In(scopeEmployee).Code(ErrCodeEmployeeSaveFailed).Wrap(err)
		}

		return projector.EmployeePositionChanged(tx, &emp)
	})
}
