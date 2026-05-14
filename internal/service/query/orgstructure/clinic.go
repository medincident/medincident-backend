package orgstructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/cursor"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Error codes emitted by ClinicReader methods.
const (
	ErrCodeClinicNotFound    = "clinic_not_found"
	ErrCodeClinicLoadFailed  = "clinic_load_failed"
	ErrCodeClinicCountFailed = "clinic_count_failed"
)

// ClinicDetails is the full card of a clinic returned by Get.
type ClinicDetails struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Description     *string
	PhysicalAddress AddressView
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ClinicListItem is the minimal view returned by paginated list
// endpoints.
type ClinicListItem struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	IsActive       bool
	UpdatedAt      time.Time
}

// ClinicListResult is returned by ListByOrganization. NextCursor is nil
// when no more pages remain.
type ClinicListResult struct {
	Items      []ClinicListItem
	NextCursor *string
}

// Get returns the ClinicDetails for the given id. Authorization:
// authz.ReaderOf.Clinic(id) — system admin, organization admin of the
// owning org, or any employee of the owning org.
//
// See: docs/services/OrgStructure.md
func (r *ClinicReader) Get(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*ClinicDetails, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(id)); err != nil {
		return nil, err
	}
	var (
		out      ClinicDetails
		addrText string
		lon, lat *float64
	)
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, organization_id, name, description,
		       physical_address_text, physical_address_longitude, physical_address_latitude,
		       is_active, created_at, updated_at
		  FROM projections.clinics
		 WHERE id = ?`, id,
	).Row().Scan(
		&out.ID, &out.OrganizationID, &out.Name, &out.Description,
		&addrText, &lon, &lat,
		&out.IsActive, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.orgstructure.clinic").
				Code(ErrCodeClinicNotFound).
				Public("Clinic not found.").
				With("clinic_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.orgstructure.clinic").
			Code(ErrCodeClinicLoadFailed).
			With("clinic_id", id).
			Wrap(err)
	}
	out.PhysicalAddress = AddressView{Text: addrText}
	if lon != nil && lat != nil {
		out.PhysicalAddress.Point = &PointView{Longitude: *lon, Latitude: *lat}
	}
	return &out, nil
}

// CountByOrganization returns the total number of clinics belonging
// to the given organization. Authorization:
// authz.ReaderOf.Organization(organizationID).
//
// See: docs/services/OrgStructure.md
func (r *ClinicReader) CountByOrganization(
	ctx context.Context,
	caller authz.Caller,
	organizationID uuid.UUID,
) (int64, error) {
	if err := r.authz.Require(
		ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(organizationID),
	); err != nil {
		return 0, err
	}
	var total int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.clinics WHERE organization_id = ?`,
		organizationID,
	).Row().Scan(&total); err != nil {
		return 0, oops.In("reader.orgstructure.clinic").
			Code(ErrCodeClinicCountFailed).
			With("organization_id", organizationID).
			Wrap(err)
	}
	return total, nil
}

// ListByOrganization returns up to q.Limit clinics belonging to the
// given organization, ordered most-recently-updated first. Authorization:
// authz.ReaderOf.Organization(organizationID). Pagination bounds are
// normalized first so a malformed Limit cannot trigger a gratuitous
// authz DB round-trip — matching the validate→authorize order used on
// the command side.
//
// See: docs/services/OrgStructure.md
func (r *ClinicReader) ListByOrganization(
	ctx context.Context,
	caller authz.Caller,
	organizationID uuid.UUID,
	q ListQuery,
) (ClinicListResult, error) {
	if err := q.normalize(); err != nil {
		return ClinicListResult{}, err
	}
	if err := r.authz.Require(
		ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(organizationID),
	); err != nil {
		return ClinicListResult{}, err
	}
	sqlBuf := `SELECT id, organization_id, name, is_active, updated_at
		  FROM projections.clinics
		 WHERE organization_id = ?`
	args := make([]any, 0, 4)
	args = append(args, organizationID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return ClinicListResult{}, oops.In("reader.orgstructure.clinic").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return ClinicListResult{}, oops.In("reader.orgstructure.clinic").
			Code(ErrCodeClinicLoadFailed).
			With("organization_id", organizationID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]ClinicListItem, 0, q.Limit+1)
	for rows.Next() {
		var v ClinicListItem
		if err := rows.Scan(&v.ID, &v.OrganizationID, &v.Name, &v.IsActive, &v.UpdatedAt); err != nil {
			return ClinicListResult{}, oops.In("reader.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("organization_id", organizationID).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return ClinicListResult{}, oops.In("reader.orgstructure.clinic").
			Code(ErrCodeClinicLoadFailed).
			With("organization_id", organizationID).
			Wrap(err)
	}

	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return ClinicListResult{Items: out, NextCursor: nextCursor}, nil
}
