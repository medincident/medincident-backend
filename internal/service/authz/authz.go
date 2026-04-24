// Package authz is the authorization service. Each Require* method
// answers one question — "is this caller allowed to act on this
// scope?" — with a single round-trip raw SQL query that combines
// system-admin membership, direct role membership, and deputy
// membership (the latter gated by an active vacation on the holder).
//
// The queries are collapsed into UNION ALL branches so that an
// unauthorized caller cannot distinguish "scope does not exist" from
// "you don't have access": for a non-system-admin caller, both cases
// end in the same zero-row EXISTS result and surface as
// permission_denied. Only system admins — who are authorized for any
// scope — can observe the service layer's not_found errors downstream.
//
// Policy.describe() strings are English-only and reach clients verbatim
// via the Public message on permission_denied. They are part of the
// external surface — changing or translating them is a breaking change.
package authz

import "gorm.io/gorm"

// Authz is a concrete authorization service. No interface — handlers
// depend on the struct type directly.
type Authz struct {
	db *gorm.DB
}

// New wires Authz with the shared gorm DB.
func New(db *gorm.DB) *Authz {
	return &Authz{db: db}
}

// activeVacationPredicate is the SQL fragment that evaluates to true
// when a vacation row is currently in effect. Shared by every deputy
// branch so the rule stays in one place: "starts_at has passed and
// ends_at is either open-ended or still in the future". Injected
// inline via fmt.Sprintf at query build time — no user input is
// concatenated.
//
// Postgres now() returns TIMESTAMPTZ, and starts_at / ends_at are
// TIMESTAMPTZ columns (see db/migrations/*employee_vacations*), so
// comparisons are timezone-correct. CURRENT_TIMESTAMP is equivalent
// here; avoid expressions like now() AT TIME ZONE 'UTC' or casts to
// timestamp, which would drop the time zone and feed a naive value
// into the comparison.
const activeVacationPredicate = `v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`
