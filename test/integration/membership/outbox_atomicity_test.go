//go:build integration

package membership_integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// TestOutbox_TransactionAtomicity proves that a failure during the
// outbox.events INSERT causes the entire command transaction —
// including the domain-table write — to roll back.
//
// Mechanism: install a Postgres trigger that raises on any insert
// into outbox.events, then run HireEmployee and assert both tables
// are empty on return.
func TestOutbox_TransactionAtomicity(t *testing.T) {
	f := takeFixture(t)

	// Install the fail-on-outbox-insert trigger. Dropped by t.Cleanup.
	trigger := `
		CREATE OR REPLACE FUNCTION test_reject_outbox() RETURNS trigger
			LANGUAGE plpgsql AS $$
			BEGIN
				RAISE EXCEPTION 'test: outbox insert rejected';
			END;
			$$;
		CREATE TRIGGER test_reject_outbox_trigger
			BEFORE INSERT ON outbox.events
			FOR EACH ROW EXECUTE FUNCTION test_reject_outbox();
	`
	require.NoError(t, testDB.Exec(trigger).Error)
	t.Cleanup(func() {
		_ = testDB.Exec(`DROP TRIGGER IF EXISTS test_reject_outbox_trigger ON outbox.events`).Error
		_ = testDB.Exec(`DROP FUNCTION IF EXISTS test_reject_outbox()`).Error
	})

	// Attempt to hire — the domain INSERT succeeds, but the outbox
	// INSERT inside the same transaction raises, rolling everything back.
	_, err := empSvc.Hire(ctxT(t), membership.HireEmployeeCommand{
		ZitadelUserID: testUserAliceID,
		DepartmentID:  f.DeptA1a,
	})
	require.Error(t, err)

	var empCount int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM domain.employees`).Scan(&empCount).Error)
	assert.Equal(t, int64(0), empCount, "employee insert must have rolled back")

	var outboxCount int64
	require.NoError(t, testDB.Raw(`SELECT count(*) FROM outbox.events`).Scan(&outboxCount).Error)
	assert.Equal(t, int64(0), outboxCount, "outbox insert must have rolled back")
}
