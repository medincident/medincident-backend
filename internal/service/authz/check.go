package authz

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by check-layer queries. Each code names the
// specific check that failed so the operator log has actionable
// signal; all of them map to codes.Internal since they indicate a DB
// fault, not a client problem.
const (
	ErrCodeAuthzSystemAdminCheckFailed      = "authz_system_admin_check_failed"
	ErrCodeAuthzOrgAccessCheckFailed        = "authz_org_access_check_failed"
	ErrCodeAuthzClinicAccessCheckFailed     = "authz_clinic_access_check_failed"
	ErrCodeAuthzDepartmentAccessCheckFailed = "authz_department_access_check_failed"
	ErrCodeAuthzEmployeeAccessCheckFailed   = "authz_employee_access_check_failed"
	ErrCodeAuthzVacationAccessCheckFailed   = "authz_vacation_access_check_failed"
	ErrCodeAuthzCategoryAccessCheckFailed   = "authz_category_access_check_failed"
	ErrCodeAuthzTypeAccessCheckFailed       = "authz_type_access_check_failed"
)

// wrapCheckErr is the common shape for turning a gorm query failure
// into a domain error with enough context to identify the check.
func wrapCheckErr(code string, err error, fields ...func(oops.OopsErrorBuilder) oops.OopsErrorBuilder) error {
	b := oops.In("service.authz").Code(code)
	for _, f := range fields {
		b = f(b)
	}
	return b.Wrap(err)
}

// scanExists runs the given raw SQL with named args and returns the
// boolean existence result. Centralises the gorm scan so each check
// method stays a 3-4 line wrapper.
func (a *Authz) scanExists(ctx context.Context, query string, args []any, code string, fields ...func(oops.OopsErrorBuilder) oops.OopsErrorBuilder) (bool, error) {
	var exists bool
	if err := a.db.WithContext(ctx).Raw(query, args...).Scan(&exists).Error; err != nil {
		return false, wrapCheckErr(code, err, fields...)
	}
	return exists, nil
}

// ---------------------------------------------------------------------
// System admin
// ---------------------------------------------------------------------

const querySystemAdmin = `
	SELECT EXISTS(
		SELECT 1 FROM domain.system_admins WHERE zitadel_user_id = @caller
	)
`

// checkSystemAdmin: 1 query. Membership in the global system_admins
// table grants unconditional access.
func (a *Authz) checkSystemAdmin(ctx context.Context, callerID string) (bool, error) {
	return a.scanExists(ctx, querySystemAdmin,
		[]any{named("caller", callerID)},
		ErrCodeAuthzSystemAdminCheckFailed,
		withCaller(callerID),
	)
}

// ---------------------------------------------------------------------
// Org access (direct org ID, no scope join)
// ---------------------------------------------------------------------

// queryOrgAccess combines three grants into one round-trip:
//  1. caller is a system admin;
//  2. caller's employee row is a direct org admin in this org;
//  3. caller's employee row is the deputy on an org admin whose
//     holder is currently on active vacation.
var queryOrgAccess = fmt.Sprintf(`
	SELECT EXISTS(
		SELECT 1 FROM domain.system_admins WHERE zitadel_user_id = @caller
		UNION ALL
		SELECT 1 FROM domain.org_admins oa
		JOIN domain.employees e ON e.id = oa.employee_id
		WHERE oa.organization_id = @scope AND e.zitadel_user_id = @caller
		UNION ALL
		SELECT 1 FROM domain.org_admins oa
		JOIN domain.employees e ON e.id = oa.deputy_employee_id
		JOIN domain.employee_vacations v ON v.employee_id = oa.employee_id
		WHERE oa.organization_id = @scope AND e.zitadel_user_id = @caller
		  AND %s
	)
`, activeVacationPredicate)

func (a *Authz) checkOrgAccess(ctx context.Context, callerID string, orgID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccess,
		[]any{named("caller", callerID), named("scope", orgID)},
		ErrCodeAuthzOrgAccessCheckFailed,
		withCaller(callerID), withScope("organization_id", orgID),
	)
}

// ---------------------------------------------------------------------
// Org access via child scope (clinic, department, employee, vacation,
// category, type)
//
// Each query joins the scope table to clinic.organization_id /
// employee.organization_id / etc. so the caller's admin/deputy
// membership is checked against the resolved organization in a single
// round-trip. The sysadmin branch deliberately does NOT join the
// scope table — sysadmins authorize for any scope, and the service
// layer is the authoritative source of "entity not found" for them.
//
// For non-sysadmins, if the scope id does not exist or belongs to a
// different organization, every JOIN branch produces zero rows and
// the overall EXISTS evaluates to false, returning permission_denied.
// This prevents cross-tenant UUID enumeration via differentiated
// NotFound vs PermissionDenied responses.
// ---------------------------------------------------------------------

