package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// CancelScheduledVacationCommand identifies the future vacation to remove.
type CancelScheduledVacationCommand struct {
	VacationID uuid.UUID `validate:"required"`
}

// CancelScheduledVacation removes a not-yet-started vacation. For a
// vacation whose starts_at has already passed, use ForceEndVacation
// instead.
func (s *EmployeeService) CancelScheduledVacation(ctx context.Context, cmd CancelScheduledVacationCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var vac model.EmployeeVacation
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", cmd.VacationID).
			First(&vac).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeVacation).
					Code(ErrCodeVacationNotFound).
					Public("Vacation not found.").
					With("vacation_id", cmd.VacationID).
					Errorf("vacation not found")
			}
			return oops.In(scopeVacation).Code(ErrCodeVacationLoadFailed).Wrap(err)
		}

		if !vac.StartsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationAlreadyStarted).
				Public("Vacation has already started — force-end it instead.").
				Hint("Use ForceEndVacation for running vacations.").
				Errorf("vacation has already started")
		}

		if err := tx.Delete(&model.EmployeeVacation{}, "id = ?", cmd.VacationID).Error; err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationDeleteFailed).Wrap(err)
		}

		return projector.VacationCancelled(tx, &vac, now)
	})
}
