package authz

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

type clinicHeadOfRole struct{}

// ClinicHeadOf returns policies asserting the caller is a ClinicHead
// of the scope.
var ClinicHeadOf clinicHeadOfRole

type clinicHeadOfClinicPolicy struct{ clinicID uuid.UUID }

// Clinic scopes the ClinicHead check to a specific clinic.
func (clinicHeadOfRole) Clinic(id uuid.UUID) Policy {
	return clinicHeadOfClinicPolicy{clinicID: id}
}

func (p clinicHeadOfClinicPolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.clinicID)
	caller := bc.addCaller()

	direct := fmt.Sprintf(`SELECT 1 FROM domain.clinic_heads h
JOIN domain.employees e ON e.id = h.employee_id
WHERE h.clinic_id = @%s
AND e.zitadel_user_id = @%s`, scope, caller)

	deputy := fmt.Sprintf(`SELECT 1 FROM domain.clinic_heads h
JOIN domain.employees e ON e.id = h.deputy_employee_id
JOIN domain.employee_vacations v ON v.employee_id = h.employee_id
WHERE h.clinic_id = @%s
AND e.zitadel_user_id = @%s AND %s`, scope, caller, activeVacationPredicate)

	return []string{direct, deputy}
}

func (p clinicHeadOfClinicPolicy) describe() string {
	return fmt.Sprintf("clinic head of %s", p.clinicID)
}

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p clinicHeadOfClinicPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With("clinic_id", p.clinicID)
}
