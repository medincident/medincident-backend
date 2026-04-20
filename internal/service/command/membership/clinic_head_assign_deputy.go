package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// AssignClinicHeadDeputyCommand carries the identifiers needed to set
// the deputy slot on an existing CH role.
type AssignClinicHeadDeputyCommand struct {
	ClinicID         uuid.UUID `validate:"required"`
	EmployeeID       uuid.UUID `validate:"required"`
	DeputyEmployeeID uuid.UUID `validate:"required"`
}

// AssignClinicHeadDeputy sets the deputy slot on an existing CH role.
// See spec §8.5 (ClinicHead variant).
func (s *EmployeeService) AssignClinicHeadDeputy(ctx context.Context, cmd AssignClinicHeadDeputyCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
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
		// same clinic as the role? Look up the department since there is
		// no denormalized clinic_id on employees.
		var deputyInClinic int64
		if err := tx.Raw(
			`SELECT count(*) FROM domain.departments WHERE id = ? AND clinic_id = ?`,
			deputy.DepartmentID, row.ClinicID,
		).Scan(&deputyInClinic).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeDepartmentLookupFailed).Wrap(err)
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

		return projector.ClinicHeadDeputyAssigned(tx, &row)
	})
}
