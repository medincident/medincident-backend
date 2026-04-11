package geo_test

import (
	"strings"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

func TestNewAddressValid(t *testing.T) {
	p, err := geo.NewPoint(30.5, 50.5)
	require.NoError(t, err)

	addr, err := geo.NewAddress("  Main St 1  ", &p)
	require.NoError(t, err)
	require.Equal(t, "Main St 1", addr.Text)
	require.NotNil(t, addr.Point)
	require.Equal(t, p, *addr.Point)
}

func TestNewAddressAcceptsNilPoint(t *testing.T) {
	addr, err := geo.NewAddress("Main St 1", nil)
	require.NoError(t, err)
	require.Equal(t, "Main St 1", addr.Text)
	require.Nil(t, addr.Point)
}

func TestNewAddressRejectsEmptyText(t *testing.T) {
	_, err := geo.NewAddress("   ", nil)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, geo.ErrCodeAddressTextEmpty, oe.Code())
}

func TestNewAddressRejectsTextTooLong(t *testing.T) {
	long := strings.Repeat("x", geo.MaxAddressTextLen+1)
	_, err := geo.NewAddress(long, nil)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, geo.ErrCodeAddressTextTooLong, oe.Code())
}
