package orgstructure

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/service/query"
)

// Error codes emitted by pagination validators.
const (
	ErrCodeListLimitOutOfRange = "list_limit_out_of_range"
	ErrCodeListBadCursor       = "list_bad_cursor"
)

// cursor is the keyset pagination token shared by all list endpoints
// in this package. Entries are ordered by (created_at DESC, id DESC).
type cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

func encodeCursor(c cursor) string {
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func decodeCursor(s string) (cursor, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return cursor{}, err
	}
	var c cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return cursor{}, err
	}
	return c, nil
}

// ListQuery is the input shared by every List method. Zero Limit maps
// to query.DefaultLimit. Non-zero Limit outside [query.MinLimit,
// query.MaxLimit] is rejected with list_limit_out_of_range. After is
// the opaque cursor from a previous response's next_cursor field.
type ListQuery struct {
	Limit int
	After *string
}

// normalize validates and normalizes the pagination fields in place.
func (q *ListQuery) normalize() error {
	if q.Limit == 0 {
		q.Limit = query.DefaultLimit
		return nil
	}
	if q.Limit < query.MinLimit || q.Limit > query.MaxLimit {
		return oops.In("reader.orgstructure").
			Code(ErrCodeListLimitOutOfRange).
			Public("List limit is out of range.").
			With("field", "limit").
			With("actual_value", q.Limit).
			With("min_value", query.MinLimit).
			With("max_value", query.MaxLimit).
			Errorf("limit out of range")
	}
	return nil
}
