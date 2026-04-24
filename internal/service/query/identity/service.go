// Package identity owns the Zitadel passthrough read surface:
//
//   - Consumer: JetStream pull consumer on the zitadel.> stream;
//     dispatches event_type values to the projector methods below.
//   - Projector: UPSERTs against projections.users and
//     projections.sessions driven by Zitadel user/session events.
//     Also back-fills projections.employee_cards denormalised columns
//     when a user profile/email changes.
//   - Reader: GetUser + GetSession over projections.users /
//     projections.sessions.
//
// Aggregate ids coming from Zitadel are opaque strings, so primary
// keys in projections.users and projections.sessions are TEXT (see
// the migrations under db/migrations/).
package identity

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Projector applies Zitadel user / session events to the
// projections schema. Methods are transactional; the consumer calls
// Apply* from a dispatch switch.
type Projector struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewProjector returns a Projector bound to the given db.
func NewProjector(db *gorm.DB, logger *zerolog.Logger) *Projector {
	return &Projector{db: db, logger: logger}
}

// Reader exposes read methods for the identity projections.
//
// Authorization model: identity rows describe the caller or their
// sessions, so reads are gated by self-or-sysadmin. GetUser compares
// the caller's Zitadel ID to the target id in Go (both are string
// identifiers); mismatches fall back to authz.SystemAdmin. GetSession
// goes through authz with AnyOf(SystemAdmin, SelfSession) so the
// ownership check is performed in the same EXISTS round trip and
// non-owners cannot enumerate session IDs.
type Reader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given db and authorization
// service.
func NewReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, authz: az, logger: logger}
}