// buildViaScopeQuery renders the via-scope query template once per
// scope type. joinFragment is SQL that joins the scope table against
// oa.organization_id (e.g. "JOIN domain.clinics c ON
// c.organization_id = oa.organization_id WHERE c.id = @scope").
// Controlled input — callers pass package-level string literals, no
// user data is interpolated.
func buildViaScopeQuery(joinFragment string) string {
	return fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM domain.system_admins WHERE zitadel_user_id = @caller
			UNION ALL
			SELECT 1 FROM domain.org_admins oa
			JOIN domain.employees e ON e.id = oa.employee_id
			%s
			AND e.zitadel_user_id = @caller
			UNION ALL
			SELECT 1 FROM domain.org_admins oa
			JOIN domain.employees e ON e.id = oa.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = oa.employee_id
			%s
			AND e.zitadel_user_id = @caller AND %s
		)
	`, joinFragment, joinFragment, activeVacationPredicate)
}

var (
	queryOrgAccessViaClinic = buildViaScopeQuery(
		`JOIN domain.clinics c ON c.organization_id = oa.organization_id WHERE c.id = @scope`)

	queryOrgAccessViaDepartment = buildViaScopeQuery(
		`JOIN domain.clinics c ON c.organization_id = oa.organization_id
		 JOIN domain.departments d ON d.clinic_id = c.id WHERE d.id = @scope`)

	queryOrgAccessViaEmployee = buildViaScopeQuery(
		`JOIN domain.employees tgt ON tgt.organization_id = oa.organization_id WHERE tgt.id = @scope`)

	queryOrgAccessViaVacation = buildViaScopeQuery(
		`JOIN domain.employees tgt ON tgt.organization_id = oa.organization_id
		 JOIN domain.employee_vacations vt ON vt.employee_id = tgt.id WHERE vt.id = @scope`)

	queryOrgAccessViaCategory = buildViaScopeQuery(
		`JOIN domain.incident_categories ic ON ic.organization_id = oa.organization_id WHERE ic.id = @scope`)

	queryOrgAccessViaType = buildViaScopeQuery(
		`JOIN domain.incident_types it ON it.organization_id = oa.organization_id WHERE it.id = @scope`)
)

func (a *Authz) checkOrgAccessViaClinic(ctx context.Context, callerID string, clinicID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccessViaClinic,
		[]any{named("caller", callerID), named("scope", clinicID)},
		ErrCodeAuthzClinicAccessCheckFailed,
		withCaller(callerID), withScope("clinic_id", clinicID),
	)
}

func (a *Authz) checkOrgAccessViaDepartment(ctx context.Context, callerID string, deptID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccessViaDepartment,
		[]any{named("caller", callerID), named("scope", deptID)},
		ErrCodeAuthzDepartmentAccessCheckFailed,
		withCaller(callerID), withScope("department_id", deptID),
	)
}

func (a *Authz) checkOrgAccessViaEmployee(ctx context.Context, callerID string, empID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccessViaEmployee,
		[]any{named("caller", callerID), named("scope", empID)},
		ErrCodeAuthzEmployeeAccessCheckFailed,
		withCaller(callerID), withScope("employee_id", empID),
	)
}

func (a *Authz) checkOrgAccessViaVacation(ctx context.Context, callerID string, vacID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccessViaVacation,
		[]any{named("caller", callerID), named("scope", vacID)},
		ErrCodeAuthzVacationAccessCheckFailed,
		withCaller(callerID), withScope("vacation_id", vacID),
	)
}

func (a *Authz) checkOrgAccessViaCategory(ctx context.Context, callerID string, catID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccessViaCategory,
		[]any{named("caller", callerID), named("scope", catID)},
		ErrCodeAuthzCategoryAccessCheckFailed,
		withCaller(callerID), withScope("category_id", catID),
	)
}

func (a *Authz) checkOrgAccessViaType(ctx context.Context, callerID string, typeID uuid.UUID) (bool, error) {
	return a.scanExists(ctx, queryOrgAccessViaType,
		[]any{named("caller", callerID), named("scope", typeID)},
		ErrCodeAuthzTypeAccessCheckFailed,
		withCaller(callerID), withScope("type_id", typeID),
	)
}

// ---------------------------------------------------------------------
// Small adapters
// ---------------------------------------------------------------------

// named is a thin shim around sql.Named (which gorm's Raw honours as
// an @name placeholder) so call sites stay readable.
func named(name string, value any) any {
	return sql.Named(name, value)
}

func withCaller(callerID string) func(oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return func(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
		return b.With("caller_id", callerID)
	}
}

func withScope(field string, id uuid.UUID) func(oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return func(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
		return b.With(field, id)
	}
}
