package membership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// AssignClinicHeadDeputyPayload carries the identifiers needed to set
// the deputy slot on an existing CH role.
type AssignClinicHeadDeputyPayload struct {
	ClinicID         string `validate:"required,uuid"`
	EmployeeID       string `validate:"required,uuid"`
	DeputyEmployeeID string `validate:"required,uuid"`
}

// AssignClinicHeadDeputyCommand = caller + payload.
type AssignClinicHeadDeputyCommand struct {
	Caller  authz.Caller
	Payload AssignClinicHeadDeputyPayload
}

// AssignClinicHeadDeputy sets the deputy slot on an existing CH role.
// See spec §8.5 (ClinicHead variant).
func (s *EmployeeService) AssignClinicHeadDeputy(ctx context.Context, cmd AssignClinicHeadDeputyCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ClinicID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	deputyEmployeeID := uuid.MustParse(cmd.Payload.DeputyEmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.ClinicHead
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("clinic_id = ? AND employee_id = ?", clinicID, employeeID).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeClinicHeadNotFound).
					Public("Clinic head not found.").
					With("clinic_id", clinicID).
					With("employee_id", employeeID).
					Errorf("not found")
			}
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadLoadFailed).Wrap(err)
		}

		var deputy model.Employee
		err = tx.Where("id = ?", deputyEmployeeID).First(&deputy).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scopeClinicHead).
					Code(ErrCodeDeputyNotFound).
					Public("Deputy employee not found.").
					With("deputy_employee_id", deputyEmployeeID).
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

		if deputyEmployeeID == employeeID {
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

		row.DeputyEmployeeID = null.ValueFrom(deputyEmployeeID)
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		return projector.ClinicHeadDeputyAssigned(tx, &row)
	})
}
