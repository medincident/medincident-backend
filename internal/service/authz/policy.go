package authz

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// ---------------------------------------------------------------------
// Error codes
// ---------------------------------------------------------------------

// ErrCodePermissionDenied is emitted when an authorization policy
// evaluates to false for the given caller. Maps to
// codes.PermissionDenied at the gRPC boundary.
const ErrCodePermissionDenied = "permission_denied"

// ErrCodeAuthzCheckFailed is emitted when the policy check query
// itself fails (DB fault). Maps to codes.Internal at the gRPC
// boundary.
const ErrCodeAuthzCheckFailed = "authz_check_failed"

// ---------------------------------------------------------------------
// Policy interface + branch-rendering context
// ---------------------------------------------------------------------

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
// Policy tree renders its branches. Caller and scope counters are
// independent so a typical composed policy yields natural names
// (caller0, caller1, ..., scope0, scope1, ...) rather than a single
// interleaved sequence.
type branchCtx struct {
	callerID   string
	args       []any
	nextCaller int
	nextScope  int
}

func (bc *branchCtx) addCaller() string {
	name := fmt.Sprintf("caller%d", bc.nextCaller)
	bc.args = append(bc.args, sql.Named(name, bc.callerID))
	bc.nextCaller++
	return name
}

func (bc *branchCtx) addScope(id uuid.UUID) string {
	name := fmt.Sprintf("scope%d", bc.nextScope)
	bc.args = append(bc.args, sql.Named(name, id))
	bc.nextScope++
	return name
}

// hasZeroScope reports whether any scope bound into bc is the zero
// UUID. The zero UUID passes go-playground/validator's "uuid" tag,
// which would otherwise let a whitespace bug at the transport layer
// turn into a "silently denies everything" query. Require
// short-circuits to ErrCodeAuthzCheckFailed when this is true.
func (bc *branchCtx) hasZeroScope() bool {
	for _, arg := range bc.args {
		named, ok := arg.(sql.NamedArg)
		if !ok {
			continue
		}
		if id, ok := named.Value.(uuid.UUID); ok && id == uuid.Nil {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------
// Authenticated — scope-less "caller is logged in" policy
// ---------------------------------------------------------------------

type authenticatedPolicy struct{}

// Authenticated authorizes any caller that reached the policy check —
// the authn interceptor upstream rejects empty caller IDs, so
// reaching Require with a non-empty caller ID is itself the proof.
// The branch is a constant `SELECT 1` that returns a row regardless
// of database state, so the composed EXISTS always evaluates true.
//
// Use for reads that must be gated by authentication but not by
// membership — e.g. the patient-facing classifier endpoints, where
// patients pick an organization to file an incident against without
// being employees of it.
var Authenticated Policy = authenticatedPolicy{}

func (authenticatedPolicy) branches(_ *branchCtx) []string {
	return []string{"SELECT 1"}
}

func (authenticatedPolicy) describe() string { return "an authenticated caller" }

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (authenticatedPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder { return b }

// ---------------------------------------------------------------------
// SystemAdmin — scope-less role
// ---------------------------------------------------------------------

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

// ---------------------------------------------------------------------
// AnyOf — disjunctive composition
// ---------------------------------------------------------------------

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

// anyOfPolicy.with accumulates context from every item, not just the
// one that would have matched — there is no "matched" item on the
// deny path. This is intentional but has a subtle consequence: if you
// compose AnyOf of two scope-bearing policies (e.g.
// AnyOf(OrgAdminOf.Organization(a), OrgAdminOf.Clinic(b))), the
// emitted error carries both "organization_id" and "clinic_id" tags.
// In practice every AdminOf.X is AnyOf(SystemAdmin, OrgAdminOf.X(id))
// and sysAdminPolicy.with is a no-op, so exactly one scope key is
// added. Preserve this shape if you add new top-level batteries.
//
//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p anyOfPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	for _, item := range p.items {
		b = item.with(b)
	}
	return b
}

// ---------------------------------------------------------------------
// OrgAdminOf — organization-admin scoped by one of seven child entities
// ---------------------------------------------------------------------

type orgAdminOfRole struct{}

// OrgAdminOf is the namespace returning OrgAdmin policies scoped by
// one of the supported child entities. Each method produces a Policy
// with two branches: "caller is the direct holder" and "caller is the
// deputy while the holder is on an active vacation".
//
// Production call sites should prefer AdminOf.X(id) — which is
// AnyOf(SystemAdmin, OrgAdminOf.X(id)) — so the sysadmin bypass stays
// visible at the call. OrgAdminOf.X(id) alone intentionally omits the
// sysadmin disjunct and is exported only because integration tests
// need to exercise the bare branch in isolation. Using it from a
// service method silently excludes system admins from the operation.
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

// ---------------------------------------------------------------------
// AdminOf — convenience battery: AnyOf(SystemAdmin, OrgAdminOf.X(id))
// ---------------------------------------------------------------------

type adminOfBattery struct{}

// AdminOf packs "system admin OR organization admin" as one
// constructor per scope. Each method returns
// AnyOf(SystemAdmin, OrgAdminOf.X(id)) so call sites stay short
// while the sysadmin bypass remains visible in the call.
var AdminOf adminOfBattery

func (adminOfBattery) Organization(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Organization(id))
}

func (adminOfBattery) Clinic(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Clinic(id))
}

