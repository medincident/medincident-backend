//go:build integration

package authz_integration_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestRequire_PlanIsIndexBound executes EXPLAIN on a representative
// AdminOf.Clinic check and asserts that every plan node uses an index
// scan (never a sequential scan). This is the regression guard for
// the auth path's latency budget on the query side.
//
// The planner is a cost-based optimiser: on a dataset as small as the
// one seedFixtures produces it can legitimately pick Seq Scan because
// the tables fit in one page. That would make the assertion pass for
// the wrong reason (no regression signal) or flake (planner changes
// its mind). Instead, we disable seqscan inside a transaction with
// SET LOCAL. With seqscan disabled, the planner still falls back to
// Seq Scan for relations that have no usable index — so "Seq Scan
// on <table>" in the resulting plan is exactly the signal of an index
// regression we care about, independent of data size.
func TestRequire_PlanIsIndexBound(t *testing.T) {
	resetDB(t)
	seedFixtures(t)

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
	err := testDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`SET LOCAL enable_seqscan = off`).Error; err != nil {
			return err
		}
		return tx.Raw(query,
			sql.Named("caller", "bob"),
			sql.Named("scope", clinicA1),
		).Scan(&lines).Error
	})
	require.NoError(t, err)

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
