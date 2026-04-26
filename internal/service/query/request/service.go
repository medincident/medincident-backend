// Package request is the query-side reader for projections.service_requests
// and the two history projections.
//
// See: docs/services/request/Requests.md
package request

import (
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
	ErrCodeListOffsetOutOfRange     = "service_request_list_offset_out_of_range"
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

// ListQuery captures pagination parameters.
type ListQuery struct {
	Limit  int
	Offset int
}

func (q *ListQuery) normalize() error {
	if q.Offset < 0 {
		return oops.In(scope).
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

func wrapRead(err error, action string) error {
	return oops.In(scope).Code(ErrCodeServiceRequestReadFailed).
		With("action", action).Wrap(err)
}
