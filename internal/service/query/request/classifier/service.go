// Package classifier exposes read methods over the request type
// projections (projections.request_types).
//
// See: docs/services/request/Classifier.md
package classifier

import (
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query"
)

// Error codes emitted by pagination validators.
const (
	ErrCodeListLimitOutOfRange  = "list_limit_out_of_range"
	ErrCodeListOffsetOutOfRange = "list_offset_out_of_range"
)

// ListQuery is the input shared by paginated list endpoints.
type ListQuery struct {
	Limit  int
	Offset int
}

// normalize validates and normalizes the pagination fields.
func (q *ListQuery) normalize() error {
	if q.Offset < 0 {
		return oops.In("reader.request.classifier").
			Code(ErrCodeListOffsetOutOfRange).
			Public("List offset must be non-negative.").
			With("field", "offset").
			With("actual_value", q.Offset).
			Errorf("offset out of range")
	}
	if q.Limit == 0 {
		q.Limit = query.DefaultLimit
		return nil
	}
	if q.Limit < query.MinLimit || q.Limit > query.MaxLimit {
		return oops.In("reader.request.classifier").
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

// Reader exposes read methods for request types.
type Reader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given db and authorization service.
func NewReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, authz: az, logger: logger}
}
