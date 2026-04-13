package membership

import (
	"strings"
	"unicode/utf8"

	"github.com/guregu/null/v6"
	"github.com/samber/oops"
)

// Position length limits.
const (
	positionMinLen = 2
	positionMaxLen = 256
)

// normalisePosition turns a free-form optional position input into a
// null.String ready for persistence:
//   - nil input → null.String{} (not set)
//   - non-nil but trims to empty → null.String{} (treated as cleared)
//   - otherwise → trimmed value
//
// Assumes validatePosition has already been called.
func normalisePosition(in *string) null.String {
	if in == nil {
		return null.String{}
	}
	trimmed := strings.TrimSpace(*in)
	if trimmed == "" {
		return null.String{}
	}
	return null.StringFrom(trimmed)
}

// validatePosition enforces min/max length on a trimmed position.
// nil and trim-to-empty are explicitly allowed (= not set).
func validatePosition(in *string) error {
	if in == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*in)
	if trimmed == "" {
		return nil
	}
	n := utf8.RuneCountInString(trimmed)
	if n < positionMinLen {
		return oops.In("services.membership.employee").
			Code(ErrCodeEmployeePositionTooShort).
			Public("Position is too short.").
			With("position_length", n).
			With("min_length", positionMinLen).
			Errorf(ErrCodeEmployeePositionTooShort)
	}
	if n > positionMaxLen {
		return oops.In("services.membership.employee").
			Code(ErrCodeEmployeePositionTooLong).
			Public("Position is too long.").
			With("position_length", n).
			With("max_length", positionMaxLen).
			Errorf(ErrCodeEmployeePositionTooLong)
	}
	return nil
}
