package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// UpdateVacationEndDateCommand carries the vacation to update and the
// desired new end date.
type UpdateVacationEndDateCommand struct {
	VacationID uuid.UUID
	EndsAt     time.Time
}

// UpdateVacationEndDate changes the end date of a not-yet-ended
// vacation. Supports setting (on unlimited), extending, or shortening.
// Uses SELECT ... FOR UPDATE on the vacation row to serialise against
// a concurrent ForceEndVacation (which would otherwise allow this
// command to resurrect a manually-ended vacation from a stale read).
func (s *EmployeeService) UpdateVacationEndDate(ctx context.Context, cmd UpdateVacationEndDateCommand) error {
	var errs []error
	if cmd.VacationID == uuid.Nil {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationIDEmpty).
			Public("Vacation ID is required.").
			Errorf("vacation id is empty"))
	}
	if cmd.EndsAt.IsZero() {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationEndBeforeStart).
			Public("Vacation end is required.").
			Errorf("vacation end must be after start"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
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

		// "Not already ended": unlimited vacation (EndsAt invalid), or
		// currently ends in the future.
		if vac.EndsAt.Valid && !vac.EndsAt.Time.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationAlreadyEnded).
				Public("Vacation has already ended.").
				With("vacation_id", cmd.VacationID).
				Errorf("vacation has already ended")
		}

		// "New end in future."
		if !cmd.EndsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndInPast).
				Public("New end date must be in the future.").
				Errorf("vacation end is in the past")
		}

		// "New end after start."
		if !cmd.EndsAt.After(vac.StartsAt) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndBeforeStart).
				Public("New end date must be after the vacation's start.").
				Errorf("vacation end must be after start")
		}

		// No-op: same end already set. Compare at microsecond resolution
		// — Postgres TIMESTAMPTZ stores microseconds, so a caller that
		// passes a Go time with sub-microsecond nanoseconds would
		// otherwise see spurious updates after the value round-trips
		// through the database.
		if vac.EndsAt.Valid && vac.EndsAt.Time.Equal(cmd.EndsAt.Truncate(time.Microsecond)) {
			return nil
		}

		vac.EndsAt = null.TimeFrom(cmd.EndsAt)
		if saveErr := tx.Save(&vac).Error; saveErr != nil {
			return oops.In(scopeVacation).With("vacation_id", vac.ID).Wrap(mapVacationInsertError(saveErr, vac.EmployeeID))
		}

		ev := &employeev1.VacationEndDateChanged{
			VacationId: vac.ID.String(),
			EndsAt:     timestamppb.New(cmd.EndsAt),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(now),
			AggregateType: AggregateTypeEmployee,
			AggregateId:   vac.EmployeeID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectVacationEndDateChanged, envelope, nil)
	})
}
