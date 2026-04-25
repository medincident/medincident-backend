// Package identity owns the Zitadel event projection for user data:
//
//   - Consumer: JetStream pull consumer on the zitadel.> stream;
//     dispatches event_type values to the projector methods.
//   - Projector: UPSERTs against projections.users driven by Zitadel
//     user events. Also back-fills projections.employee_cards
//     denormalised columns when a user profile or email changes.
//
// Aggregate ids coming from Zitadel are opaque strings, so the
// primary key in projections.users is TEXT (see db/migrations/).
package identity

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Projector applies Zitadel user events to the projections schema.
// Methods are transactional; the consumer calls Apply* from a dispatch switch.
type Projector struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewProjector returns a Projector bound to the given db.
func NewProjector(db *gorm.DB, logger *zerolog.Logger) *Projector {
	return &Projector{db: db, logger: logger}
}
