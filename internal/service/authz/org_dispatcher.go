package authz

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

type orgDispatcherOfRole struct{}

// OrgDispatcherOf returns policies asserting the caller is an
// OrgDispatcher of the scope.
var OrgDispatcherOf orgDispatcherOfRole

type orgDispatcherOfOrgPolicy struct{ orgID uuid.UUID }

// Organization scopes the OrgDispatcher check to a specific organization.
func (orgDispatcherOfRole) Organization(id uuid.UUID) Policy {
	return orgDispatcherOfOrgPolicy{orgID: id}
}

func (p orgDispatcherOfOrgPolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.orgID)
	caller := bc.addCaller()

	direct := fmt.Sprintf(`SELECT 1 FROM domain.org_dispatchers d
JOIN domain.employees e ON e.id = d.employee_id
WHERE d.organization_id = @%s
AND e.zitadel_user_id = @%s`, scope, caller)

	deputy := fmt.Sprintf(`SELECT 1 FROM domain.org_dispatchers d
JOIN domain.employees e ON e.id = d.deputy_employee_id
JOIN domain.employee_vacations v ON v.employee_id = d.employee_id
WHERE d.organization_id = @%s
AND e.zitadel_user_id = @%s AND %s`, scope, caller, activeVacationPredicate)

	return []string{direct, deputy}
}

func (p orgDispatcherOfOrgPolicy) describe() string {
	return fmt.Sprintf("organization dispatcher of %s", p.orgID)
}

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p orgDispatcherOfOrgPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With("organization_id", p.orgID)
}