func (adminOfBattery) Department(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Department(id))
}

func (adminOfBattery) Employee(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Employee(id))
}

func (adminOfBattery) Vacation(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Vacation(id))
}

func (adminOfBattery) Category(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Category(id))
}

func (adminOfBattery) IncidentType(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.IncidentType(id))
}

// ---------------------------------------------------------------------
// MemberOf — organization membership (any employee of the scope)
// ---------------------------------------------------------------------

type memberOfRole struct{}

// MemberOf is the namespace returning "caller is an employee of the
// organization that owns the supplied scope" policies. Scopes narrow
// syntactically (Clinic, Department) but every check widens to the
// containing organization at the SQL level, mirroring OrgAdminOf: both
// "admin of clinic X" and "member of clinic X" mean the organization
// owning clinic X. This matches the read model where any employee of
// an organization can see the whole organization's catalog.
//
// Unlike OrgAdminOf, MemberOf has no deputy-during-vacation branch:
// "member" is a static fact over domain.employees, not a role with
// delegation.
//
// MemberOf is exported for test isolation; production reads should
// prefer ReaderOf.X(id), which bundles SystemAdmin, OrgAdminOf.X(id),
// and MemberOf.X(id) into one battery so the privilege ladder stays
// visible at the call site.
var MemberOf memberOfRole

type memberOfPolicy struct {
	field string
	id    uuid.UUID
	// clauseFmt is a SQL fragment that joins domain.employees `e` to
	// the scope being checked and restricts it by the scope id. The
	// fragment MUST be an optional JOIN chain followed by a WHERE
	// clause on the scope placeholder — branches appends
	// "AND e.zitadel_user_id = @caller..." so the fragment is required
	// to end on a filter predicate rather than a plain join. Every
	// constructor below follows that shape; break the convention and
	// rendered SQL becomes invalid at runtime.
	clauseFmt string // "(JOIN ...)* WHERE <scope>.id = @%s"
}

func (p memberOfPolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.id)
	caller := bc.addCaller()
	clause := fmt.Sprintf(p.clauseFmt, scope)
	return []string{fmt.Sprintf(`SELECT 1 FROM domain.employees e
%s
AND e.zitadel_user_id = @%s`, clause, caller)}
}

func (memberOfPolicy) describe() string { return "organization member" }

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p memberOfPolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With(p.field, p.id)
}

func (memberOfRole) Organization(id uuid.UUID) Policy {
	return memberOfPolicy{
		field:     "organization_id",
		id:        id,
		clauseFmt: "WHERE e.organization_id = @%s",
	}
}

func (memberOfRole) Clinic(id uuid.UUID) Policy {
	return memberOfPolicy{
		field:     "clinic_id",
		id:        id,
		clauseFmt: "JOIN domain.clinics c ON c.organization_id = e.organization_id WHERE c.id = @%s",
	}
}

func (memberOfRole) Department(id uuid.UUID) Policy {
	return memberOfPolicy{
		field: "department_id",
		id:    id,
		clauseFmt: "JOIN domain.clinics c ON c.organization_id = e.organization_id " +
			"JOIN domain.departments d ON d.clinic_id = c.id WHERE d.id = @%s",
	}
}

func (memberOfRole) Employee(id uuid.UUID) Policy {
	return memberOfPolicy{
		field:     "employee_id",
		id:        id,
		clauseFmt: "JOIN domain.employees tgt ON tgt.organization_id = e.organization_id WHERE tgt.id = @%s",
	}
}

func (memberOfRole) Category(id uuid.UUID) Policy {
	return memberOfPolicy{
		field:     "category_id",
		id:        id,
		clauseFmt: "JOIN domain.incident_categories ic ON ic.organization_id = e.organization_id WHERE ic.id = @%s",
	}
}

func (memberOfRole) IncidentType(id uuid.UUID) Policy {
	return memberOfPolicy{
		field:     "type_id",
		id:        id,
		clauseFmt: "JOIN domain.incident_types it ON it.organization_id = e.organization_id WHERE it.id = @%s",
	}
}

// ---------------------------------------------------------------------
// SelfEmployee — caller IS the target employee
// ---------------------------------------------------------------------

type selfEmployeePolicy struct{ id uuid.UUID }

