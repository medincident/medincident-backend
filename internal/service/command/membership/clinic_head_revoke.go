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
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// RevokeClinicHeadCommand carries the identifiers needed to remove an
// employee's clinic head role.
type RevokeClinicHeadCommand struct {
	ClinicID   uuid.UUID `validate:"required"`
	EmployeeID uuid.UUID `validate:"required"`
}

// RevokeClinicHead removes the role row and publishes the Revoked
// event. If a deputy was assigned, a DeputyRemoved event is published
// FIRST (Rule 2 — cleanup before terminate). See spec §8.6.
func (s *EmployeeService) RevokeClinicHead(ctx context.Context, cmd RevokeClinicHeadCommand) error {
	if err := validation.Struct(cmd); err != nil {
		return err
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

		// Rule 2 — cleanup before terminate.
		if row.DeputyEmployeeID.Valid {
			if err := publishClinicHeadDeputyRemoved(tx, cmd.ClinicID, cmd.EmployeeID, now); err != nil {
				return err
			}
		}

		if err := tx.Delete(&model.ClinicHead{}, "clinic_id = ? AND employee_id = ?",
			cmd.ClinicID, cmd.EmployeeID).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadDeleteFailed).Wrap(err)
		}

		return publishClinicHeadRevoked(tx, cmd.ClinicID, cmd.EmployeeID, now)
	})
}

// publishClinicHeadRevoked is a shared helper for Revoke,
// cascade-on-transfer, and cascade-on-terminate. It does NOT delete
// the domain row; the caller owns that. The projector call here
// deletes the matching projection row so explicit and cascaded
// revokes stay in sync.
func publishClinicHeadRevoked(tx *gorm.DB, clinicID, employeeID uuid.UUID, _ time.Time) error {
	return projector.ClinicHeadRevoked(tx, clinicID, employeeID)
}

// publishClinicHeadDeputyRemoved is a shared helper; it clears the
// deputy slot on the projection row. The caller owns the domain-row
// update.
func publishClinicHeadDeputyRemoved(tx *gorm.DB, clinicID, employeeID uuid.UUID, now time.Time) error {
	return projector.ClinicHeadDeputyRemoved(tx, clinicID, employeeID, now)
}
