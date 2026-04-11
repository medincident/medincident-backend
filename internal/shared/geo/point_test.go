package geo_test

import (
	"errors"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

func TestNewPointAccepts(t *testing.T) {
	cases := []struct {
		name                string
		longitude, latitude float64
	}{
		{"null_island", 0, 0},
		{"north_pole", 0, 90},
		{"south_pole", 0, -90},
		{"antimeridian_east", 180, 0},
		{"antimeridian_west", -180, 0},
		{"realistic", 30.5234, 50.4501},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p, err := geo.NewPoint(c.longitude, c.latitude)
			require.NoError(t, err)
			require.Equal(t, c.longitude, p.Longitude)
			require.Equal(t, c.latitude, p.Latitude)
		})
	}
}

func TestNewPointRejectsLongitudeOutOfRange(t *testing.T) {
	_, err := geo.NewPoint(180.1, 0)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, geo.ErrCodePointLongitudeOutOfRange, oe.Code())
	require.Equal(t, "longitude", oe.Context()["field"])
}

func TestNewPointRejectsLatitudeOutOfRange(t *testing.T) {
	_, err := geo.NewPoint(0, 90.1)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, geo.ErrCodePointLatitudeOutOfRange, oe.Code())
	require.Equal(t, "latitude", oe.Context()["field"])
}

func TestNewPointCollectsBothErrors(t *testing.T) {
	_, err := geo.NewPoint(-999, 999)
	require.Error(t, err)

	var joined interface{ Unwrap() []error }
	require.True(t, errors.As(err, &joined))
	require.Len(t, joined.Unwrap(), 2)
}
