package authz

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// ErrCodePermissionDenied is emitted when an authorization policy
// evaluates to false for the given caller. Maps to
// codes.PermissionDenied at the gRPC boundary.
const ErrCodePermissionDenied = "permission_denied"

// ErrCodeAuthzCheckFailed is emitted when the policy check query
// itself fails (DB fault). Maps to codes.Internal at the gRPC
// boundary.
const ErrCodeAuthzCheckFailed = "authz_check_failed"

// Policy is a single authorization rule. It renders to one or more
// SELECT branches that produce rows iff the caller is authorized
// under this rule. Policies compose via AnyOf into a single
// UNION ALL-based EXISTS query so that the full check stays one
// database round-trip.
//
// The with method mirrors oops.OopsErrorBuilder's value-chaining
// convention so policy decoration composes naturally with the rest
// of the error-building code in this package.
type Policy interface {
	branches(bc *branchCtx) []string
	describe() string
	with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder
}

// branchCtx accumulates SQL placeholder names and named args as a
// Policy tree renders its branches. Each Policy node requests a
// fresh caller placeholder via addCaller and a fresh scope
// placeholder via addScope so that multi-branch policies do not
// collide on parameter names.
type branchCtx struct {
	callerID string
	args     []any
	next     int
}

func (bc *branchCtx) addCaller() string {
	name := fmt.Sprintf("caller%d", bc.next)
	bc.args = append(bc.args, sql.Named(name, bc.callerID))
	bc.next++
	return name
}

func (bc *branchCtx) addScope(id uuid.UUID) string {
	name := fmt.Sprintf("scope%d", bc.next)
	bc.args = append(bc.args, sql.Named(name, id))
	bc.next++
	return name
}

type sysAdminPolicy struct{}

// SystemAdmin authorizes any caller with a row in domain.system_admins.
// Value (not function) — carries no parameters, safe to reuse as a
// singleton across call sites.
var SystemAdmin Policy = sysAdminPolicy{}

func (sysAdminPolicy) branches(bc *branchCtx) []string {
	c := bc.addCaller()
	return []string{fmt.Sprintf(
		"SELECT 1 FROM domain.system_admins WHERE zitadel_user_id = @%s", c)}
}

func (sysAdminPolicy) describe() string { return "system administrator" }

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (sysAdminPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder { return b }

type anyOfPolicy struct{ items []Policy }

// AnyOf composes policies disjunctively. The rendered SQL UNIONs each
// item's branches into one EXISTS so the combined check remains a
// single round-trip regardless of how many items are stacked.
func AnyOf(items ...Policy) Policy { return anyOfPolicy{items: items} }

func (p anyOfPolicy) branches(bc *branchCtx) []string {
	all := make([]string, 0, len(p.items))
	for _, item := range p.items {
		all = append(all, item.branches(bc)...)
	}
	return all
}

func (p anyOfPolicy) describe() string {
	parts := make([]string, 0, len(p.items))
	for _, item := range p.items {
		parts = append(parts, item.describe())
	}
	return strings.Join(parts, " or ")
}

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p anyOfPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	for _, item := range p.items {
		b = item.with(b)
	}
	return b
}

type orgAdminOfRole struct{}

// OrgAdminOf is the namespace returning OrgAdmin policies scoped by
// one of the supported child entities. Each method produces a Policy
// with two branches: "caller is the direct holder" and "caller is the
// deputy while the holder is on an active vacation".
var OrgAdminOf orgAdminOfRole

type orgAdminPolicy struct {
	field     string
	id        uuid.UUID
	clauseFmt string // format string with one %s for scope placeholder
}

func (p orgAdminPolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.id)
	caller := bc.addCaller()
	clause := fmt.Sprintf(p.clauseFmt, scope)

	direct := fmt.Sprintf(`SELECT 1 FROM domain.org_admins oa
JOIN domain.employees e ON e.id = oa.employee_id
%s
AND e.zitadel_user_id = @%s`, clause, caller)

	deputy := fmt.Sprintf(`SELECT 1 FROM domain.org_admins oa
JOIN domain.employees e ON e.id = oa.deputy_employee_id
JOIN domain.employee_vacations v ON v.employee_id = oa.employee_id
%s
AND e.zitadel_user_id = @%s AND %s`, clause, caller, activeVacationPredicate)

	return []string{direct, deputy}
}

func (orgAdminPolicy) describe() string { return "organization administrator" }

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p orgAdminPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With(p.field, p.id)
}

func (orgAdminOfRole) Organization(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field:     "organization_id",
		id:        id,
		clauseFmt: "WHERE oa.organization_id = @%s",
	}
}

func (orgAdminOfRole) Clinic(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field:     "clinic_id",
		id:        id,
		clauseFmt: "JOIN domain.clinics c ON c.organization_id = oa.organization_id WHERE c.id = @%s",
	}
}

func (orgAdminOfRole) Department(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field: "department_id",
		id:    id,
		clauseFmt: "JOIN domain.clinics c ON c.organization_id = oa.organization_id " +
			"JOIN domain.departments d ON d.clinic_id = c.id WHERE d.id = @%s",
	}
}

func (orgAdminOfRole) Employee(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field:     "employee_id",
		id:        id,
		clauseFmt: "JOIN domain.employees tgt ON tgt.organization_id = oa.organization_id WHERE tgt.id = @%s",
	}
}

func (orgAdminOfRole) Vacation(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field: "vacation_id",
		id:    id,
		clauseFmt: "JOIN domain.employees tgt ON tgt.organization_id = oa.organization_id " +
			"JOIN domain.employee_vacations vt ON vt.employee_id = tgt.id WHERE vt.id = @%s",
	}
}

func (orgAdminOfRole) Category(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field:     "category_id",
		id:        id,
		clauseFmt: "JOIN domain.incident_categories ic ON ic.organization_id = oa.organization_id WHERE ic.id = @%s",
	}
}

func (orgAdminOfRole) IncidentType(id uuid.UUID) Policy {
	return orgAdminPolicy{
		field:     "type_id",
		id:        id,
		clauseFmt: "JOIN domain.incident_types it ON it.organization_id = oa.organization_id WHERE it.id = @%s",
	}
}
