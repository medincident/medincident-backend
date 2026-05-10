// Package self exposes read methods that answer "who am I?" for the
// authenticated caller. Every method queries the caller's own data
// only; authorization is handled entirely by the authn interceptor
// upstream — no authz.Require calls are needed here.
//
// See: docs/services/self/Self.md
package self

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// SelfReader exposes read methods scoped to the calling user's own
// identity, memberships, and role assignments.
type SelfReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewSelfReader returns a SelfReader bound to the given database.
func NewSelfReader(db *gorm.DB, logger *zerolog.Logger) *SelfReader {
	return &SelfReader{db: db, logger: logger}
}
