package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
	vacv1 "github.com/medincident/medincident-backend/pkg/event/vacation/v1"
)

// StartVacationNowPayload carries everything the service needs to
// start a vacation at the current time.
type StartVacationNowPayload struct {
	EmployeeID string     `validate:"required,uuid"`
	EndsAt     *time.Time `validate:"omitnil"`
}

// StartVacationNowCommand = caller + payload.
type StartVacationNowCommand struct {
	Caller  authz.Caller
	Payload StartVacationNowPayload
}

// StartVacationNowResult holds the ID of the newly created vacation.
type StartVacationNowResult struct {
	ID uuid.UUID `validate:"required"`
}

// StartVacationNow starts a vacation at the current time. EndsAt may
// be nil (unlimited). Overlap with existing vacations of the same
// employee is rejected by the exclusion constraint on the table.
//
// See: docs/services/Membership.md
func (s *EmployeeService) StartVacationNow(ctx context.Context, cmd StartVacationNowCommand) (StartVacationNowResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return StartVacationNowResult{}, err
	}
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Employee(employeeID)); err != nil {
		return StartVacationNowResult{}, err
	}

	var result StartVacationNowResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		if cmd.Payload.EndsAt != nil && !cmd.Payload.EndsAt.After(now) {
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
			EmployeeID: employeeID,
			StartsAt:   now,
			EndsAt:     null.TimeFromPtr(cmd.Payload.EndsAt),
		}
		if err := tx.Create(&vac).Error; err != nil {
			return mapVacationInsertError(err, employeeID)
		}

		env, err := buildVacationStartedEnvelope(&vac)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.vacation.v1.started", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}

func buildVacationStartedEnvelope(v *model.EmployeeVacation) (*eventv1.Envelope, error) {
	msg := &vacv1.VacationStarted{
		VacationId: v.ID.String(),
		StartsAt:   timestamppb.New(v.StartsAt),
		CreatedAt:  timestamppb.New(v.CreatedAt),
	}
	if v.EndsAt.Valid {
		msg.EndsAt = timestamppb.New(v.EndsAt.Time)
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scopeVacation).Code(ErrCodeVacationSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(v.StartsAt),
		AggregateType: "vacation",
		AggregateId:   v.EmployeeID.String(),
		Payload:       payload,
	}, nil
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
