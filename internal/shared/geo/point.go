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

// Error codes emitted by the Point validators. Specific per violation
// and declared next to the emitting validators. Constant names keep the
// model prefix so they stay unique across the codebase; the string
// values drop it because the emitting package (`shared.geo`) and field
// (`longitude` / `latitude`) already convey the model context.
const (
	ErrCodePointLongitudeOutOfRange = "longitude_out_of_range"
	ErrCodePointLatitudeOutOfRange  = "latitude_out_of_range"
)

// Point is an immutable geographic point. It is created only via NewPoint.
// Exported fields enable trusted struct-literal construction in the
// persistence path. No struct tags — the domain layer does not know about
// serialisation formats; the outbox layer calls encoding/json on these
// types using their exported field names directly, which is acceptable
// because the domain type is the sole source of truth for both writers
// and readers of outbox JSONB payloads.
type Point struct {
	Longitude float64
	Latitude  float64
}

// NewPoint is the ONE constructor of Point. It validates both coordinates
// and collects errors via errors.Join (multi-error inside this single VO).
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
