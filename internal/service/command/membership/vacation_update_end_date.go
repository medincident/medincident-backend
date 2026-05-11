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

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
	vacv1 "github.com/medincident/medincident-backend/pkg/event/vacation/v1"
)

// UpdateVacationEndDatePayload carries the vacation to update and the
// desired new end date.
type UpdateVacationEndDatePayload struct {
	VacationID string    `validate:"required,uuid"`
	EndsAt     time.Time `validate:"required"`
}

// UpdateVacationEndDateCommand = caller + payload.
type UpdateVacationEndDateCommand struct {
	Caller  authz.Caller
	Payload UpdateVacationEndDatePayload
}

// UpdateVacationEndDate changes the end date of a not-yet-ended
// vacation. Supports setting (on unlimited), extending, or shortening.
// Uses SELECT ... FOR UPDATE on the vacation row to serialise against
// a concurrent ForceEndVacation (which would otherwise allow this
// command to resurrect a manually-ended vacation from a stale read).
//
// See: docs/services/Membership.md
func (s *EmployeeService) UpdateVacationEndDate(ctx context.Context, cmd UpdateVacationEndDateCommand) error {
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

		// "Not already ended": unlimited vacation (EndsAt invalid), or
		// currently ends in the future.
		if vac.EndsAt.Valid && !vac.EndsAt.Time.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationAlreadyEnded).
				Public("Vacation has already ended.").
				With("vacation_id", vacationID).
				Errorf("vacation has already ended")
		}

		// "New end in future."
		if !cmd.Payload.EndsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndInPast).
				Public("New end date must be in the future.").
				Errorf("vacation end is in the past")
		}

		// "New end after start."
		if !cmd.Payload.EndsAt.After(vac.StartsAt) {
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
		if vac.EndsAt.Valid && vac.EndsAt.Time.Equal(cmd.Payload.EndsAt.Truncate(time.Microsecond)) {
			return nil
		}

		vac.EndsAt = null.TimeFrom(cmd.Payload.EndsAt)
		if saveErr := tx.Save(&vac).Error; saveErr != nil {
			return oops.In(scopeVacation).With("vacation_id", vac.ID).Wrap(mapVacationInsertError(saveErr, vac.EmployeeID))
		}

		env, err := buildVacationEndDateChangedEnvelope(&vac)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.vacation.v1.end_date_changed", env)
	})
}

func buildVacationEndDateChangedEnvelope(v *model.EmployeeVacation) (*eventv1.Envelope, error) {
	msg := &vacv1.VacationEndDateChanged{
		VacationId: v.ID.String(),
		UpdatedAt:  timestamppb.New(v.UpdatedAt),
	}
	if v.EndsAt.Valid {
		msg.EndsAt = timestamppb.New(v.EndsAt.Time)
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(v.UpdatedAt),
		AggregateType: "vacation",
		AggregateId:   v.EmployeeID.String(),
		Payload:       payload,
	}, nil
}
