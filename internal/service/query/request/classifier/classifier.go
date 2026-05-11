package classifier

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

// Error codes emitted by the request type reader.
const (
	ErrCodeRequestTypeNotFound   = "request_type_not_found"
	ErrCodeRequestTypeLoadFailed = "request_type_load_failed"
)

// RequestTypeView mirrors projections.request_types.
type RequestTypeView struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    *string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

const selectRequestType = `
	SELECT id, organization_id, name, description,
	       is_active, created_at, updated_at
	  FROM projections.request_types`

// scanRequestType scans one row into a RequestTypeView.
func scanRequestType(scanner interface {
	Scan(dest ...any) error
}, out *RequestTypeView,
) error {
	return scanner.Scan(
		&out.ID, &out.OrganizationID, &out.Name, &out.Description,
		&out.IsActive, &out.CreatedAt, &out.UpdatedAt,
	)
}

// GetRequestType returns one request type by id. Authorization:
// authz.ReaderOf.RequestType(id).
//
// See: docs/services/request/Classifier.md
func (r *Reader) GetRequestType(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*RequestTypeView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.RequestType(id)); err != nil {
		return nil, err
	}
	var out RequestTypeView
	err := scanRequestType(r.db.WithContext(ctx).Raw(selectRequestType+` WHERE id = ?`, id).Row(), &out)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.request.classifier").
				Code(ErrCodeRequestTypeNotFound).
				Public("Request type not found.").
				With("request_type_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			With("request_type_id", id).
			Wrap(err)
	}
	return &out, nil
}

// ListRequestTypesByOrganization paginates request types for one org.
// Authorization: authz.ReaderOf.Organization(orgID).
//
// See: docs/services/request/Classifier.md
func (r *Reader) ListRequestTypesByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (RequestTypeListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return RequestTypeListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return RequestTypeListResult{}, err
	}
	sqlBuf := selectRequestType + ` WHERE organization_id = ?`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return RequestTypeListResult{}, oops.In("reader.request.classifier").
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
		return RequestTypeListResult{}, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]RequestTypeView, 0, q.Limit+1)
	for rows.Next() {
		var v RequestTypeView
		if err := scanRequestType(rows, &v); err != nil {
			return RequestTypeListResult{}, oops.In("reader.request.classifier").
				Code(ErrCodeRequestTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return RequestTypeListResult{}, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return RequestTypeListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListActiveRequestTypesByOrganization returns every active request
// type for one org. Authorization: authz.ReaderOf.Organization(orgID).
//
// See: docs/services/request/Classifier.md
func (r *Reader) ListActiveRequestTypesByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (RequestTypeListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return RequestTypeListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return RequestTypeListResult{}, err
	}
	sqlBuf := selectRequestType + ` WHERE organization_id = ? AND is_active = TRUE`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return RequestTypeListResult{}, oops.In("reader.request.classifier").
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
		return RequestTypeListResult{}, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]RequestTypeView, 0, q.Limit+1)
	for rows.Next() {
		var v RequestTypeView
		if err := scanRequestType(rows, &v); err != nil {
			return RequestTypeListResult{}, oops.In("reader.request.classifier").
				Code(ErrCodeRequestTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return RequestTypeListResult{}, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return RequestTypeListResult{Items: out, NextCursor: nextCursor}, nil
}
