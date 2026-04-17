package orgstructure

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/samber/oops"
)

// AddressInput is the command-layer input for an address. Text is
// required (non-empty, bounded). Point is optional; when present, both
// coordinates must be inside their valid ranges.
type AddressInput struct {
	Text  string
	Point *PointInput
}

// PointInput is the command-layer input for optional geographic
// coordinates.
type PointInput struct {
	Longitude float64
	Latitude  float64
}

// Invariant limits for Address text length and coordinate ranges.
const (
	addressMinTextLen = 4
	addressMaxTextLen = 128
	minLongitude      = -180.0
	maxLongitude      = 180.0
	minLatitude       = -90.0
	maxLatitude       = 90.0
)

// Error codes emitted by the Address/Point validators.
const (
	ErrCodeAddressTextEmpty           = "address_text_empty"
	ErrCodeAddressTextTooShort        = "address_text_too_short"
	ErrCodeAddressTextTooLong         = "address_text_too_long"
	ErrCodeAddressLongitudeOutOfRange = "address_longitude_out_of_range"
	ErrCodeAddressLatitudeOutOfRange  = "address_latitude_out_of_range"
)

// validateAddressInput checks that addr.Text is present and within the
// bounded range and, when a Point is provided, delegates to
// validatePointInput. All violations are collected into a single
// multi-error so the caller sees every field error at once.
func validateAddressInput(addr AddressInput) error {
	var errs []error
	text := strings.TrimSpace(addr.Text)
	if text == "" {
		errs = append(errs, oops.In("services.orgstructure.address").
			Code(ErrCodeAddressTextEmpty).
			Public("Address text is required.").
			With("field", "text").
			Errorf("address text is empty"))
	} else {
		n := utf8.RuneCountInString(text)
		if n < addressMinTextLen {
			errs = append(errs, oops.In("services.orgstructure.address").
				Code(ErrCodeAddressTextTooShort).
				Public("Address text is too short.").
				With("field", "text").
				With("actual_length", n).
				With("min_length", addressMinTextLen).
				Errorf("address text too short"))
		}
		if n > addressMaxTextLen {
			errs = append(errs, oops.In("services.orgstructure.address").
				Code(ErrCodeAddressTextTooLong).
				Public("Address text is too long.").
				With("field", "text").
				With("actual_length", n).
				With("max_length", addressMaxTextLen).
				Errorf("address text too long"))
		}
	}
	if addr.Point != nil {
		if err := validatePointInput(*addr.Point); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// validatePointInput checks that BOTH longitude and latitude are in
// range; violations are collected and joined so callers see both field
// errors at once. The nil check is handled by validateAddressInput
// before calling this function.
func validatePointInput(p PointInput) error {
	var errs []error
	if p.Longitude < minLongitude || p.Longitude > maxLongitude {
		errs = append(errs, oops.In("services.orgstructure.address").
			Code(ErrCodeAddressLongitudeOutOfRange).
			Public("Longitude is out of range.").
			With("field", "longitude").
			With("actual_value", p.Longitude).
			With("min_value", minLongitude).
			With("max_value", maxLongitude).
			Errorf("longitude out of range"))
	}
	if p.Latitude < minLatitude || p.Latitude > maxLatitude {
		errs = append(errs, oops.In("services.orgstructure.address").
			Code(ErrCodeAddressLatitudeOutOfRange).
			Public("Latitude is out of range.").
			With("field", "latitude").
			With("actual_value", p.Latitude).
			With("min_value", minLatitude).
			With("max_value", maxLatitude).
			Errorf("latitude out of range"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
