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

// Error codes emitted by DepartmentReader methods.
const (
	ErrCodeDepartmentNotFound    = "department_not_found"
	ErrCodeDepartmentLoadFailed  = "department_load_failed"
	ErrCodeDepartmentCountFailed = "department_count_failed"
)

// DepartmentDetails is the full card of a department returned by Get.
type DepartmentDetails struct {
	ID          uuid.UUID
	ClinicID    uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DepartmentListItem is the minimal view returned by paginated list
// endpoints.
type DepartmentListItem struct {
	ID        uuid.UUID
	ClinicID  uuid.UUID
	Name      string
	UpdatedAt time.Time
}

// DepartmentListResult is returned by ListByClinic. NextCursor is nil
// when no more pages remain.
type DepartmentListResult struct {
	Items      []DepartmentListItem
	NextCursor *string
}

// Get returns the DepartmentDetails for the given id. Authorization:
// authz.ReaderOf.Department(id) — system admin, organization admin
// of the owning org, or any employee of the owning org.
//
// See: docs/services/OrgStructure.md
func (r *DepartmentReader) Get(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*DepartmentDetails, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Department(id)); err != nil {
		return nil, err
	}
	var out DepartmentDetails
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, clinic_id, name, description, created_at, updated_at
		  FROM projections.departments
		 WHERE id = ?`, id,
	).Row().Scan(
		&out.ID, &out.ClinicID, &out.Name, &out.Description, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.orgstructure.department").
				Code(ErrCodeDepartmentNotFound).
				Public("Department not found.").
				With("department_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.orgstructure.department").
			Code(ErrCodeDepartmentLoadFailed).
			With("department_id", id).
			Wrap(err)
	}
	return &out, nil
}

// CountByClinic returns the total number of departments belonging to
// the given clinic. Authorization: authz.ReaderOf.Clinic(clinicID).
//
// See: docs/services/OrgStructure.md
func (r *DepartmentReader) CountByClinic(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
) (int64, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID)); err != nil {
		return 0, err
	}
	var total int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.departments WHERE clinic_id = ?`,
		clinicID,
	).Row().Scan(&total); err != nil {
		return 0, oops.In("reader.orgstructure.department").
			Code(ErrCodeDepartmentCountFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	return total, nil
}

// ListByClinic returns up to q.Limit departments belonging to the
// given clinic, ordered most-recently-updated first. Authorization:
// authz.ReaderOf.Clinic(clinicID). Pagination bounds are normalized
// first so a malformed Limit cannot trigger a gratuitous authz
// DB round-trip — matching the validate→authorize order used on the
// command side.
//
// See: docs/services/OrgStructure.md
func (r *DepartmentReader) ListByClinic(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
	q ListQuery,
) (DepartmentListResult, error) {
	if err := q.normalize(); err != nil {
		return DepartmentListResult{}, err
	}
	if err := r.authz.Require(
		ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID),
	); err != nil {
		return DepartmentListResult{}, err
	}
	sqlBuf := `SELECT id, clinic_id, name, updated_at
		  FROM projections.departments
		 WHERE clinic_id = ?`
	args := make([]any, 0, 4)
	args = append(args, clinicID)
	if q.After != nil {
		t, idStr, err := cursor.Decode(*q.After)
		if err != nil {
			return DepartmentListResult{}, oops.In("reader.orgstructure.department").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		id, _ := uuid.Parse(idStr)
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, t, id)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return DepartmentListResult{}, oops.In("reader.orgstructure.department").
			Code(ErrCodeDepartmentLoadFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]DepartmentListItem, 0, q.Limit+1)
	for rows.Next() {
		var v DepartmentListItem
		if err := rows.Scan(&v.ID, &v.ClinicID, &v.Name, &v.UpdatedAt); err != nil {
			return DepartmentListResult{}, oops.In("reader.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return DepartmentListResult{}, oops.In("reader.orgstructure.department").
			Code(ErrCodeDepartmentLoadFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}

	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return DepartmentListResult{Items: out, NextCursor: nextCursor}, nil
}
