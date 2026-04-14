package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

const scopeClinicHead = "services.membership.clinic_head"

// AssignClinicHeadCommand carries the identifiers needed to link an
// employee to a clinic as its head.
type AssignClinicHeadCommand struct {
	ClinicID   uuid.UUID
	EmployeeID uuid.UUID
}

// AssignClinicHead links the employee to the clinic as its head. The
// employee must currently work in a department that belongs to the
// clinic. See spec §8.2 (ClinicHead variant).
func (s *EmployeeService) AssignClinicHead(ctx context.Context, cmd AssignClinicHeadCommand) error {
	var errs []error
	if cmd.ClinicID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).
			Code(ErrCodeClinicIDEmpty).
			Public("Clinic ID is required.").
			Errorf("clinic id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).
			Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").
			Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var clinicExists int64
		if err := tx.Raw(`SELECT count(*) FROM domain.clinics WHERE id = ?`, cmd.ClinicID).
			Scan(&clinicExists).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicLookupFailed).Wrap(err)
		}
		if clinicExists == 0 {
			return oops.In(scopeClinicHead).
				Code(ErrCodeClinicNotFound).
				Public("Clinic not found.").
				With("clinic_id", cmd.ClinicID).
				Errorf("clinic not found")
		}

		var emp model.Employee
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("id = ?", cmd.EmployeeID).
			First(&emp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeEmployeeNotFound).
					Public("Employee not found.").
					With("employee_id", cmd.EmployeeID).
					Errorf("employee not found")
			}
			return oops.In(scopeClinicHead).Code(ErrCodeEmployeeLoadFailed).Wrap(err)
		}

		// Scope check: does this employee's current department belong to the
		// target clinic? There is no denormalized clinic_id on employees, so
		// we JOIN through departments.
		var inClinic int64
		if err := tx.Raw(
			`SELECT count(*) FROM domain.departments WHERE id = ? AND clinic_id = ?`,
			emp.DepartmentID, cmd.ClinicID,
		).Scan(&inClinic).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicLookupFailed).Wrap(err)
		}
		if inClinic == 0 {
			return oops.In(scopeClinicHead).
				Code(ErrCodeEmployeeNotInClinic).
				Public("Employee does not belong to this clinic.").
				With("employee_id", cmd.EmployeeID).
				With("employee_department_id", emp.DepartmentID).
				With("target_clinic_id", cmd.ClinicID).
				Errorf("employee not in target clinic")
		}

		row := model.ClinicHead{
			ClinicID:   cmd.ClinicID,
			EmployeeID: cmd.EmployeeID,
		}
		if err := tx.Create(&row).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
				return oops.In(scopeClinicHead).
					Code(ErrCodeClinicHeadAlreadyAssigned).
					Public("This employee is already a clinic head.").
					With("clinic_id", cmd.ClinicID).
					With("employee_id", cmd.EmployeeID).
					Wrap(err)
			}
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		ev := &clinicv1.ClinicHeadAssigned{
			EmployeeId: cmd.EmployeeID.String(),
		}
		payload, err := anypb.New(ev)
		if err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadEventBuildFailed).Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(row.CreatedAt),
			AggregateType: AggregateTypeClinic,
			AggregateId:   cmd.ClinicID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectClinicHeadAssigned, envelope)
	})
}
