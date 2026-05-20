// Package classifier exposes read methods over the incident classifier
// projections (projections.incident_categories,
// projections.incident_types). Subtree walks use PostgreSQL recursive
// CTEs in raw SQL.
//
// Authorization model: all list/subtree RPCs accept any authenticated
// caller (authz.Authenticated). Employees (authz.ReaderOf.Organization /
// authz.ReaderOf.Category) receive the full unfiltered result set.
// Non-employee authenticated callers (patients) automatically receive
// only active categories whose subtree contains a patient-allowed type,
// and only active, patient-allowed types. The branching is done via
// authz.Satisfies so no permission_denied is raised for patient callers.
// Get methods (GetCategory, GetType) remain employee-only via
// authz.ReaderOf.{Category,IncidentType}. Cross-org reads by employees
// return the full set for the requested org; patients see only the
// patient-visible subset regardless of their own org.
//
// See: docs/services/incident/Classifier.md
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
	ErrCodeListLimitOutOfRange = "list_limit_out_of_range"
	ErrCodeListBadCursor       = "incident_classifier_bad_cursor"
)

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
		return oops.In("reader.incident.classifier").
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

// CategoryListResult is returned by paginated category list methods.
type CategoryListResult struct {
	Items      []CategoryView
	NextCursor *string
}

// TypeListResult is returned by paginated type list methods.
type TypeListResult struct {
	Items      []TypeView
	NextCursor *string
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
