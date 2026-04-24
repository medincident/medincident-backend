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
	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// CancelScheduledVacationPayload identifies the future vacation to remove.
type CancelScheduledVacationPayload struct {
	VacationID string `validate:"required,uuid"`
}

// CancelScheduledVacationCommand = caller + payload.
type CancelScheduledVacationCommand struct {
	Caller  authz.Caller
	Payload CancelScheduledVacationPayload
}

// CancelScheduledVacation removes a not-yet-started vacation. For a
// vacation whose starts_at has already passed, use ForceEndVacation
// instead.
func (s *EmployeeService) CancelScheduledVacation(ctx context.Context, cmd CancelScheduledVacationCommand) error {
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

		if !vac.StartsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationAlreadyStarted).
				Public("Vacation has already started — force-end it instead.").
				Hint("Use ForceEndVacation for running vacations.").
				Errorf("vacation has already started")
		}

		if err := tx.Delete(&model.EmployeeVacation{}, "id = ?", vacationID).Error; err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationDeleteFailed).Wrap(err)
		}

		return projector.VacationCancelled(tx, &vac, now)
	})
}
