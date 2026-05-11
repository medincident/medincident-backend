// Package buffer is the query-side reader for projections.patient_incident_buffer.
//
// See: docs/services/incident/Buffer.md
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

	"github.com/medincident/medincident-backend/internal/cursor"
	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	queryincident "github.com/medincident/medincident-backend/internal/service/query/incident"
)

const (
	ErrCodeBufferReadFailed = "buffer_query_read_failed"
	ErrCodeBufferNotFound   = "buffer_query_not_found"
	ErrCodeListBadCursor    = "incident_bad_cursor"
)

// BufferListResult is returned by paginated buffer list methods.
type BufferListResult struct {
	Items      []BufferEntryView
	NextCursor *string
}

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
//
// See: docs/services/incident/Buffer.md
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
	After    *string
}

// ListBufferEntries returns buffer entries in an org. Visible only to
// roles that may see the buffer for that org (SystemAdmin, OrgAdmin,
// OrgDispatcher).
//
// See: docs/services/incident/Buffer.md
func (r *Reader) ListBufferEntries(
	ctx context.Context, callerID string, orgID uuid.UUID, f *ListBufferFilters,
) (BufferListResult, error) {
	cc, err := r.incidentRdr.ResolveCaller(ctx, callerID)
	if err != nil {
		return BufferListResult{}, err
	}
	if !cc.CanSeeBufferForOrg(orgID) {
		return BufferListResult{}, oops.In(scope).Code(authz.ErrCodePermissionDenied).
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
	limit := normLimit(0)
	if f != nil {
		limit = normLimit(f.Limit)
		if f.After != nil {
			t, idStr, err := cursor.Decode(*f.After)
			if err != nil {
				return BufferListResult{}, oops.In(scope).
					Code(ErrCodeListBadCursor).
					Public("Invalid pagination cursor.").
					Wrap(err)
			}
			bufID, _ := uuid.Parse(idStr)
			conds = append(conds, "(updated_at, id) < (?, ?)")
			args = append(args, t, bufID)
		}
	}
	args = append(args, limit+1)
	q := `SELECT ` + bufferSelect + ` FROM projections.patient_incident_buffer WHERE ` +
		strings.Join(conds, " AND ") + ` ORDER BY updated_at DESC, id DESC LIMIT ?`

	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return BufferListResult{}, wrapRead(err, "list buffer")
	}
	defer func() { _ = rows.Close() }()
	out := make([]BufferEntryView, 0, limit+1)
	for rows.Next() {
		var v BufferEntryView
		if err := scanBuffer(rows, &v); err != nil {
			return BufferListResult{}, wrapRead(err, "scan buffer row")
		}
		v.PatientPerspective = false
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return BufferListResult{}, wrapRead(err, "iterate buffer rows")
	}
	var nextCursor *string
	if len(out) > limit {
		out = out[:limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return BufferListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListMyBufferEntries returns buffer entries the caller submitted as
// a patient. For staff callers it can also surface their own
// (the filter is by patient_zitadel_user_id, regardless of role).
//
// See: docs/services/incident/Buffer.md
func (r *Reader) ListMyBufferEntries(
	ctx context.Context, callerID string, limit int, after *string,
) (BufferListResult, error) {
	cc, err := r.incidentRdr.ResolveCaller(ctx, callerID)
	if err != nil {
		return BufferListResult{}, err
	}
	limit = normLimit(limit)
	args := []any{callerID}
	cursorCond := ""
	if after != nil {
		t, idStr, err := cursor.Decode(*after)
		if err != nil {
			return BufferListResult{}, oops.In(scope).
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		bufID, _ := uuid.Parse(idStr)
		cursorCond = ` AND (updated_at, id) < (?, ?)`
		args = append(args, t, bufID)
	}
	args = append(args, limit+1)
	q := `SELECT ` + bufferSelect + ` FROM projections.patient_incident_buffer
		WHERE patient_zitadel_user_id = ?` + cursorCond +
		` ORDER BY updated_at DESC, id DESC LIMIT ?`
	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return BufferListResult{}, wrapRead(err, "list my buffer")
	}
	defer func() { _ = rows.Close() }()
	out := make([]BufferEntryView, 0, limit+1)
	for rows.Next() {
		var v BufferEntryView
		if err := scanBuffer(rows, &v); err != nil {
			return BufferListResult{}, wrapRead(err, "scan my buffer row")
		}
		v.PatientPerspective = cc.IsPatient()
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return BufferListResult{}, wrapRead(err, "iterate my buffer rows")
	}
	var nextCursor *string
	if len(out) > limit {
		out = out[:limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return BufferListResult{Items: out, NextCursor: nextCursor}, nil
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

// normLimit delegates to the incident reader's exported helper
// so the two readers don't drift.
func normLimit(limit int) int { return queryincident.NormLimit(limit) }

func wrapRead(err error, action string) error {
	return oops.In(scope).Code(ErrCodeBufferReadFailed).
		With("action", action).Wrap(err)
}
