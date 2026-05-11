package cursor

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/samber/oops"
)

const ErrCodeInvalidCursor = "invalid_cursor"

// Cursor encodes the pagination position as (updated_at, id). Both fields
// are needed because updated_at is not unique — identical timestamps exist
// when rows are bulk-inserted.
type Cursor struct {
	T int64  `json:"t"` // UnixNano of updated_at
	I string `json:"i"` // row id (UUID string or Zitadel user id)
}

// Encode serialises a cursor to base64url(JSON). The result is URL-safe,
// deterministic, and carries no secrets.
func Encode(updatedAt time.Time, id string) string {
	b, _ := json.Marshal(Cursor{T: updatedAt.UnixNano(), I: id})
	return base64.RawURLEncoding.EncodeToString(b)
}

// Decode parses a cursor produced by Encode. Callers must check
// after != "" before calling Decode — an empty string is never valid.
// Returns an oops error with code ErrCodeInvalidCursor on any failure.
func Decode(s string) (Cursor, error) {
	if s == "" {
		return Cursor{}, oops.In("cursor").
			Code(ErrCodeInvalidCursor).
			Public("Invalid pagination cursor.").
			Errorf("empty cursor")
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Cursor{}, oops.In("cursor").
			Code(ErrCodeInvalidCursor).
			Public("Invalid pagination cursor.").
			Wrap(err)
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return Cursor{}, oops.In("cursor").
			Code(ErrCodeInvalidCursor).
			Public("Invalid pagination cursor.").
			Wrap(err)
	}
	if c.T == 0 || c.I == "" {
		return Cursor{}, oops.In("cursor").
			Code(ErrCodeInvalidCursor).
			Public("Invalid pagination cursor.").
			Errorf("cursor missing required fields")
	}
	return c, nil
}

// Time converts the stored UnixNano back to a UTC time.Time.
func (c Cursor) Time() time.Time {
	return time.Unix(0, c.T).UTC()
}
