package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/outbox"
	employeev1 "github.com/medincident/medincident-command-service/pkg/event/employee/v1"
)

// StartVacationNowCommand carries everything the service needs to
// start a vacation at the current time.
type StartVacationNowCommand struct {
	EmployeeID uuid.UUID
	EndsAt     *time.Time
}

// StartVacationNowResult holds the ID of the newly created vacation.
type StartVacationNowResult struct {
	ID uuid.UUID
}

// StartVacationNow starts a vacation at the current time. EndsAt may
// be nil (unlimited). Overlap with existing vacations of the same
// employee is rejected by the exclusion constraint on the table.
func (s *EmployeeService) StartVacationNow(ctx context.Context, cmd StartVacationNowCommand) (StartVacationNowResult, error) {
	if cmd.EmployeeID == uuid.Nil {
		return StartVacationNowResult{}, oops.In(scopeVacation).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id is empty")
	}

	var result StartVacationNowResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		if cmd.EndsAt != nil && !cmd.EndsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndBeforeStart).
				Public("Vacation end must be after its start.").
				Errorf("vacation end must be after start")
		}

		id, err := uuid.NewV7()
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationIDGenerationFailed).Wrap(err)
		}
		vac := model.EmployeeVacation{
			ID:         id,
			EmployeeID: cmd.EmployeeID,
			StartsAt:   now,
			EndsAt:     null.TimeFromPtr(cmd.EndsAt),
		}
		if err := tx.Create(&vac).Error; err != nil {
			return mapVacationInsertError(err, cmd.EmployeeID)
		}

		ev := &employeev1.VacationStarted{
			VacationId: id.String(),
			StartsAt:   timestamppb.New(now),
		}
		if cmd.EndsAt != nil {
			ev.EndsAt = timestamppb.New(*cmd.EndsAt)
		}
		if err := outbox.Publish(tx, SubjectVacationStarted, AggregateTypeEmployee, cmd.EmployeeID.String(), now, ev); err != nil {
			return err
		}

		result.ID = id
		return nil
	})
	return result, err
}

// mapVacationInsertError translates Postgres constraint errors on
// employee_vacations INSERT/UPDATE into domain error codes. Used by
// all vacation write operations. employeeID is attached to the
// NotFound branch so operators can identify which employee was missing.
// isExclusionViolation reports whether err wraps a Postgres exclusion
// violation (23P01). GORM does not translate this code, so we check
// the underlying pgconn.PgError directly.
func isExclusionViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23P01"
}

func mapVacationInsertError(err error, employeeID uuid.UUID) error {
	if isExclusionViolation(err) {
		return oops.In(scopeVacation).
			Code(ErrCodeVacationOverlap).
			Public("Vacation overlaps with an existing one.").
			Wrap(err)
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return oops.In(scopeVacation).
			Code(ErrCodeEmployeeNotFound).
			Public("Employee not found.").
			With("employee_id", employeeID).
			Wrap(err)
	}
	return oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).Wrap(err)
}
