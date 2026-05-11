package cursor_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-backend/internal/cursor"
)

func TestEncodeDecode_RoundTrip(t *testing.T) {
	ts := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	id := "9f4e2a1b-0000-0000-0000-000000000001"

	s := cursor.Encode(ts, id)
	assert.NotEmpty(t, s)

	got, err := cursor.Decode(s)
	require.NoError(t, err)
	assert.Equal(t, ts.UnixNano(), got.T)
	assert.Equal(t, id, got.I)
	assert.Equal(t, ts, got.Time())
}

func TestDecode_InvalidBase64_ReturnsInvalidCursor(t *testing.T) {
	_, err := cursor.Decode("not!!!valid@@")
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	assert.Equal(t, cursor.ErrCodeInvalidCursor, oe.Code())
}

func TestDecode_InvalidJSON_ReturnsInvalidCursor(t *testing.T) {
	s := base64.RawURLEncoding.EncodeToString([]byte(`{"not":"valid`))
	_, err := cursor.Decode(s)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	assert.Equal(t, cursor.ErrCodeInvalidCursor, oe.Code())
}

func TestDecode_EmptyString_ReturnsInvalidCursor(t *testing.T) {
	_, err := cursor.Decode("")
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	assert.Equal(t, cursor.ErrCodeInvalidCursor, oe.Code())
}

func TestEncode_SameTimeDifferentID_DifferentCursors(t *testing.T) {
	ts := time.Now().UTC()
	assert.NotEqual(t, cursor.Encode(ts, "id-a"), cursor.Encode(ts, "id-b"))
}

func TestEncode_DifferentTime_DifferentCursors(t *testing.T) {
	ts1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ts2 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	assert.NotEqual(t, cursor.Encode(ts1, "same-id"), cursor.Encode(ts2, "same-id"))
}
