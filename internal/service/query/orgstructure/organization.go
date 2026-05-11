package orgstructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/cursor"
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
	ID        uuid.UUID
	Name      string
	UpdatedAt time.Time
}

// OrganizationListResult is returned by List and Search. NextCursor is
// nil when no more pages remain.
type OrganizationListResult struct {
	Items      []OrganizationListItem
	NextCursor *string
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

// List returns up to q.Limit organizations, ordered most-recently-updated first.
//
// See: docs/services/OrgStructure.md
func (r *OrganizationReader) List(ctx context.Context, q ListQuery) (OrganizationListResult, error) {
	if err := q.normalize(); err != nil {
		return OrganizationListResult{}, err
	}
	sqlBuf := `SELECT id, name, updated_at FROM projections.organizations`
	args := make([]any, 0, 3)
	if q.After != nil {
		t, idStr, err := cursor.Decode(*q.After)
		if err != nil {
			return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		id, _ := uuid.Parse(idStr)
		sqlBuf += ` WHERE (updated_at, id) < (?, ?)`
		args = append(args, t, id)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]OrganizationListItem, 0, q.Limit+1)
	for rows.Next() {
		var v OrganizationListItem
		if err := rows.Scan(&v.ID, &v.Name, &v.UpdatedAt); err != nil {
			return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	return buildListResult(out, q.Limit), nil
}

// Search returns organizations whose name contains the given substring
// (case-insensitive, ILIKE %query%). An empty query (after trimming at
// the handler boundary) degenerates to the same SQL shape as List so
// callers can swap endpoints without reshaping their page model. The
// query is length-capped BEFORE any DB round-trip; the pattern is
// always bound positionally so users cannot inject SQL.
//
// See: docs/services/OrgStructure.md
func (r *OrganizationReader) Search(ctx context.Context, query string, q ListQuery) (OrganizationListResult, error) {
	if len(query) > organizationSearchMaxQueryLength {
		return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationSearchQueryTooLong).
			Public("Search query is too long.").
			With("max_length", organizationSearchMaxQueryLength).
			With("actual_length", len(query)).
			Errorf("search query too long")
	}
	if err := q.normalize(); err != nil {
		return OrganizationListResult{}, err
	}
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 4)
	if query != "" {
		clauses = append(clauses, `COALESCE(name, '') ILIKE ?`)
		args = append(args, "%"+like.EscapePattern(query)+"%")
	}
	if q.After != nil {
		t, idStr, err := cursor.Decode(*q.After)
		if err != nil {
			return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		id, _ := uuid.Parse(idStr)
		clauses = append(clauses, `(updated_at, id) < (?, ?)`)
		args = append(args, t, id)
	}
	sqlBuf := `SELECT id, name, updated_at FROM projections.organizations`
	for i, c := range clauses {
		if i == 0 {
			sqlBuf += ` WHERE ` + c
		} else {
			sqlBuf += ` AND ` + c
		}
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]OrganizationListItem, 0, q.Limit+1)
	for rows.Next() {
		var v OrganizationListItem
		if err := rows.Scan(&v.ID, &v.Name, &v.UpdatedAt); err != nil {
			return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return OrganizationListResult{}, oops.In("reader.orgstructure.organization").
			Code(ErrCodeOrganizationLoadFailed).
			Wrap(err)
	}
	return buildListResult(out, q.Limit), nil
}

// buildListResult trims the over-fetched row and computes the next
// cursor from the last item in the page.
func buildListResult(rows []OrganizationListItem, limit int) OrganizationListResult {
	var nextCursor *string
	if len(rows) > limit {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return OrganizationListResult{Items: rows, NextCursor: nextCursor}
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
