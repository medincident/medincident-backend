package geo

import (
	"strings"
	"unicode/utf8"

	"github.com/samber/oops"
)

// MaxAddressTextLen is the maximum length of an address text in runes.
const MaxAddressTextLen = 500

// Error codes emitted by the Address validators. Each code is specific
// to a single violation and declared in the same file as the validator
// that emits it — per project convention, codes live next to the point
// of application, not in a central codes file. Constant names keep the
// model prefix for unique grep; string values drop it — the emitting
// package (`shared.geo`) and field already convey the model context.
const (
	ErrCodeAddressTextEmpty   = "text_empty"
	ErrCodeAddressTextTooLong = "text_too_long"
)

// Address is an immutable value object. It holds the free-form text of
// a postal address and an optional geographic point. Use NewAddress to
// construct; direct struct literals are only permitted in the trusted
// persistence path. No struct tags — the domain layer does not know
// about serialisation formats.
type Address struct {
	Text  string
	Point *Point
}

// NewAddress is the ONE constructor of Address. It validates text.
// The Point parameter is accepted as trusted (already built via NewPoint
// or explicitly nil).
func NewAddress(text string, point *Point) (Address, error) {
	text = strings.TrimSpace(text)
	if err := validateAddressText(text); err != nil {
		return Address{}, err
	}
	return Address{Text: text, Point: point}, nil
}

func validateAddressText(text string) error {
	if text == "" {
		return oops.In("shared.geo").
			Code(ErrCodeAddressTextEmpty).
			Public("Address text is required.").
			With("field", "text").
			Errorf("address text is empty")
	}
	if n := utf8.RuneCountInString(text); n > MaxAddressTextLen {
		return oops.In("shared.geo").
			Code(ErrCodeAddressTextTooLong).
			Public("Address text is too long.").
			With("field", "text").
			With("actual_length", n).
			With("max_length", MaxAddressTextLen).
			Errorf("address text too long")
	}
	return nil
}
