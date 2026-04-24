// Package classifier exposes read methods over the incident classifier
// projections (projections.incident_categories,
// projections.incident_types). Subtree walks use PostgreSQL recursive
// CTEs in raw SQL.
//
// Authorization model: catalog reads (Get / List / Subtree) are gated
// by authz.ReaderOf.{Category,IncidentType,Organization} — any employee
// of the owning organization may read. Patient-facing endpoints
// (ListPatientAllowed*, ListPatientVisible*) run under authz.Authenticated:
// patients are not employees, but they must pick an organization to
// file an incident against, so the classifier menu stays open to any
// logged-in caller. Cross-org reads by employees fail permission_denied.
package classifier

import (
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Bounds applied to List pagination.
const (
	listMinLimit     = 1
	listMaxLimit     = 500
	listDefaultLimit = 50
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
		return oops.In("reader.incident.classifier").
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
		return oops.In("reader.incident.classifier").
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

// Reader exposes read methods for both aggregates in this domain.
type Reader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given db and authorization
// service.
func NewReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, authz: az, logger: logger}
}
