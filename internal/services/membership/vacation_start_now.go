package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

const scopeVacation = "services.membership.vacation"

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
			Errorf(ErrCodeEmployeeIDEmpty)
	}

	var result StartVacationNowResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		// Pre-check employee exists for a clean error code.
		var exists int64
		if err := tx.Raw(`SELECT count(*) FROM domain.employees WHERE id = ?`, cmd.EmployeeID).Scan(&exists).Error; err != nil {
			return oops.In(scopeVacation).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}
		if exists == 0 {
			return oops.In(scopeVacation).
				Code(ErrCodeEmployeeNotFound).
				Public("Employee not found.").
				With("employee_id", cmd.EmployeeID).
				Errorf(ErrCodeEmployeeNotFound)
		}

		if cmd.EndsAt != nil && !cmd.EndsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationEndBeforeStart).
				Public("Vacation end must be after its start.").
				Errorf(ErrCodeVacationEndBeforeStart)
		}

		id, err := uuid.NewV7()
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationIDGenerationFailed).Wrap(err)
		}
		vac := model.EmployeeVacation{
			ID:         id,
			EmployeeID: cmd.EmployeeID,
			StartsAt:   now,
			EndsAt:     nullTimeFromPtr(cmd.EndsAt),
		}
		if err := tx.Create(&vac).Error; err != nil {
			return mapVacationInsertError(err)
		}

		ev := &employeev1.VacationStarted{
			VacationId: id.String(),
			StartsAt:   timestamppb.New(now),
		}
		if cmd.EndsAt != nil {
			ev.EndsAt = timestamppb.New(*cmd.EndsAt)
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(now),
			AggregateType: AggregateTypeEmployee,
			AggregateId:   cmd.EmployeeID.String(),
			Payload:       payload,
		}
		if err := AppendMembershipOutboxEvent(tx, SubjectVacationStarted, envelope, nil); err != nil {
			return err
		}

		result.ID = id
		return nil
	})
	return result, err
}

// nullTimeFromPtr is a shared helper used by vacation commands.
func nullTimeFromPtr(t *time.Time) null.Time {
	if t == nil {
		return null.Time{}
	}
	return null.TimeFrom(*t)
}

// mapVacationInsertError translates Postgres constraint errors on
// employee_vacations INSERT/UPDATE into domain error codes. Used by
// all vacation write operations.
func mapVacationInsertError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgErrCodeExclusionViolation:
			return oops.In(scopeVacation).
				Code(ErrCodeVacationOverlap).
				Public("Vacation overlaps with an existing one.").
				Wrap(err)
		case pgErrCodeForeignKeyViolation:
			return oops.In(scopeVacation).
				Code(ErrCodeEmployeeNotFound).
				Public("Employee not found.").
				Wrap(err)
		}
	}
	return oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).Wrap(err)
}
