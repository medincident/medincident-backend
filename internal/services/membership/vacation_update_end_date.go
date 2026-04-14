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
			Errorf(ErrCodeVacationIDEmpty))
	}
	if cmd.EndsAt.IsZero() {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationEndBeforeStart).
			Public("Vacation end is required.").
			Errorf(ErrCodeVacationEndBeforeStart))
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
					Errorf(ErrCodeVacationNotFound)
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
				Errorf(ErrCodeVacationAlreadyEnded)
		}

		// "New end in future."
		if !cmd.EndsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndInPast).
				Public("New end date must be in the future.").
				Errorf(ErrCodeVacationEndInPast)
		}

		// "New end after start."
		if !cmd.EndsAt.After(vac.StartsAt) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndBeforeStart).
				Public("New end date must be after the vacation's start.").
				Errorf(ErrCodeVacationEndBeforeStart)
		}

		// No-op: same end already set.
		if vac.EndsAt.Valid && vac.EndsAt.Time.Equal(cmd.EndsAt) {
			return nil
		}

		vac.EndsAt = null.TimeFrom(cmd.EndsAt)
		if err := tx.Save(&vac).Error; err != nil {
			return mapVacationInsertError(err, vac.EmployeeID)
		}

		ev := &employeev1.VacationEndDateChanged{
			VacationId: vac.ID.String(),
			EndsAt:     timestamppb.New(cmd.EndsAt),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(now),
			AggregateType: AggregateTypeEmployee,
			AggregateId:   vac.EmployeeID.String(),
			Payload:       payload,
		}
		return AppendMembershipOutboxEvent(tx, SubjectVacationEndDateChanged, envelope, nil)
	})
}
