// Package classifier is the query-side reader for the request type classifier.
//
// See: docs/services/request/Classifier.md
package classifier

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

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
) ([]RequestTypeView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}
	if err := q.normalize(); err != nil {
		return nil, err
	}
	rows, err := r.db.WithContext(ctx).Raw(selectRequestType+`
		 WHERE organization_id = ?
		 ORDER BY created_at DESC, id DESC
		 LIMIT ? OFFSET ?`, orgID, q.Limit, q.Offset,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]RequestTypeView, 0, q.Limit)
	for rows.Next() {
		var v RequestTypeView
		if err := scanRequestType(rows, &v); err != nil {
			return nil, oops.In("reader.request.classifier").
				Code(ErrCodeRequestTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			Wrap(err)
	}
	return out, nil
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
) ([]RequestTypeView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}
	if err := q.normalize(); err != nil {
		return nil, err
	}
	rows, err := r.db.WithContext(ctx).Raw(selectRequestType+`
		 WHERE organization_id = ? AND is_active = TRUE
		 ORDER BY name ASC, id ASC
		 LIMIT ? OFFSET ?`, orgID, q.Limit, q.Offset,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]RequestTypeView, 0, q.Limit)
	for rows.Next() {
		var v RequestTypeView
		if err := scanRequestType(rows, &v); err != nil {
			return nil, oops.In("reader.request.classifier").
				Code(ErrCodeRequestTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.request.classifier").
			Code(ErrCodeRequestTypeLoadFailed).
			Wrap(err)
	}
	return out, nil
}
