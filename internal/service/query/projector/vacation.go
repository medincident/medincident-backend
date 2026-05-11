package projector

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	vacv1 "github.com/medincident/medincident-backend/pkg/event/vacation/v1"
)

const (
	vacationStateScheduled = "scheduled"
	vacationStateActive    = "active"
	vacationStateEnded     = "ended"
	vacationStateCancelled = "cancelled"
)

func VacationScheduled(tx *gorm.DB, aggregateID string, _ time.Time, ev *vacv1.VacationScheduled) error {
	empID := uuid.MustParse(aggregateID)
	vacID := uuid.MustParse(ev.GetVacationId())
	createdAt := ev.GetCreatedAt().AsTime()
	var endsAt null.Time
	if ts := ev.GetEndsAt(); ts != nil && ts.IsValid() {
		endsAt = null.TimeFrom(ts.AsTime())
	}

	if err := tx.Exec(`INSERT INTO projections.employee_vacations (id, employee_id, state, starts_at, ends_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
		vacID, empID, vacationStateScheduled, ev.GetStartsAt().AsTime(), endsAt, createdAt, createdAt).Error; err != nil {
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("vacation_id", ev.GetVacationId()).Wrap(err)
	}
	return refreshCardVacations(tx, empID, createdAt)
}

func VacationStarted(tx *gorm.DB, aggregateID string, _ time.Time, ev *vacv1.VacationStarted) error {
	empID := uuid.MustParse(aggregateID)
	vacID := uuid.MustParse(ev.GetVacationId())
	createdAt := ev.GetCreatedAt().AsTime()
	var endsAt null.Time
	if ts := ev.GetEndsAt(); ts != nil && ts.IsValid() {
		endsAt = null.TimeFrom(ts.AsTime())
	}

	if err := tx.Exec(`INSERT INTO projections.employee_vacations (id, employee_id, state, starts_at, ends_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
		vacID, empID, vacationStateActive, ev.GetStartsAt().AsTime(), endsAt, createdAt, createdAt).Error; err != nil {
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("vacation_id", ev.GetVacationId()).Wrap(err)
	}
	return refreshCardVacations(tx, empID, createdAt)
}

func VacationEnded(tx *gorm.DB, aggregateID string, _ time.Time, ev *vacv1.VacationEnded) error {
	empID := uuid.MustParse(aggregateID)
	vacID := uuid.MustParse(ev.GetVacationId())
	endsAt := ev.GetEndsAt().AsTime()

	if err := tx.Exec(`UPDATE projections.employee_vacations SET state = ?, ends_at = ?, updated_at = ? WHERE id = ?`,
		vacationStateEnded, endsAt, endsAt, vacID).Error; err != nil {
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("vacation_id", ev.GetVacationId()).Wrap(err)
	}
	return refreshCardVacations(tx, empID, endsAt)
}

func VacationCancelled(tx *gorm.DB, aggregateID string, _ time.Time, ev *vacv1.VacationCancelled) error {
	empID := uuid.MustParse(aggregateID)
	vacID := uuid.MustParse(ev.GetVacationId())
	cancelledAt := ev.GetCancelledAt().AsTime()

	if err := tx.Exec(`UPDATE projections.employee_vacations SET state = ?, updated_at = ? WHERE id = ?`,
		vacationStateCancelled, cancelledAt, vacID).Error; err != nil {
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("vacation_id", ev.GetVacationId()).Wrap(err)
	}
	return refreshCardVacations(tx, empID, cancelledAt)
}

func VacationEndDateChanged(tx *gorm.DB, aggregateID string, _ time.Time, ev *vacv1.VacationEndDateChanged) error {
	empID := uuid.MustParse(aggregateID)
	vacID := uuid.MustParse(ev.GetVacationId())
	updatedAt := ev.GetUpdatedAt().AsTime()
	var endsAt null.Time
	if ts := ev.GetEndsAt(); ts != nil && ts.IsValid() {
		endsAt = null.TimeFrom(ts.AsTime())
	}

	if err := tx.Exec(`UPDATE projections.employee_vacations SET ends_at = ?, updated_at = ? WHERE id = ?`,
		endsAt, updatedAt, vacID).Error; err != nil {
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("vacation_id", ev.GetVacationId()).Wrap(err)
	}
	return refreshCardVacations(tx, empID, updatedAt)
}

// refreshCardVacations re-reads current+next vacation from projections.employee_vacations
// and updates the employee_card pointer columns.
func refreshCardVacations(tx *gorm.DB, empID uuid.UUID, now time.Time) error {
	var activeID uuid.UUID
	var activeEndsAt null.Time
	activeErr := tx.Raw(`SELECT id, ends_at FROM projections.employee_vacations WHERE employee_id = ? AND state = ? ORDER BY starts_at DESC LIMIT 1`,
		empID, vacationStateActive).Row().Scan(&activeID, &activeEndsAt)

	var nextID uuid.UUID
	var nextStartsAt time.Time
	nextErr := tx.Raw(`SELECT id, starts_at FROM projections.employee_vacations WHERE employee_id = ? AND state = ? AND starts_at > ? ORDER BY starts_at ASC LIMIT 1`,
		empID, vacationStateScheduled, now).Row().Scan(&nextID, &nextStartsAt)

	var (
		currentVacationID     *uuid.UUID
		currentVacationEndsAt null.Time
		nextVacationID        *uuid.UUID
		nextVacationStartsAt  null.Time
	)
	switch {
	case activeErr == nil:
		id := activeID
		currentVacationID = &id
		currentVacationEndsAt = activeEndsAt
	case errors.Is(activeErr, sql.ErrNoRows):
	default:
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("employee_id", empID).Wrap(activeErr)
	}
	switch {
	case nextErr == nil:
		id := nextID
		nextVacationID = &id
		nextVacationStartsAt = null.TimeFrom(nextStartsAt)
	case errors.Is(nextErr, sql.ErrNoRows):
	default:
		return oops.In("projector.vacation").Code(ErrCodeVacationProjectionFailed).With("employee_id", empID).Wrap(nextErr)
	}

	return tx.Exec(`UPDATE projections.employee_cards SET current_vacation_id = ?, current_vacation_ends_at = ?, next_vacation_id = ?, next_vacation_starts_at = ?, updated_at = ? WHERE employee_id = ?`,
		currentVacationID, currentVacationEndsAt, nextVacationID, nextVacationStartsAt, now, empID).Error
}

// Forwarding methods on *Projectors

func (p *Projectors) VacationScheduled(tx *gorm.DB, id string, t time.Time, ev *vacv1.VacationScheduled) error {
	return VacationScheduled(tx, id, t, ev)
}

func (p *Projectors) VacationStarted(tx *gorm.DB, id string, t time.Time, ev *vacv1.VacationStarted) error {
	return VacationStarted(tx, id, t, ev)
}

func (p *Projectors) VacationEnded(tx *gorm.DB, id string, t time.Time, ev *vacv1.VacationEnded) error {
	return VacationEnded(tx, id, t, ev)
}

func (p *Projectors) VacationCancelled(tx *gorm.DB, id string, t time.Time, ev *vacv1.VacationCancelled) error {
	return VacationCancelled(tx, id, t, ev)
}

func (p *Projectors) VacationEndDateChanged(tx *gorm.DB, id string, t time.Time, ev *vacv1.VacationEndDateChanged) error {
	return VacationEndDateChanged(tx, id, t, ev)
}
