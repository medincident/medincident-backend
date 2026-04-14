package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-command-service/internal/model"
)

// RemoveClinicHeadDeputyCommand carries the identifiers needed to
// clear the deputy slot on a CH role.
type RemoveClinicHeadDeputyCommand struct {
	ClinicID   uuid.UUID
	EmployeeID uuid.UUID
}

// RemoveClinicHeadDeputy clears the deputy slot. Fails if the slot is
// already empty (no idempotent no-op per spec §4.7).
func (s *EmployeeService) RemoveClinicHeadDeputy(ctx context.Context, cmd RemoveClinicHeadDeputyCommand) error {
	var errs []error
	if cmd.ClinicID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).Code(ErrCodeClinicIDEmpty).
			Public("Clinic ID is required.").Errorf("clinic id empty"))
	}
	if cmd.EmployeeID == uuid.Nil {
		errs = append(errs, oops.In(scopeClinicHead).Code(ErrCodeEmployeeIDEmpty).
			Public("Employee ID is required.").Errorf("employee id empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		if row.DeputyEmployeeID == nil {
			return oops.In(scopeClinicHead).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("clinic_id", cmd.ClinicID).
				With("employee_id", cmd.EmployeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = nil
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		return publishClinicHeadDeputyRemoved(tx, cmd.ClinicID, cmd.EmployeeID, now)
	})
}
