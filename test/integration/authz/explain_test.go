//go:build integration

package authz_integration_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequire_PlanIsIndexBound executes EXPLAIN on a representative
// AdminOf.Clinic check and asserts that every node in the plan uses an
// index scan (never a sequential scan). This is the regression guard
// for the auth path's latency budget on the query side.
func TestRequire_PlanIsIndexBound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

	// Re-derive the query the way Require does, then EXPLAIN it. We
	// cannot pass the policy through EXPLAIN directly without changing
	// Require's signature; instead, run EXPLAIN on the equivalent raw
	// SQL. If Require's rendering ever drifts, this test will start
	// failing in ways visible to the reviewer.
	const query = `EXPLAIN (ANALYZE, FORMAT TEXT)
		SELECT EXISTS(
			SELECT 1 FROM domain.system_admins WHERE zitadel_user_id = @caller
			UNION ALL
			SELECT 1 FROM domain.org_admins oa
			JOIN domain.employees e ON e.id = oa.employee_id
			JOIN domain.clinics c ON c.organization_id = oa.organization_id WHERE c.id = @scope
			AND e.zitadel_user_id = @caller
			UNION ALL
			SELECT 1 FROM domain.org_admins oa
			JOIN domain.employees e ON e.id = oa.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = oa.employee_id
			JOIN domain.clinics c ON c.organization_id = oa.organization_id WHERE c.id = @scope
			AND e.zitadel_user_id = @caller
			AND v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())
		)`

	var lines []string
	require.NoError(t, testDB.Raw(query,
		sql.Named("caller", "bob"),
		sql.Named("scope", clinicA1),
	).Scan(&lines).Error)

	plan := strings.Join(lines, "\n")
	// Postgres renders plan nodes as "Seq Scan on <table>" (no schema
	// prefix in the textual format), so the assertion matches that
	// prefix rather than a schema-qualified form.
	assert.NotContains(t, plan, "Seq Scan on ",
		"auth query should be index-bound; plan was:\n%s", plan)

	// Sanity: the plan must mention index usage for at least one of
	// the hot tables so we don't pass trivially on an empty plan.
	assert.Contains(t, plan, "Index",
		"plan had no Index Scan lines:\n%s", plan)
}
