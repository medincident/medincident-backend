package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// ScheduleVacationCommand carries everything the service needs to
// schedule a future-dated vacation.
type ScheduleVacationCommand struct {
	EmployeeID uuid.UUID `validate:"required"`
	StartsAt   time.Time
	EndsAt     *time.Time
}

// ScheduleVacationResult holds the ID of the newly scheduled vacation.
type ScheduleVacationResult struct {
	ID uuid.UUID `validate:"required"`
}

// ScheduleVacation creates a future-dated vacation. StartsAt must be
// strictly in the future; EndsAt may be nil (unlimited) or must be
// strictly after StartsAt. Overlap with existing vacations of the same
// employee is rejected by the exclusion constraint on the table.
func (s *EmployeeService) ScheduleVacation(ctx context.Context, cmd ScheduleVacationCommand) (ScheduleVacationResult, error) {
	if err := validation.Struct(cmd); err != nil {
		return ScheduleVacationResult{}, err
	}

	now := time.Now().UTC()
	var errs []error
	if !cmd.StartsAt.After(now) {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationStartInPast).
			Public("Scheduled vacation must start in the future.").
			Errorf("vacation start is in the past"))
	}
	if cmd.EndsAt != nil && !cmd.EndsAt.After(cmd.StartsAt) {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationEndBeforeStart).
			Public("Vacation end must be after its start.").
			Errorf("vacation end must be after start"))
	}
	if len(errs) > 0 {
		return ScheduleVacationResult{}, errors.Join(errs...)
	}

	var result ScheduleVacationResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		id, err := uuid.NewV7()
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationIDGenerationFailed).Wrap(err)
		}
		vac := model.EmployeeVacation{
			ID:         id,
			EmployeeID: cmd.EmployeeID,
			StartsAt:   cmd.StartsAt,
			EndsAt:     null.TimeFromPtr(cmd.EndsAt),
		}
		if err := tx.Create(&vac).Error; err != nil {
			mapped := mapVacationInsertError(err, cmd.EmployeeID)
			return mapped
		}

		if err := projector.VacationScheduled(tx, &vac); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
