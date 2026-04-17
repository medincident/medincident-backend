package orgstructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
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
	ID       uuid.UUID
	ClinicID uuid.UUID
	Name     string
}

// Get returns the DepartmentDetails for the given id.
func (r *DepartmentReader) Get(ctx context.Context, id uuid.UUID) (*DepartmentDetails, error) {
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

// ListByClinic returns up to q.Limit departments belonging to the
// given clinic, ordered most-recently-created first.
func (r *DepartmentReader) ListByClinic(
	ctx context.Context,
	clinicID uuid.UUID,
	q ListQuery,
) ([]DepartmentListItem, error) {
	if err := q.normalize(); err != nil {
		return nil, err
	}
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT id, clinic_id, name
		  FROM projections.departments
		 WHERE clinic_id = ?
		 ORDER BY created_at DESC, id DESC
		 LIMIT ? OFFSET ?`, clinicID, q.Limit, q.Offset,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.orgstructure.department").
			Code(ErrCodeDepartmentLoadFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]DepartmentListItem, 0, q.Limit)
	for rows.Next() {
		var v DepartmentListItem
		if err := rows.Scan(&v.ID, &v.ClinicID, &v.Name); err != nil {
			return nil, oops.In("reader.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("clinic_id", clinicID).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.orgstructure.department").
			Code(ErrCodeDepartmentLoadFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	return out, nil
}
