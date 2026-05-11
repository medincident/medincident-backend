package orgstructure

import (
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/service/query"
)

// Error codes emitted by pagination validators.
const (
	ErrCodeListLimitOutOfRange = "list_limit_out_of_range"
	ErrCodeListBadCursor       = "orgstructure_bad_cursor"
)

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
