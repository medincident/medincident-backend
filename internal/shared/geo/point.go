package geo

import (
	"errors"

	"github.com/samber/oops"
)

// Geographic coordinate bounds (WGS84).
const (
	MinLongitude = -180.0
	MaxLongitude = 180.0
	MinLatitude  = -90.0
	MaxLatitude  = 90.0
)

// Error codes emitted by the Point validators.
const (
	ErrCodePointLongitudeOutOfRange = "longitude_out_of_range"
	ErrCodePointLatitudeOutOfRange  = "latitude_out_of_range"
)

// Point is an immutable geographic point in WGS84 decimal degrees.
type Point struct {
	Longitude float64
	Latitude  float64
}

// NewPoint validates both coordinates and collects failures via
// errors.Join so a single call surfaces every violation at once.
func NewPoint(longitude, latitude float64) (Point, error) {
	var errs []error
	if err := validateLongitude(longitude); err != nil {
		errs = append(errs, err)
	}
	if err := validateLatitude(latitude); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return Point{}, errors.Join(errs...)
	}
	return Point{Longitude: longitude, Latitude: latitude}, nil
}

// Equal reports whether p and other represent the same geographic point.
func (p Point) Equal(other Point) bool {
	return p.Longitude == other.Longitude && p.Latitude == other.Latitude
}

func validateLongitude(longitude float64) error {
	if longitude < MinLongitude || longitude > MaxLongitude {
		return oops.In("shared.geo").
			Code(ErrCodePointLongitudeOutOfRange).
			Public("Longitude is out of range.").
			With("field", "longitude").
			With("actual_value", longitude).
			With("min_value", MinLongitude).
			With("max_value", MaxLongitude).
			Errorf("longitude out of range")
	}
	return nil
}

func validateLatitude(latitude float64) error {
	if latitude < MinLatitude || latitude > MaxLatitude {
		return oops.In("shared.geo").
			Code(ErrCodePointLatitudeOutOfRange).
			Public("Latitude is out of range.").
			With("field", "latitude").
			With("actual_value", latitude).
			With("min_value", MinLatitude).
			With("max_value", MaxLatitude).
			Errorf("latitude out of range")
	}
	return nil
}
