package orgstructure

import "github.com/samber/oops"

// Error codes emitted by pagination validators. Readers reuse the same
// constants so callers can match on a single code no matter which list
// endpoint rejected the request.
const (
	ErrCodeListLimitOutOfRange  = "list_limit_out_of_range"
	ErrCodeListOffsetOutOfRange = "list_offset_out_of_range"
)

// ListQuery is the input shared by every List method. Zero Limit maps
// to listDefaultLimit. Non-zero Limit outside [listMinLimit,
// listMaxLimit] is rejected with list_limit_out_of_range. Offset must
// be non-negative.
type ListQuery struct {
	Limit  int
	Offset int
}

// normalize validates and normalizes the pagination fields in place.
func (q *ListQuery) normalize() error {
	if q.Offset < 0 {
		return oops.In("reader.orgstructure").
			Code(ErrCodeListOffsetOutOfRange).
			Public("List offset must be non-negative.").
			With("field", "offset").
			With("actual_value", q.Offset).
			Errorf("offset out of range")
	}
	if q.Limit == 0 {
		q.Limit = listDefaultLimit
		return nil
	}
	if q.Limit < listMinLimit || q.Limit > listMaxLimit {
		return oops.In("reader.orgstructure").
			Code(ErrCodeListLimitOutOfRange).
			Public("List limit is out of range.").
			With("field", "limit").
			With("actual_value", q.Limit).
			With("min_value", listMinLimit).
			With("max_value", listMaxLimit).
			Errorf("limit out of range")
	}
	return nil
}
