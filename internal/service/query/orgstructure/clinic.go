package orgstructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

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
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ClinicListItem is the minimal view returned by paginated list
// endpoints.
type ClinicListItem struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
}

// Get returns the ClinicDetails for the given id. Authorization:
// authz.ReaderOf.Clinic(id) — system admin, organization admin of the
// owning org, or any employee of the owning org.
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
		       created_at, updated_at
		  FROM projections.clinics
		 WHERE id = ?`, id,
	).Row().Scan(
		&out.ID, &out.OrganizationID, &out.Name, &out.Description,
		&addrText, &lon, &lat,
		&out.CreatedAt, &out.UpdatedAt,
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

// ListByOrganization returns up to q.Limit clinics belonging to the
// given organization, ordered most-recently-created first. Authorization:
// authz.ReaderOf.Organization(organizationID).
func (r *ClinicReader) ListByOrganization(
	ctx context.Context,
	caller authz.Caller,
	organizationID uuid.UUID,
	q ListQuery,
) ([]ClinicListItem, error) {
	if err := r.authz.Require(
		ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(organizationID),
	); err != nil {
		return nil, err
	}
	if err := q.normalize(); err != nil {
		return nil, err
	}
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT id, organization_id, name
		  FROM projections.clinics
		 WHERE organization_id = ?
		 ORDER BY created_at DESC, id DESC
		 LIMIT ? OFFSET ?`, organizationID, q.Limit, q.Offset,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.orgstructure.clinic").
			Code(ErrCodeClinicLoadFailed).
			With("organization_id", organizationID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]ClinicListItem, 0, q.Limit)
	for rows.Next() {
		var v ClinicListItem
		if err := rows.Scan(&v.ID, &v.OrganizationID, &v.Name); err != nil {
			return nil, oops.In("reader.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("organization_id", organizationID).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.orgstructure.clinic").
			Code(ErrCodeClinicLoadFailed).
			With("organization_id", organizationID).
			Wrap(err)
	}
	return out, nil
}
