package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	employeev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/employee/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
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
	now := time.Now().UTC()

	var errs []error
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id is empty"))
	}
	if cmd.StartsAt.IsZero() {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationStartRequired).
			Public("Vacation start is required.").
			Errorf("vacation start is required"))
	}
	if !cmd.StartsAt.IsZero() && !cmd.StartsAt.After(now) {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationStartInPast).
			Public("Scheduled vacation must start in the future.").
			Errorf("vacation start is in the past"))
	}
	if cmd.EndsAt != nil && !cmd.EndsAt.After(cmd.StartsAt) {
		errs = append(errs, oops.In(scopeVacation).
			Code(ErrCodeVacationEndBeforeStart).
			Public("Vacation end must be after its start.").
			Errorf("vacation end must be after start"))
	}
	if len(errs) > 0 {
		return ScheduleVacationResult{}, errors.Join(errs...)
	}

	var result ScheduleVacationResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
			mapped := mapVacationInsertError(err, cmd.EmployeeID)
			return mapped
		}

		ev := &employeev1.VacationScheduled{
			VacationId: id.String(),
			StartsAt:   timestamppb.New(cmd.StartsAt),
		}
		if cmd.EndsAt != nil {
			ev.EndsAt = timestamppb.New(*cmd.EndsAt)
		}
		if err := outbox.Publish(tx, SubjectVacationScheduled, AggregateTypeEmployee, cmd.EmployeeID.String(), now, ev); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
