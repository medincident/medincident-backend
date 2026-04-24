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
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// ForceEndVacationPayload identifies the running vacation to close.
type ForceEndVacationPayload struct {
	VacationID string `validate:"required,uuid"`
}

// ForceEndVacationCommand = caller + payload.
type ForceEndVacationCommand struct {
	Caller  authz.Caller
	Payload ForceEndVacationPayload
}

// ForceEndVacation closes a running vacation at the current moment.
// Only applicable to vacations whose starts_at has already passed.
// For scheduled (future) vacations, use CancelScheduledVacation.
func (s *EmployeeService) ForceEndVacation(ctx context.Context, cmd ForceEndVacationCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	vacationID := uuid.MustParse(cmd.Payload.VacationID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Vacation(vacationID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		var vac model.EmployeeVacation
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", vacationID).
			First(&vac).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeVacation).
					Code(ErrCodeVacationNotFound).
					Public("Vacation not found.").
					With("vacation_id", vacationID).
					Errorf("vacation not found")
			}
			return oops.In(scopeVacation).Code(ErrCodeVacationLoadFailed).Wrap(err)
		}

		if vac.StartsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationNotStarted).
				Public("Vacation has not started yet — cancel it instead.").
				Hint("Use CancelScheduledVacation for scheduled vacations.").
				Errorf("vacation has not started yet")
		}
		if vac.EndsAt.Valid && !vac.EndsAt.Time.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationAlreadyEnded).
				Public("Vacation has already ended.").
				Errorf("vacation has already ended")
		}

		vac.EndsAt = null.TimeFrom(now)
		if err := tx.Save(&vac).Error; err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).With("vacation_id", vac.ID).Wrap(err)
		}

		return projector.VacationEnded(tx, &vac, now)
	})
}
