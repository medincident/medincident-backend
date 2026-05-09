package self

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by GetMyClinicRole.
const (
	ErrCodeSelfClinicRoleNotFound = "self_clinic_role_not_found"
	ErrCodeSelfClinicRoleFailed   = "self_clinic_role_failed"
)

// ClinicRoleView carries the caller's named role within a clinic.
type ClinicRoleView struct {
	IsClinicHead bool
}

// GetMyClinicRole returns whether the caller is the clinic head of the
// given clinic. Returns ErrCodeSelfClinicRoleNotFound when the caller
// is not an active employee of that clinic.
//
// See: docs/services/self/Self.md
func (r *SelfReader) GetMyClinicRole(
	ctx context.Context,
	callerID string,
	clinicID uuid.UUID,
) (*ClinicRoleView, error) {
	var out ClinicRoleView
	err := r.db.WithContext(ctx).Raw(`
		SELECT ch.employee_id IS NOT NULL AS is_clinic_head
		  FROM projections.employee_cards ec
		  LEFT JOIN projections.clinic_heads ch
		         ON ch.employee_id = ec.employee_id
		        AND ch.clinic_id = ?
		 WHERE ec.zitadel_user_id = ?
		   AND ec.clinic_id = ?
		   AND ec.terminated_at IS NULL`,
		clinicID, callerID, clinicID,
	).Row().Scan(&out.IsClinicHead)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.self").
				Code(ErrCodeSelfClinicRoleNotFound).
				Public("Employment record not found in this clinic.").
				With("caller_id", callerID).
				With("clinic_id", clinicID).
				Errorf("not found")
		}
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfClinicRoleFailed).
			With("caller_id", callerID).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	return &out, nil
}
