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

	clinicv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/clinic/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

// AssignClinicHeadDeputyCommand carries the identifiers needed to set
// the deputy slot on an existing CH role.
type AssignClinicHeadDeputyCommand struct {
	ClinicID         uuid.UUID
	EmployeeID       uuid.UUID
	DeputyEmployeeID uuid.UUID
}

// AssignClinicHeadDeputy sets the deputy slot on an existing CH role.
// See spec §8.5 (ClinicHead variant).
func (s *EmployeeService) AssignClinicHeadDeputy(ctx context.Context, cmd AssignClinicHeadDeputyCommand) error {
	var errs []error
	if cmd.ClinicID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).Code(ErrCodeClinicIDEmpty).
			Public("Clinic ID is required.").Errorf("clinic id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if cmd.DeputyEmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).Code(ErrCodeDeputyEmployeeIDEmpty).
			Public("Deputy employee ID is required.").Errorf("deputy id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.ClinicHead
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("clinic_id = ? AND employee_id = ?", cmd.ClinicID, cmd.EmployeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeClinicHeadNotFound).
					Public("Clinic head not found.").
					With("clinic_id", cmd.ClinicID).
					With("employee_id", cmd.EmployeeID).
					Errorf("not found")
			}
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", cmd.DeputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", cmd.DeputyEmployeeID).
					Errorf("deputy not found")
			}
			return oops.In(scopeClinicHead).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: does the deputy's current department belong to the
		// same clinic as the role? JOIN through departments since there is no
		// denormalized clinic_id on employees.
		var deputyInClinic int64
		if err := tx.Raw(
			`SELECT count(*) FROM domain.departments WHERE id = ? AND clinic_id = ?`,
			deputy.DepartmentID, row.ClinicID,
		).Scan(&deputyInClinic).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicLookupFailed).Wrap(err)
		}
		if deputyInClinic == 0 {
			return oops.In(scopeClinicHead).
				Code(ErrCodeDeputyNotInClinic).
				Public("Deputy employee does not belong to this clinic.").
				With("deputy_department_id", deputy.DepartmentID).
				With("target_clinic_id", row.ClinicID).
				Errorf("deputy not in clinic")
		}

		if cmd.DeputyEmployeeID == cmd.EmployeeID {
			return oops.In(scopeClinicHead).
				Code(ErrCodeDeputyIsHolder).
				Public("Deputy cannot be the same employee as the role holder.").
				Errorf("deputy is holder")
		}

		if row.DeputyEmployeeID.Valid {
			return oops.In(scopeClinicHead).
				Code(ErrCodeDeputyAlreadyAssigned).
				Public("Deputy slot is already occupied.").
				With("current_deputy_employee_id", row.DeputyEmployeeID.V).
				Errorf("deputy already assigned")
		}

		row.DeputyEmployeeID = null.ValueFrom(cmd.DeputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		ev := &clinicv1.ClinicHeadDeputyAssigned{
			EmployeeId:       cmd.EmployeeID.String(),
			DeputyEmployeeId: cmd.DeputyEmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(time.Now().UTC()),
			AggregateType: AggregateTypeClinic,
			AggregateId:   cmd.ClinicID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectClinicHeadDeputyAssigned, envelope)
	})
}
