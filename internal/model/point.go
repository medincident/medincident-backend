package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/oops"
)

// Point is a geographic coordinate pair. Both fields are always valid;
// absence of a point is expressed via Valid=false on the
// null.Value[Point] wrapper in the owning Address struct.
type Point struct {
	Longitude float64
	Latitude  float64
}

// Equal reports whether p and other have identical coordinates.
func (p Point) Equal(other Point) bool {
	return p.Longitude == other.Longitude && p.Latitude == other.Latitude
}

// Value serialises Point into the Postgres composite text format
// (longitude,latitude) accepted by the domain.geo_point type.
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)",
		strconv.FormatFloat(p.Longitude, 'f', -1, 64),
		strconv.FormatFloat(p.Latitude, 'f', -1, 64),
	), nil
}

// Scan parses the Postgres composite text format (longitude,latitude).
func (p *Point) Scan(src any) error {
	if src == nil {
		return oops.Errorf("cannot scan nil into Point")
	}
	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return oops.Errorf("unsupported scan type %T for Point", src)
	}
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(s), ")"), "("))
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return oops.Errorf("expected (lon,lat), got %q", s)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return oops.Wrap(err)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return oops.Wrap(err)
	}
	p.Longitude, p.Latitude = lon, lat
	return nil
}
