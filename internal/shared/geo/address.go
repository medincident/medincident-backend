package geo

import (
	"strings"
	"unicode/utf8"

	"github.com/samber/oops"
)

// MaxAddressTextLen is the maximum length of an address text in runes.
const MaxAddressTextLen = 500

// Error codes emitted by the Address validators.
const (
	ErrCodeAddressTextEmpty   = "text_empty"
	ErrCodeAddressTextTooLong = "text_too_long"
)

// Address is an immutable value object carrying the free-form postal
// text and an optional geographic point.
type Address struct {
	Text  string
	Point *Point
}

// NewAddress validates text and constructs an Address. The Point
// parameter is accepted as trusted (built via NewPoint or explicitly nil).
func NewAddress(text string, point *Point) (Address, error) {
	text = strings.TrimSpace(text)
	if err := validateAddressText(text); err != nil {
		return Address{}, err
	}
	return Address{Text: text, Point: point}, nil
}

// Equal reports whether a and other represent the same address.
func (a Address) Equal(other Address) bool {
	if a.Text != other.Text {
		return false
	}
	if a.Point == nil && other.Point == nil {
		return true
	}
	if a.Point == nil || other.Point == nil {
		return false
	}
	return a.Point.Equal(*other.Point)
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
