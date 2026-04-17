package projector

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// Vacation lifecycle states as stored in projections.employee_vacations.
// Mirrors the query-service enum the old projector populated.
const (
	vacationStateScheduled = "scheduled"
	vacationStateActive    = "active"
	vacationStateEnded     = "ended"
	vacationStateCancelled = "cancelled"
)

// VacationScheduled inserts a new row in projections.employee_vacations
// with state='scheduled', then refreshes the employee_card current /
// next vacation pointers.
func VacationScheduled(tx *gorm.DB, v *model.EmployeeVacation) error {
	if err := tx.Exec(`
		INSERT INTO projections.employee_vacations
		    (id, employee_id, state, starts_at, ends_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		v.ID, v.EmployeeID, vacationStateScheduled,
		v.StartsAt, v.EndsAt, v.CreatedAt, v.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("vacation_id", v.ID).
			Wrap(err)
	}
	return refreshCardVacations(tx, v.EmployeeID, v.UpdatedAt)
}

// VacationStarted inserts a new row in projections.employee_vacations
// with state='active' (the command service's StartNow path always
// creates a brand-new vacation — there is no prior scheduled row to
// transition).
func VacationStarted(tx *gorm.DB, v *model.EmployeeVacation) error {
	if err := tx.Exec(`
		INSERT INTO projections.employee_vacations
		    (id, employee_id, state, starts_at, ends_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		v.ID, v.EmployeeID, vacationStateActive,
		v.StartsAt, v.EndsAt, v.CreatedAt, v.UpdatedAt,
	).Error; err != nil {
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("vacation_id", v.ID).
			Wrap(err)
	}
	return refreshCardVacations(tx, v.EmployeeID, v.UpdatedAt)
}

// VacationEnded marks a running vacation as ended and stamps the final
// ends_at, then refreshes the employee_card pointers.
func VacationEnded(tx *gorm.DB, v *model.EmployeeVacation, endsAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.employee_vacations
		   SET state = ?, ends_at = ?, updated_at = ?
		 WHERE id = ?`,
		vacationStateEnded, endsAt, endsAt, v.ID,
	).Error; err != nil {
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("vacation_id", v.ID).
			Wrap(err)
	}
	return refreshCardVacations(tx, v.EmployeeID, endsAt)
}

// VacationCancelled marks a scheduled vacation as cancelled, then
// refreshes the employee_card pointers. The domain row is deleted by
// the service, but the projection preserves cancelled rows for audit.
func VacationCancelled(tx *gorm.DB, v *model.EmployeeVacation, cancelledAt time.Time) error {
	if err := tx.Exec(`
		UPDATE projections.employee_vacations
		   SET state = ?, updated_at = ?
		 WHERE id = ?`,
		vacationStateCancelled, cancelledAt, v.ID,
	).Error; err != nil {
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("vacation_id", v.ID).
			Wrap(err)
	}
	return refreshCardVacations(tx, v.EmployeeID, cancelledAt)
}

// VacationEndDateChanged updates ends_at on an existing vacation row,
// then refreshes the employee_card pointers.
func VacationEndDateChanged(tx *gorm.DB, v *model.EmployeeVacation) error {
	if err := tx.Exec(`
		UPDATE projections.employee_vacations
		   SET ends_at = ?, updated_at = ?
		 WHERE id = ?`,
		v.EndsAt, v.UpdatedAt, v.ID,
	).Error; err != nil {
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("vacation_id", v.ID).
			Wrap(err)
	}
	return refreshCardVacations(tx, v.EmployeeID, v.UpdatedAt)
}

// refreshCardVacations re-reads the employee's current-active and
// next-scheduled vacation from employee_vacations and stamps the
// derived pointer columns on the employee_card row. Called at the end
// of every vacation handler to sidestep incremental-maintenance races.
// Mirrors the query-service projector's behaviour one-for-one.
func refreshCardVacations(tx *gorm.DB, empID uuid.UUID, now time.Time) error {
	var activeID uuid.UUID
	var activeEndsAt null.Time
	activeErr := tx.Raw(`
		SELECT id, ends_at FROM projections.employee_vacations
		 WHERE employee_id = ? AND state = ?
		 ORDER BY starts_at DESC
		 LIMIT 1`,
		empID, vacationStateActive,
	).Row().Scan(&activeID, &activeEndsAt)

	var nextID uuid.UUID
	var nextStartsAt time.Time
	nextErr := tx.Raw(`
		SELECT id, starts_at FROM projections.employee_vacations
		 WHERE employee_id = ? AND state = ? AND starts_at > ?
		 ORDER BY starts_at ASC
		 LIMIT 1`,
		empID, vacationStateScheduled, now,
	).Row().Scan(&nextID, &nextStartsAt)

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
		// no active vacation — leave nil
	default:
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("employee_id", empID).
			Wrap(activeErr)
	}
	switch {
	case nextErr == nil:
		id := nextID
		nextVacationID = &id
		nextVacationStartsAt = null.TimeFrom(nextStartsAt)
	case errors.Is(nextErr, sql.ErrNoRows):
		// no next vacation — leave nil
	default:
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("employee_id", empID).
			Wrap(nextErr)
	}

	if err := tx.Exec(`
		UPDATE projections.employee_cards
		   SET current_vacation_id      = ?,
		       current_vacation_ends_at = ?,
		       next_vacation_id         = ?,
		       next_vacation_starts_at  = ?,
		       updated_at               = ?
		 WHERE employee_id = ?`,
		currentVacationID, currentVacationEndsAt,
		nextVacationID, nextVacationStartsAt,
		now, empID,
	).Error; err != nil {
		return oops.In("projector.vacation").
			Code(ErrCodeVacationProjectionFailed).
			With("employee_id", empID).
			Wrap(err)
	}
	return nil
}
