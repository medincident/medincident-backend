package authz

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// orgHeadOfRole is a namespace for OrgHead-based policies.
type orgHeadOfRole struct{}

// OrgHeadOf returns policies asserting the caller is OrgHead of the
// scope. Like OrgAdminOf, the deputy branch is gated by an active
// vacation on the principal.
var OrgHeadOf orgHeadOfRole

type orgHeadOfOrgPolicy struct{ orgID uuid.UUID }

// Organization scopes the OrgHead check to a specific organization.
func (orgHeadOfRole) Organization(id uuid.UUID) Policy {
	return orgHeadOfOrgPolicy{orgID: id}
}

func (p orgHeadOfOrgPolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.orgID)
	caller := bc.addCaller()

	direct := fmt.Sprintf(`SELECT 1 FROM domain.org_heads h
JOIN domain.employees e ON e.id = h.employee_id
WHERE h.organization_id = @%s
AND e.zitadel_user_id = @%s`, scope, caller)

	deputy := fmt.Sprintf(`SELECT 1 FROM domain.org_heads h
JOIN domain.employees e ON e.id = h.deputy_employee_id
JOIN domain.employee_vacations v ON v.employee_id = h.employee_id
WHERE h.organization_id = @%s
AND e.zitadel_user_id = @%s AND %s`, scope, caller, activeVacationPredicate)

	return []string{direct, deputy}
}

func (p orgHeadOfOrgPolicy) describe() string {
	return fmt.Sprintf("organization head of %s", p.orgID)
}

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p orgHeadOfOrgPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With("organization_id", p.orgID)
}
