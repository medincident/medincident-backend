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

// ForceEndVacationCommand identifies the running vacation to close.
type ForceEndVacationCommand struct {
	VacationID uuid.UUID
}

// ForceEndVacation closes a running vacation at the current moment.
// Only applicable to vacations whose starts_at has already passed.
// For scheduled (future) vacations, use CancelScheduledVacation.
func (s *EmployeeService) ForceEndVacation(ctx context.Context, cmd ForceEndVacationCommand) error {
	if cmd.VacationID == uuid.Nil {
		return oops.In(scopeVacation).
			Code(ErrCodeVacationIDEmpty).
			Public("Vacation ID is required.").
			Errorf("vacation id is empty")
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

		ev := &employeev1.VacationEnded{
			VacationId: vac.ID.String(),
			EndsAt:     timestamppb.New(now),
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
		return outbox.AppendEvent(tx, SubjectVacationEnded, envelope, nil)
	})
}
