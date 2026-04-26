package membership

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// RemoveClinicHeadDeputyPayload carries the identifiers needed to
// clear the deputy slot on a CH role.
type RemoveClinicHeadDeputyPayload struct {
	ClinicID   string `validate:"required,uuid"`
	EmployeeID string `validate:"required,uuid"`
}

// RemoveClinicHeadDeputyCommand = caller + payload.
type RemoveClinicHeadDeputyCommand struct {
	Caller  authz.Caller
	Payload RemoveClinicHeadDeputyPayload
}

// RemoveClinicHeadDeputy clears the deputy slot. Fails if the slot is
// already empty (no idempotent no-op per spec §4.7).
//
// See: docs/services/Membership.md
func (s *EmployeeService) RemoveClinicHeadDeputy(ctx context.Context, cmd RemoveClinicHeadDeputyCommand) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	clinicID := uuid.MustParse(cmd.Payload.ClinicID)
	employeeID := uuid.MustParse(cmd.Payload.EmployeeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

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

		if !row.DeputyEmployeeID.Valid {
			return oops.In(scopeClinicHead).
				Code(ErrCodeDeputyNotAssigned).
				Public("No deputy is assigned to this role.").
				With("clinic_id", clinicID).
				With("employee_id", employeeID).
				Errorf("deputy not assigned")
		}

		row.DeputyEmployeeID = null.Value[uuid.UUID]{}
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}

		return publishClinicHeadDeputyRemoved(tx, clinicID, employeeID, now)
	})
}
