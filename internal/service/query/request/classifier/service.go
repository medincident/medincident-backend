// Package classifier exposes read methods over the request type
// projections (projections.request_types).
//
// See: docs/services/request/Classifier.md
package classifier

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query"
)

// Error codes emitted by pagination validators.
const (
	ErrCodeListLimitOutOfRange = "list_limit_out_of_range"
	ErrCodeListBadCursor       = "list_bad_cursor"
)

// createdAtCursor is the keyset pagination token for lists ordered by
// (created_at DESC, id DESC).
type createdAtCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

// nameCursor is the keyset pagination token for lists ordered by
// (name ASC, id ASC).
type nameCursor struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

func encodeCursor[T any](c T) string {
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func decodeCursor[T any](s string) (T, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	var zero T
	if err != nil {
		return zero, err
	}
	var c T
	if err := json.Unmarshal(b, &c); err != nil {
		return zero, err
	}
	return c, nil
}

// ListQuery is the input shared by paginated list endpoints.
type ListQuery struct {
	Limit int
	After *string
}

// normalize validates and normalizes the pagination fields.
func (q *ListQuery) normalize() error {
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

// RequestTypeListResult is returned by paginated request type list methods.
type RequestTypeListResult struct {
	Items      []RequestTypeView
	NextCursor *string
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
