// Package request is the query-side reader for projections.service_requests
// and the two history projections.
//
// See: docs/services/request/Requests.md
package request

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

const (
	ErrCodeServiceRequestReadFailed = "service_request_query_read_failed"
	ErrCodeServiceRequestNotFound   = "service_request_query_not_found"
	ErrCodeIncidentNotFound         = "service_request_query_incident_not_found"
	ErrCodeListLimitOutOfRange      = "service_request_list_limit_out_of_range"
	ErrCodeListBadCursor            = "service_request_list_bad_cursor"
)

const scope = "services.query.request"

// Reader is the query-side service for service requests.
type Reader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given db, authz, and logger.
func NewReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, authz: az, logger: logger}
}

// createdAtCursor is the keyset pagination token for lists ordered by
// (created_at DESC, id DESC).
type createdAtCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
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

// ListQuery captures pagination parameters.
type ListQuery struct {
	Limit int
	After *string
}

func (q *ListQuery) normalize() error {
	if q.Limit == 0 {
		q.Limit = query.DefaultLimit
		return nil
	}
	if q.Limit < query.MinLimit || q.Limit > query.MaxLimit {
		return oops.In(scope).
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

// ServiceRequestListResult is returned by paginated list methods.
type ServiceRequestListResult struct {
	Items      []ServiceRequestView
	NextCursor *string
}

func wrapRead(err error, action string) error {
	return oops.In(scope).Code(ErrCodeServiceRequestReadFailed).
		With("action", action).Wrap(err)
}
