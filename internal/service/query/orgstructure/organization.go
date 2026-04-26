package orgstructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/util/like"
)

// Error codes emitted by OrganizationReader methods.
const (
	ErrCodeOrganizationNotFound           = "organization_not_found"
	ErrCodeOrganizationLoadFailed         = "organization_load_failed"
	ErrCodeOrganizationCountFailed        = "organization_count_failed"
	ErrCodeOrganizationSearchQueryTooLong = "organization_search_query_too_long"
)

// organizationSearchMaxQueryLength caps the user-supplied query for
// Search. Over-long inputs are rejected before any DB round-trip so a
// multi-megabyte pattern cannot tie up the projection.
const organizationSearchMaxQueryLength = 256

// PointView is the optional coordinate pair attached to an AddressView.
type PointView struct {
	Longitude float64
	Latitude  float64
}

// AddressView is the read-side shape of a stored address. Point is nil
// when the backing row has no longitude/latitude.
type AddressView struct {
	Text  string
	Point *PointView
}

// OrganizationDetails is the full card of an organization returned by
// Get. Timestamps are NOT NULL in the projection so they are plain
// time.Time values, not pointers.
type OrganizationDetails struct {
	ID           uuid.UUID
	Name         string
	Description  *string
	LegalAddress AddressView
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// OrganizationListItem is the minimal view returned by paginated list
// endpoints.
type OrganizationListItem struct {
	ID   uuid.UUID
	Name string
}

// Get returns the OrganizationDetails for the given id. Returns an
// oops error with code organization_not_found when the row is absent.
//
// See: docs/services/OrgStructure.md
func (r *OrganizationReader) Get(ctx context.Context, id uuid.UUID) (*OrganizationDetails, error) {
	var (
		out      OrganizationDetails
		addrText string
		lon, lat *float64
	)
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, name, description,
		       legal_address_text, legal_address_longitude, legal_address_latitude,
		       created_at, updated_at
		  FROM projections.organizations
		 WHERE id = ?`, id,
	).Row().Scan(
		&out.ID, &out.Name, &out.Description,
		&addrText, &lon, &lat,
		&out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.orgstructure.organization").
				Code(ErrCodeOrganizationNotFound).
				Public("Organization not found.").
				With("organization_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			With("organization_id", id).
			Wrap(err)
	}
	out.LegalAddress = AddressView{Text: addrText}
	if lon != nil && lat != nil {
		out.LegalAddress.Point = &PointView{Longitude: *lon, Latitude: *lat}
	}
	return &out, nil
}

// List returns up to q.Limit organizations, ordered most-recently-created first.
//
// See: docs/services/OrgStructure.md
func (r *OrganizationReader) List(ctx context.Context, q ListQuery) ([]OrganizationListItem, error) {
	if err := q.normalize(); err != nil {
		return nil, err
	}
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT id, name
		  FROM projections.organizations
		 ORDER BY created_at DESC, id DESC
		 LIMIT ? OFFSET ?`, q.Limit, q.Offset,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]OrganizationListItem, 0, q.Limit)
	for rows.Next() {
		var v OrganizationListItem
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			return nil, oops.In("reader.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	return out, nil
}

// Search returns organizations whose name contains the given substring
// (case-insensitive, ILIKE %query%). An empty query (after trimming at
// the handler boundary) degenerates to the same SQL shape as List so
// callers can swap endpoints without reshaping their page model. The
// query is length-capped BEFORE any DB round-trip; the pattern is
// always bound positionally so users cannot inject SQL.
//
// See: docs/services/OrgStructure.md
func (r *OrganizationReader) Search(ctx context.Context, query string, q ListQuery) ([]OrganizationListItem, error) {
	if len(query) > organizationSearchMaxQueryLength {
		return nil, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationSearchQueryTooLong).
			Public("Search query is too long.").
			With("max_length", organizationSearchMaxQueryLength).
			With("actual_length", len(query)).
			Errorf("search query too long")
	}
	if err := q.normalize(); err != nil {
		return nil, err
	}
	sqlBuf := `SELECT id, name FROM projections.organizations`
	args := make([]any, 0, 3)
	if query != "" {
		sqlBuf += ` WHERE COALESCE(name, '') ILIKE ?`
		args = append(args, "%"+like.EscapePattern(query)+"%")
	}
	sqlBuf += ` ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`
	args = append(args, q.Limit, q.Offset)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return nil, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]OrganizationListItem, 0, q.Limit)
	for rows.Next() {
		var v OrganizationListItem
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			return nil, oops.In("reader.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	return out, nil
}

// Count returns the total number of organizations in the projection.
//
// See: docs/services/OrgStructure.md
func (r *OrganizationReader) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.organizations`,
	).Row().Scan(&total); err != nil {
		return 0, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationCountFailed).
			Wrap(err)
	}
	return total, nil
}
