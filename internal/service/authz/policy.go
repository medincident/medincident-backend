package authz

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

const ErrCodeAuthzCheckFailed = "authz_check_failed"

// Policy is a single authorization rule. It renders to one or more
// SELECT branches that produce rows iff the caller is authorized
// under this rule. Policies compose via AnyOf into a single
// UNION ALL-based EXISTS query so that the full check stays one
// database round-trip.
type Policy interface {
	branches(bc *branchCtx) []string
	describe() string
	with(b *oops.OopsErrorBuilder) *oops.OopsErrorBuilder
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

func (sysAdminPolicy) describe() string                                     { return "system administrator" }
func (sysAdminPolicy) with(b *oops.OopsErrorBuilder) *oops.OopsErrorBuilder { return b }

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

func (p anyOfPolicy) with(b *oops.OopsErrorBuilder) *oops.OopsErrorBuilder {
	for _, item := range p.items {
		b = item.with(b)
	}
	return b
}
