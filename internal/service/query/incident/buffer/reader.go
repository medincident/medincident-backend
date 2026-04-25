// Package buffer is the query-side reader for projections.patient_incident_buffer.
package buffer

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	queryincident "github.com/medincident/medincident-backend/internal/service/query/incident"
)

const (
	ErrCodeBufferReadFailed = "buffer_query_read_failed"
	ErrCodeBufferNotFound   = "buffer_query_not_found"
)

const scope = "services.query.incident.buffer"

type Reader struct {
	db          *gorm.DB
	logger      *zerolog.Logger
	incidentRdr *queryincident.Reader
}

func NewReader(db *gorm.DB, logger *zerolog.Logger, incRdr *queryincident.Reader) *Reader {
	return &Reader{db: db, logger: logger, incidentRdr: incRdr}
}

// BufferEntryView is a projection row.
type BufferEntryView struct {
	ID                   uuid.UUID
	OrganizationID       uuid.UUID
	PatientZitadelUserID string
	CategoryID           uuid.NullUUID
	TypeID               uuid.NullUUID
	Description          null.String
	OccurredAt           null.Time
	Status               model.BufferStatus
	PublishedIncidentID  uuid.NullUUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	// True when the caller is treated as a patient — handlers redact
	// fields that should not surface to non-staff.
	PatientPerspective bool
}

const bufferSelect = `id, organization_id, patient_zitadel_user_id,
	category_id, type_id, description, occurred_at, status, published_incident_id,
	created_at, updated_at`

func scanBuffer(row interface{ Scan(...any) error }, v *BufferEntryView) error {
	return row.Scan(
		&v.ID, &v.OrganizationID, &v.PatientZitadelUserID,
		&v.CategoryID, &v.TypeID, &v.Description, &v.OccurredAt,
		&v.Status, &v.PublishedIncidentID, &v.CreatedAt, &v.UpdatedAt,
	)
}

// GetBufferEntry returns one buffer row if the caller may see it.
func (r *Reader) GetBufferEntry(
	ctx context.Context, callerID string, id uuid.UUID,
) (*BufferEntryView, error) {
	cc, err := r.incidentRdr.ResolveCaller(ctx, callerID)
	if err != nil {
		return nil, err
	}
	q := `SELECT ` + bufferSelect + ` FROM projections.patient_incident_buffer WHERE id = ? LIMIT 1`
	row := r.db.WithContext(ctx).Raw(q, id).Row()
	var v BufferEntryView
	if err := scanBuffer(row, &v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In(scope).Code(ErrCodeBufferNotFound).
				Public("Patient incident not found.").With("buffer_id", id).Wrap(err)
		}
		return nil, wrapRead(err, "scan buffer")
	}
	if !canSeeBuffer(cc, &v) {
		return nil, oops.In(scope).Code(ErrCodeBufferNotFound).
			Public("Patient incident not found.").With("buffer_id", id).
			Errorf("not visible")
	}
	v.PatientPerspective = cc.IsPatient()
	return &v, nil
}

// ListBufferFilters captures the filters for listing.
type ListBufferFilters struct {
	Statuses []model.BufferStatus
	Limit    int
	Offset   int
}

// ListBufferEntries returns buffer entries in an org. Visible only to
// roles that may see the buffer for that org (SystemAdmin, OrgAdmin,
// OrgDispatcher).
func (r *Reader) ListBufferEntries(
	ctx context.Context, callerID string, orgID uuid.UUID, f *ListBufferFilters,
) ([]BufferEntryView, error) {
	cc, err := r.incidentRdr.ResolveCaller(ctx, callerID)
	if err != nil {
		return nil, err
	}
	if !cc.CanSeeBufferForOrg(orgID) {
		return nil, oops.In(scope).Code("permission_denied").
			Public("Not authorized to view this organization's patient buffer.").
			With("organization_id", orgID).Errorf("denied")
	}
	conds := []string{"organization_id = ?"}
	args := []any{orgID}
	if f != nil && len(f.Statuses) > 0 {
		ph := strings.TrimSuffix(strings.Repeat("?,", len(f.Statuses)), ",")
		conds = append(conds, "status IN ("+ph+")")
		for _, s := range f.Statuses {
			args = append(args, s)
		}
	}
	limit, offset := paginationDefaults(0, 0)
	if f != nil {
		limit, offset = paginationDefaults(f.Limit, f.Offset)
	}
	args = append(args, limit, offset)
	q := `SELECT ` + bufferSelect + ` FROM projections.patient_incident_buffer WHERE ` +
		strings.Join(conds, " AND ") + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return nil, wrapRead(err, "list buffer")
	}
	defer func() { _ = rows.Close() }()
	out := []BufferEntryView{}
	for rows.Next() {
		var v BufferEntryView
		if err := scanBuffer(rows, &v); err != nil {
			return nil, wrapRead(err, "scan buffer row")
		}
		v.PatientPerspective = false
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapRead(err, "iterate buffer rows")
	}
	return out, nil
}

// ListMyBufferEntries returns buffer entries the caller submitted as
// a patient. For staff callers it can also surface their own
// (the filter is by patient_zitadel_user_id, regardless of role).
func (r *Reader) ListMyBufferEntries(
	ctx context.Context, callerID string, limit, offset int,
) ([]BufferEntryView, error) {
	cc, err := r.incidentRdr.ResolveCaller(ctx, callerID)
	if err != nil {
		return nil, err
	}
	limit, offset = paginationDefaults(limit, offset)
	q := `SELECT ` + bufferSelect + ` FROM projections.patient_incident_buffer
		WHERE patient_zitadel_user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.WithContext(ctx).Raw(q, callerID, limit, offset).Rows()
	if err != nil {
		return nil, wrapRead(err, "list my buffer")
	}
	defer func() { _ = rows.Close() }()
	out := []BufferEntryView{}
	for rows.Next() {
		var v BufferEntryView
		if err := scanBuffer(rows, &v); err != nil {
			return nil, wrapRead(err, "scan my buffer row")
		}
		v.PatientPerspective = cc.IsPatient()
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapRead(err, "iterate my buffer rows")
	}
	return out, nil
}

// canSeeBuffer encapsulates the per-row visibility test for GetBufferEntry.
func canSeeBuffer(cc *queryincident.CallerContext, b *BufferEntryView) bool {
	if cc.IsSystemAdmin() {
		return true
	}
	if cc.IsOrgAdminOf(b.OrganizationID) || cc.IsOrgDispatcherOf(b.OrganizationID) {
		return true
	}
	return b.PatientZitadelUserID == cc.ZitadelID()
}

func paginationDefaults(limit, offset int) (outLimit, outOffset int) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func wrapRead(err error, action string) error {
	return oops.In(scope).Code(ErrCodeBufferReadFailed).
		With("action", action).Wrap(err)
}
