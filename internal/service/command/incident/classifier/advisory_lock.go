package classifier

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"
)

// ErrCodeIncidentClassifierLockFailed is reported when the Postgres
// transaction-scoped advisory lock call itself fails (driver error,
// connection gone, etc.). It should never surface under normal load.
const ErrCodeIncidentClassifierLockFailed = "incident_classifier_lock_failed"

// lockClassifierOrg serializes every tree-structural mutation (create,
// move, delete, activate/deactivate) for a single organisation inside
// the current transaction. Without this serialisation, two concurrent
// Move operations on unrelated categories can race the depth / cycle
// checks and violate incidentClassifierMaxDepth.
//
// The lock is a Postgres transaction-scoped advisory lock keyed on
// hashtext("classifier:" + organizationID), automatically released on
// COMMIT or ROLLBACK. Different organisations hash to different keys,
// so cross-org mutations still run in parallel.
//
// Call this after resolving the target row's organization_id and
// before doing any depth / cycle / subtree scans.
func lockClassifierOrg(tx *gorm.DB, organizationID uuid.UUID) error {
	key := "classifier:" + organizationID.String()
	if err := tx.Exec(`SELECT pg_advisory_xact_lock(hashtext(?))`, key).Error; err != nil {
		return oops.In("services.incident.classifier").
			Code(ErrCodeIncidentClassifierLockFailed).
			With("organization_id", organizationID).
			Wrap(err)
	}
	return nil
}
