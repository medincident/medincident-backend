package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

// ScheduleVacationCommand carries everything the service needs to
// schedule a future-dated vacation.
type ScheduleVacationCommand struct {
	EmployeeID uuid.UUID
	StartsAt   time.Time
	EndsAt     *time.Time
}

// ScheduleVacationResult holds the ID of the newly scheduled vacation.
type ScheduleVacationResult struct {
	ID uuid.UUID
}

// ScheduleVacation creates a future-dated vacation. StartsAt must be
// strictly in the future; EndsAt may be nil (unlimited) or must be
// strictly after StartsAt. Overlap with existing vacations of the same
// employee is rejected by the exclusion constraint on the table.
func (s *EmployeeService) ScheduleVacation(ctx context.Context, cmd ScheduleVacationCommand) (ScheduleVacationResult, error) {
	var errs []error
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf(ErrCodeEmployeeIDEmpty))
	}
	if cmd.EndsAt != nil && !cmd.EndsAt.After(cmd.StartsAt) {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationEndBeforeStart).
			Public("Vacation end must be after its start.").
			Errorf(ErrCodeVacationEndBeforeStart))
	}
	if len(errs) > 0 {
		return ScheduleVacationResult{}, errors.Join(errs...)
	}

	var result ScheduleVacationResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()
		if !cmd.StartsAt.After(now) {
			return oops.In(scopeVacation).
				Code(ErrCodeVacationStartInPast).
				Public("Scheduled vacation must start in the future.").
				Errorf(ErrCodeVacationStartInPast)
		}

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

		id, err := uuid.NewV7()
		if err != nil {
			return oops.In(scopeVacation).Code(ErrCodeVacationIDGenerationFailed).Wrap(err)
		}
		vac := model.EmployeeVacation{
			ID:         id,
			EmployeeID: cmd.EmployeeID,
			StartsAt:   cmd.StartsAt,
			EndsAt:     nullTimeFromPtr(cmd.EndsAt),
		}
		if err := tx.Create(&vac).Error; err != nil {
			return mapVacationInsertError(err)
		}

		ev := &employeev1.VacationScheduled{
			VacationId: id.String(),
			StartsAt:   timestamppb.New(cmd.StartsAt),
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
		if err := AppendMembershipOutboxEvent(tx, SubjectVacationScheduled, envelope, nil); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