// SelfEmployee authorizes the caller when they are the employee row
// identified by id. Used for self-service reads where "only admins or
// the subject themselves" is the right gate — e.g. an employee's own
// vacation history.
func SelfEmployee(id uuid.UUID) Policy { return selfEmployeePolicy{id: id} }

func (p selfEmployeePolicy) branches(bc *branchCtx) []string {
	scope := bc.addScope(p.id)
	caller := bc.addCaller()
	return []string{fmt.Sprintf(
		`SELECT 1 FROM domain.employees e
WHERE e.id = @%s AND e.zitadel_user_id = @%s`, scope, caller)}
}

func (selfEmployeePolicy) describe() string { return "the employee themselves" }

//nolint:gocritic // hugeParam: mirrors oops.OopsErrorBuilder's value-chaining API.
func (p selfEmployeePolicy) with(b oops.OopsErrorBuilder) oops.OopsErrorBuilder {
	return b.With("employee_id", p.id)
}

// ---------------------------------------------------------------------
// ReaderOf — read-side battery: SystemAdmin OR OrgAdminOf.X OR MemberOf.X
// ---------------------------------------------------------------------

type readerOfBattery struct{}

// ReaderOf packs the three-level read privilege — "system admin OR
// organization admin OR organization member" — as one constructor per
// scope. It is the standard policy for query-side handlers: any
// employee of the owning organization may read catalog data, org
// admins retain their write-path privilege, and system admins bypass
// both. Scopes below the organization (Clinic, Department) narrow the
// member check to the containing subtree while the admin/sysadmin
// branches stay organization-wide via existing OrgAdminOf clauses.
var ReaderOf readerOfBattery

func (readerOfBattery) Organization(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Organization(id), MemberOf.Organization(id))
}

func (readerOfBattery) Clinic(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Clinic(id), MemberOf.Clinic(id))
}

func (readerOfBattery) Department(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Department(id), MemberOf.Department(id))
}

func (readerOfBattery) Employee(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Employee(id), MemberOf.Employee(id))
}

func (readerOfBattery) Category(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.Category(id), MemberOf.Category(id))
}

func (readerOfBattery) IncidentType(id uuid.UUID) Policy {
	return AnyOf(SystemAdmin, OrgAdminOf.IncidentType(id), MemberOf.IncidentType(id))
}

// ---------------------------------------------------------------------
// Require — entry point: renders policy to SQL, executes, shapes error
// ---------------------------------------------------------------------

// Require evaluates policy p for callerID. It renders the policy tree
// into a single SELECT EXISTS(... UNION ALL ...) query so the whole
// check is one database round-trip regardless of how many sub-policies
// compose through AnyOf. Returns:
//   - nil on success;
//   - an ErrCodePermissionDenied domain error on deny (Public message
//     derived from policy.describe());
//   - an ErrCodeAuthzCheckFailed domain error wrapping the DB error if
//     the check query itself cannot run.
func (a *Authz) Require(ctx context.Context, callerID string, p Policy) error {
	if p == nil {
		return oops.In("service.authz").
			Code(ErrCodeAuthzCheckFailed).
			With("caller_id", callerID).
			Errorf("nil policy")
	}
	// Authenticated-only gate: the authn interceptor upstream rejects
	// empty caller IDs, but Require is exported and can be called from
	// tests or future code paths that skip the interceptor. Guard
	// explicitly and short-circuit the DB round-trip — Authenticated
	// never depends on a scope or row, so the EXISTS query would only
	// add latency.
	if _, ok := p.(authenticatedPolicy); ok {
		if callerID == "" {
			return oops.In("service.authz").
				Code(ErrCodePermissionDenied).
				Public("Access denied: requires an authenticated caller.").
				Errorf("empty caller id for authenticated policy")
		}
		return nil
	}
	bc := &branchCtx{callerID: callerID}
	branches := p.branches(bc)
	if len(branches) == 0 {
		return oops.In("service.authz").
			Code(ErrCodeAuthzCheckFailed).
			With("caller_id", callerID).
			Errorf("policy produced no branches")
	}
	if bc.hasZeroScope() {
		return oops.In("service.authz").
			Code(ErrCodeAuthzCheckFailed).
			With("caller_id", callerID).
			Errorf("policy scope is zero uuid")
	}
	query := "SELECT EXISTS(" + strings.Join(branches, " UNION ALL ") + ")"

	var ok bool
	if err := a.db.WithContext(ctx).Raw(query, bc.args...).Scan(&ok).Error; err != nil {
		return oops.In("service.authz").
			Code(ErrCodeAuthzCheckFailed).
			With("caller_id", callerID).
			Wrap(err)
	}
	if ok {
		return nil
	}

	b := oops.In("service.authz").
		Code(ErrCodePermissionDenied).
		Public(fmt.Sprintf("Access denied: requires %s privileges.", p.describe())).
		With("caller_id", callerID)
	b = p.with(b)
	return b.Errorf("permission denied")
}
