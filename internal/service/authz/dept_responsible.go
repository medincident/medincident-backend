package authz

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

type deptResponsibleOfRole struct{}

// DeptResponsibleOf returns policies asserting the caller is a
// DepartmentResponsible of the scope.
var DeptResponsibleOf deptResponsibleOfRole

type deptResponsibleOfDeptPolicy struct{ deptID uuid.UUID }

// Department scopes the DeptResponsible check to a specific department.
func (deptResponsibleOfRole) Department(id uuid.UUID) Policy {
	return deptResponsibleOfDeptPolicy{deptID: id}
}

func (p deptResponsibleOfDeptPolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.deptID)
	caller := bc.addCaller()

	direct := fmt.Sprintf(`SELECT 1 FROM domain.department_responsibles dr
JOIN domain.employees e ON e.id = dr.employee_id
WHERE dr.department_id = @%s
AND e.zitadel_user_id = @%s`, scope, caller)

	deputy := fmt.Sprintf(`SELECT 1 FROM domain.department_responsibles dr
JOIN domain.employees e ON e.id = dr.deputy_employee_id
JOIN domain.employee_vacations v ON v.employee_id = dr.employee_id
WHERE dr.department_id = @%s
AND e.zitadel_user_id = @%s AND %s`, scope, caller, activeVacationPredicate)

	return []string{direct, deputy}
}

func (p deptResponsibleOfDeptPolicy) describe() string {
	return fmt.Sprintf("department responsible of %s", p.deptID)
}

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p deptResponsibleOfDeptPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With("department_id", p.deptID)
}
