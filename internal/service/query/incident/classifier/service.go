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
//
// See: docs/services/incident/Classifier.md
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
