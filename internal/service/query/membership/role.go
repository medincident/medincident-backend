package membership

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by RoleReader.
const (
	ErrCodeRoleLoadFailed = "role_load_failed"
)

// ErrRoleVacant is returned by single-holder role lookups when the
// parent aggregate has no assignment. Call sites use errors.Is to
// distinguish this from transport / DB failures.
var ErrRoleVacant = errors.New("role vacant")

// RoleHolderView represents a single role assignment (employee_id
// plus optional deputy).
type RoleHolderView struct {
	EmployeeID       uuid.UUID
	DeputyEmployeeID *uuid.UUID
}

// GetClinicHead returns the clinic-head assignment for the clinic, or
// ErrRoleVacant if none exists.
func (r *RoleReader) GetClinicHead(ctx context.Context, clinicID uuid.UUID) (*RoleHolderView, error) {
	return r.oneRoleByParent(ctx, "projections.clinic_heads", "clinic_id", clinicID)
}

// GetDepartmentResponsible returns the responsible assignment for the
// department, or ErrRoleVacant if none exists.
func (r *RoleReader) GetDepartmentResponsible(ctx context.Context, departmentID uuid.UUID) (*RoleHolderView, error) {
	return r.oneRoleByParent(ctx, "projections.department_responsibles", "department_id", departmentID)
}

// ListOrgAdmins returns all org-admin holders for the organization.
func (r *RoleReader) ListOrgAdmins(ctx context.Context, orgID uuid.UUID) ([]RoleHolderView, error) {
	return r.listRolesByParent(ctx, "projections.org_admins", "organization_id", orgID)
}

// ListOrgDispatchers returns all org-dispatcher holders for the organization.
func (r *RoleReader) ListOrgDispatchers(ctx context.Context, orgID uuid.UUID) ([]RoleHolderView, error) {
	return r.listRolesByParent(ctx, "projections.org_dispatchers", "organization_id", orgID)
}

// ListOrgHeads returns all org-head holders for the organization.
func (r *RoleReader) ListOrgHeads(ctx context.Context, orgID uuid.UUID) ([]RoleHolderView, error) {
	return r.listRolesByParent(ctx, "projections.org_heads", "organization_id", orgID)
}

// SystemAdminView mirrors projections.system_admins.
type SystemAdminView struct {
	ZitadelUserID string
	CreatedAt     time.Time
}

// ListSystemAdmins returns every system-admin row.
func (r *RoleReader) ListSystemAdmins(ctx context.Context) ([]SystemAdminView, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT zitadel_user_id, created_at
		  FROM projections.system_admins
		 ORDER BY created_at DESC, zitadel_user_id DESC`).Rows()
	if err != nil {
		return nil, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("role", "system_admin").
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]SystemAdminView, 0)
	for rows.Next() {
		var v SystemAdminView
		if err := rows.Scan(&v.ZitadelUserID, &v.CreatedAt); err != nil {
			return nil, oops.In("reader.membership.role").
				Code(ErrCodeRoleLoadFailed).
				With("role", "system_admin").
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("role", "system_admin").
			Wrap(err)
	}
	return out, nil
}

// oneRoleByParent returns the first role row (by employee_id asc) for
// a given parent aggregate, or nil when the role is vacant. Clinic
// head and department responsible are schema-level one-employee-per-
// parent but the PK permits multiple rows; taking the first by
// employee_id keeps the read idempotent.
func (r *RoleReader) oneRoleByParent(ctx context.Context, table, parentField string, parentID uuid.UUID) (*RoleHolderView, error) {
	var v RoleHolderView
	err := r.db.WithContext(ctx).Raw(
		`SELECT employee_id, deputy_employee_id
		   FROM `+table+`
		  WHERE `+parentField+` = ?
		  ORDER BY employee_id ASC
		  LIMIT 1`, parentID,
	).Row().Scan(&v.EmployeeID, &v.DeputyEmployeeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleVacant
		}
		return nil, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("table", table).
			With(parentField, parentID).
			Wrap(err)
	}
	return &v, nil
}

// listRolesByParent returns every role row for a given parent aggregate
// (used by org-level roles which legitimately can have multiple holders).
func (r *RoleReader) listRolesByParent(ctx context.Context, table, parentField string, parentID uuid.UUID) ([]RoleHolderView, error) {
	rows, err := r.db.WithContext(ctx).Raw(
		`SELECT employee_id, deputy_employee_id
		   FROM `+table+`
		  WHERE `+parentField+` = ?
		  ORDER BY employee_id ASC`, parentID,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("table", table).
			With(parentField, parentID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]RoleHolderView, 0)
	for rows.Next() {
		var v RoleHolderView
		if err := rows.Scan(&v.EmployeeID, &v.DeputyEmployeeID); err != nil {
			return nil, oops.In("reader.membership.role").
				Code(ErrCodeRoleLoadFailed).
				With("table", table).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("table", table).
			Wrap(err)
	}
	return out, nil
}
