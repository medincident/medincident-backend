package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	organizationapp "github.com/medincident/medincident-command-service/internal/orgstructure/organization/app"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
	sqlcgen "github.com/medincident/medincident-command-service/internal/storage/postgres/sqlc/gen"
	"github.com/medincident/medincident-command-service/internal/tx"
)

// Error codes emitted by the Organization repository in this file.
const (
	CodeOrganizationNotFound        = "postgres_organization_not_found"
	CodeOrganizationGetFailed       = "postgres_organization_get_failed"
	CodeOrganizationSaveFailed      = "postgres_organization_save_failed"
	CodeOrganizationListFailed      = "postgres_organization_list_failed"
	CodeOrganizationListBuildFailed = "postgres_organization_list_build_failed"
)

// OrganizationRepo implements organizationapp.Repository on Postgres.
type OrganizationRepo struct {
	pool *Pool
}

// NewOrganizationRepo returns a repository bound to the given pool.
func NewOrganizationRepo(pool *Pool) *OrganizationRepo {
	return &OrganizationRepo{pool: pool}
}

// executor returns the underlying pgx handle — either the tx in ctx
// (for writes that must share a transaction with the outbox) or the
// pool for plain reads.
func (r *OrganizationRepo) executor(ctx context.Context) (sqlcgen.DBTX, error) {
	if t, ok := tx.FromContext(ctx); ok {
		pg, err := unwrap(t)
		if err != nil {
			return nil, err
		}
		return pg.raw, nil
	}
	return r.pool.Pool, nil
}

// GetByID loads an Organization by primary key via the sqlc-generated
// query, then calls organization.Hydrate to build the aggregate
// (no validation, no events).
func (r *OrganizationRepo) GetByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error) {
	exec, err := r.executor(ctx)
	if err != nil {
		return nil, err
	}
	q := sqlcgen.New(exec)
	row, err := q.GetOrganization(ctx, uuidToPg(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, oops.In("storage.postgres").
				Code(CodeOrganizationNotFound).
				Public("Organization not found.").
				With("id", id).
				Errorf("no rows")
		}
		return nil, oops.In("storage.postgres").
			Code(CodeOrganizationGetFailed).
			With("id", id).
			Wrap(err)
	}

	return organization.Hydrate(
		pgToUUID(row.ID),
		row.Name,
		derefStr(row.Description),
		addressFromRow(row.LegalAddressText, row.LegalAddressLng, row.LegalAddressLat),
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}

// Save upserts the aggregate's current state via the sqlc-generated
// UpsertOrganization query.
func (r *OrganizationRepo) Save(ctx context.Context, o *organization.Organization) error {
	exec, err := r.executor(ctx)
	if err != nil {
		return err
	}
	q := sqlcgen.New(exec)

	addrText, addrLng, addrLat := addressToRow(o.LegalAddress)
	descPtr := strPtrIfNotEmpty(o.Description)

	if err := q.UpsertOrganization(ctx, sqlcgen.UpsertOrganizationParams{
		ID:               uuidToPg(o.ID),
		Name:             o.Name,
		Description:      descPtr,
		LegalAddressText: addrText,
		LegalAddressLng:  addrLng,
		LegalAddressLat:  addrLat,
		CreatedAt:        pgtype.Timestamptz{Time: o.CreatedAt, Valid: true},
		UpdatedAt:        pgtype.Timestamptz{Time: o.UpdatedAt, Valid: true},
	}); err != nil {
		return oops.In("storage.postgres").
			Code(CodeOrganizationSaveFailed).
			With("id", o.ID).
			Wrap(err)
	}
	return nil
}

// List returns organizations matching the filter. Uses squirrel to
// build the optional WHERE / LIMIT / OFFSET fragments; sqlc is not a
// fit here because the query shape varies with the filter.
func (r *OrganizationRepo) List(ctx context.Context, f organizationapp.ListFilter) ([]*organization.Organization, error) {
	exec, err := r.executor(ctx)
	if err != nil {
		return nil, err
	}

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"name",
			"description",
			"legal_address_text",
			"legal_address_lng",
			"legal_address_lat",
			"created_at",
			"updated_at",
		).
		From("domain.organizations").
		OrderBy("created_at DESC")

	if f.NameLike != "" {
		qb = qb.Where(sq.ILike{"name": "%" + f.NameLike + "%"})
	}
	if f.Limit > 0 {
		qb = qb.Limit(uint64(f.Limit))
	}
	if f.Offset > 0 {
		qb = qb.Offset(uint64(f.Offset))
	}

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, oops.In("storage.postgres").
			Code(CodeOrganizationListBuildFailed).
			Wrap(err)
	}

	rows, err := exec.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, oops.In("storage.postgres").
			Code(CodeOrganizationListFailed).
			Wrap(err)
	}
	defer rows.Close()

	var out []*organization.Organization
	for rows.Next() {
		var (
			id          pgtype.UUID
			name        string
			descPtr     *string
			addrTextPtr *string
			addrLngPtr  *float64
			addrLatPtr  *float64
			createdAt   pgtype.Timestamptz
			updatedAt   pgtype.Timestamptz
		)
		if err := rows.Scan(&id, &name, &descPtr, &addrTextPtr, &addrLngPtr, &addrLatPtr, &createdAt, &updatedAt); err != nil {
			return nil, oops.In("storage.postgres").
				Code(CodeOrganizationListFailed).
				Wrap(err)
		}
		out = append(out, organization.Hydrate(
			pgToUUID(id),
			name,
			derefStr(descPtr),
			addressFromRow(addrTextPtr, addrLngPtr, addrLatPtr),
			createdAt.Time,
			updatedAt.Time,
		))
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("storage.postgres").
			Code(CodeOrganizationListFailed).
			Wrap(err)
	}
	return out, nil
}

var _ organizationapp.Repository = (*OrganizationRepo)(nil)

// --- scalar <-> VO conversion helpers (flat columns, no composite types) ---

// addressFromRow builds a *geo.Address from the three flat columns.
// Returns nil when text is nil (no address stored at all). If text is
// present but the coordinates are not, returns an Address with Point=nil.
func addressFromRow(text *string, lng, lat *float64) *geo.Address {
	if text == nil {
		return nil
	}
	a := &geo.Address{Text: *text}
	if lng != nil && lat != nil {
		a.Point = &geo.Point{Longitude: *lng, Latitude: *lat}
	}
	return a
}

// addressToRow decomposes a *geo.Address into the three flat column
// arguments. nil address → all three nil pointers. Address without
// Point → text set, lng/lat nil.
func addressToRow(a *geo.Address) (text *string, lng, lat *float64) {
	if a == nil {
		return nil, nil, nil
	}
	t := a.Text
	text = &t
	if a.Point != nil {
		lngV := a.Point.Longitude
		latV := a.Point.Latitude
		lng = &lngV
		lat = &latV
	}
	return text, lng, lat
}

func strPtrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// uuidToPg converts a uuid.UUID into the pgtype.UUID form sqlc expects.
func uuidToPg(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// pgToUUID converts a pgtype.UUID back into a uuid.UUID. A NULL pg value
// returns uuid.Nil — caller is responsible for null handling, but in
// practice every id column we read is NOT NULL.
func pgToUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	var out uuid.UUID
	copy(out[:], p.Bytes[:])
	return out
}
