package membership

import (
	"strings"
	"unicode/utf8"

	"github.com/samber/oops"
)

// Position length limits.
const (
	positionMinLen = 2
	positionMaxLen = 256
)

// validatePosition enforces min/max length on a trimmed position.
// nil is explicitly allowed (= not set); empty-after-trim is rejected.
func validatePosition(in *string) error {
	if in == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*in)
	if trimmed == "" {
		return oops.In("services.membership.employee").
			Code(ErrCodeEmployeePositionTooShort).
			Public("Position is too short.").
			With("position_length", 0).
			With("min_length", positionMinLen).
			Errorf("position too short")
	}
	n := utf8.RuneCountInString(trimmed)
	if n < positionMinLen {
		return oops.In("services.membership.employee").
			Code(ErrCodeEmployeePositionTooShort).
			Public("Position is too short.").
			With("position_length", n).
			With("min_length", positionMinLen).
			Errorf("position too short")
	}
	if n > positionMaxLen {
		return oops.In("services.membership.employee").
			Code(ErrCodeEmployeePositionTooLong).
			Public("Position is too long.").
			With("position_length", n).
			With("max_length", positionMaxLen).
			Errorf("position too long")
	}
	return nil
}
