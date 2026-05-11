package membership

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

// Error codes emitted by RoleReader.
const (
	ErrCodeRoleLoadFailed = "role_load_failed"
)

// ErrRoleVacant is returned by single-holder role lookups when the
// parent aggregate has no assignment. Call sites use errors.Is to
// distinguish this from transport / DB failures.
var ErrRoleVacant = errors.New("role vacant")

// RoleAssignmentView is a role row enriched with the denormalised
// employee_cards row for the holder and, when present, the deputy.
// The reader joins projections.employee_cards in the same round-trip
// as the role table so downstream callers can render a role without
// an N+1 GetEmployee follow-up.
type RoleAssignmentView struct {
	UpdatedAt time.Time
	Holder    EmployeeCardView
	Deputy    *EmployeeCardView
}

// RoleListResult is returned by all paginated role list methods.
// NextCursor is nil when no more pages remain.
type RoleListResult struct {
	Items      []RoleAssignmentView
	NextCursor *string
}

// SystemAdminView mirrors projections.system_admins.
type SystemAdminView struct {
	ZitadelUserID string
	CreatedAt     time.Time
}

// SystemAdminListResult is returned by ListSystemAdmins. NextCursor is
// nil when no more pages remain.
type SystemAdminListResult struct {
	Items      []SystemAdminView
	NextCursor *string
}

// employeeCardColumns returns the ordered SELECT list for one
// employee_cards row, with every column qualified by the given table
// alias. The column order matches ScanEmployeeCard / scanEmployeeCardPtrs.
func employeeCardColumns(alias string) string {
	return alias + `.employee_id, ` + alias + `.zitadel_user_id, ` +
		alias + `.first_name, ` + alias + `.last_name, ` +
		alias + `.display_name, ` + alias + `.email, ` +
		alias + `.organization_id, ` + alias + `.organization_name, ` +
		alias + `.clinic_id, ` + alias + `.clinic_name, ` +
		alias + `.department_id, ` + alias + `.department_name, ` +
		alias + `.position, ` + alias + `.terminated_at, ` +
		alias + `.current_vacation_ends_at, ` + alias + `.next_vacation_starts_at, ` +
		alias + `.updated_at`
}

// deputyCardScan carries every column of the deputy's employee_cards
// row through pointers so a LEFT JOIN that missed the deputy side
// scans cleanly into nils rather than tripping a NULL-to-UUID error.
// Non-nullable columns (employee_id, zitadel_user_id, organization_id,
// department_id) become nil for a vacant deputy; when all are nil the
// scanner returns a nil *EmployeeCardView.
type deputyCardScan struct {
	EmployeeID            *uuid.UUID
	ZitadelUserID         *string
	FirstName             *string
	LastName              *string
	DisplayName           *string
	Email                 *string
	OrganizationID        *uuid.UUID
	OrganizationName      *string
	ClinicID              *uuid.UUID
	ClinicName            *string
	DepartmentID          *uuid.UUID
	DepartmentName        *string
	Position              *string
	TerminatedAt          *time.Time
	CurrentVacationEndsAt *time.Time
	NextVacationStartsAt  *time.Time
	UpdatedAt             *time.Time
}

// scanTargets returns the pointer slice to pass into rows.Scan in the
// same order as employeeCardColumns.
func (d *deputyCardScan) scanTargets() []any {
	return []any{
		&d.EmployeeID, &d.ZitadelUserID,
		&d.FirstName, &d.LastName, &d.DisplayName, &d.Email,
		&d.OrganizationID, &d.OrganizationName,
		&d.ClinicID, &d.ClinicName,
		&d.DepartmentID, &d.DepartmentName,
		&d.Position, &d.TerminatedAt,
		&d.CurrentVacationEndsAt, &d.NextVacationStartsAt,
		&d.UpdatedAt,
	}
}

// toView materialises a full EmployeeCardView when the deputy was
// present, or nil when the LEFT JOIN produced NULLs across the board.
// Presence is keyed off the primary key (employee_id) — if it is set
// the other NOT-NULL columns are guaranteed to be set too.
func (d *deputyCardScan) toView() *EmployeeCardView {
	if d.EmployeeID == nil {
		return nil
	}
	out := &EmployeeCardView{
		EmployeeID:            *d.EmployeeID,
		FirstName:             d.FirstName,
		LastName:              d.LastName,
		DisplayName:           d.DisplayName,
		Email:                 d.Email,
		OrganizationName:      d.OrganizationName,
		ClinicID:              d.ClinicID,
		ClinicName:            d.ClinicName,
		DepartmentName:        d.DepartmentName,
		Position:              d.Position,
		TerminatedAt:          d.TerminatedAt,
		CurrentVacationEndsAt: d.CurrentVacationEndsAt,
		NextVacationStartsAt:  d.NextVacationStartsAt,
	}
	if d.ZitadelUserID != nil {
		out.ZitadelUserID = *d.ZitadelUserID
	}
	if d.OrganizationID != nil {
		out.OrganizationID = *d.OrganizationID
	}
	if d.DepartmentID != nil {
		out.DepartmentID = *d.DepartmentID
	}
	if d.UpdatedAt != nil {
		out.UpdatedAt = *d.UpdatedAt
	}
	return out
}

// scanAssignment reads one joined row into UpdatedAt + Holder + (optional) Deputy.
// role.updated_at is scanned first, then holder columns (INNER JOIN guarantees
// presence), then deputy columns via deputyCardScan so NULL-on-LEFT-JOIN stays
// nil rather than erroring.
func scanAssignment(scanner interface {
	Scan(dest ...any) error
}, out *RoleAssignmentView,
) error {
	var dep deputyCardScan
	dests := make([]any, 0, 35)
	dests = append(dests,
		&out.UpdatedAt,
		&out.Holder.EmployeeID, &out.Holder.ZitadelUserID,
		&out.Holder.FirstName, &out.Holder.LastName, &out.Holder.DisplayName, &out.Holder.Email,
		&out.Holder.OrganizationID, &out.Holder.OrganizationName,
		&out.Holder.ClinicID, &out.Holder.ClinicName,
		&out.Holder.DepartmentID, &out.Holder.DepartmentName,
		&out.Holder.Position, &out.Holder.TerminatedAt,
		&out.Holder.CurrentVacationEndsAt, &out.Holder.NextVacationStartsAt,
		&out.Holder.UpdatedAt,
	)
	dests = append(dests, dep.scanTargets()...)
	if err := scanner.Scan(dests...); err != nil {
		return err
	}
	out.Deputy = dep.toView()
	return nil
}

// selectRoleAssignment composes the SELECT list: role.updated_at first,
// then every holder card column followed by every deputy card column,
// in the same order scanAssignment expects.
var selectRoleAssignment = `
	SELECT role.updated_at, ` + employeeCardColumns("holder") + `,
	       ` + employeeCardColumns("deputy")

// GetClinicHead returns the clinic-head assignment for the clinic, or
// ErrRoleVacant if none exists. Authorization: authz.ReaderOf.Clinic.
//
// See: docs/services/Membership.md
func (r *RoleReader) GetClinicHead(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
) (*RoleAssignmentView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID)); err != nil {
		return nil, err
	}
	return r.oneRoleByParent(ctx, "projections.clinic_heads", "clinic_id", clinicID)
}

// GetDepartmentResponsible returns the responsible assignment for the
// department, or ErrRoleVacant if none exists. Authorization:
// authz.ReaderOf.Department.
//
// See: docs/services/Membership.md
func (r *RoleReader) GetDepartmentResponsible(
	ctx context.Context,
	caller authz.Caller,
	departmentID uuid.UUID,
) (*RoleAssignmentView, error) {
	if err := r.authz.Require(
		ctx, caller.ZitadelUserID, authz.ReaderOf.Department(departmentID),
	); err != nil {
		return nil, err
	}
	return r.oneRoleByParent(ctx, "projections.department_responsibles", "department_id", departmentID)
}

// ListOrgAdmins returns all org-admin assignments for the organization.
// Authorization: authz.ReaderOf.Organization.
//
// See: docs/services/Membership.md
func (r *RoleReader) ListOrgAdmins(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (RoleListResult, error) {
	if err := q.normalize(); err != nil {
		return RoleListResult{}, err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return RoleListResult{}, err
	}
	return r.listRolesByParent(ctx, "projections.org_admins", "organization_id", orgID, q)
}

// ListOrgDispatchers returns all org-dispatcher assignments for the
// organization. Authorization: authz.ReaderOf.Organization.
//
// See: docs/services/Membership.md
func (r *RoleReader) ListOrgDispatchers(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (RoleListResult, error) {
	if err := q.normalize(); err != nil {
		return RoleListResult{}, err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return RoleListResult{}, err
	}
	return r.listRolesByParent(ctx, "projections.org_dispatchers", "organization_id", orgID, q)
}

// ListOrgHeads returns all org-head assignments for the organization.
// Authorization: authz.ReaderOf.Organization.
//
// See: docs/services/Membership.md
func (r *RoleReader) ListOrgHeads(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (RoleListResult, error) {
	if err := q.normalize(); err != nil {
		return RoleListResult{}, err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return RoleListResult{}, err
	}
	return r.listRolesByParent(ctx, "projections.org_heads", "organization_id", orgID, q)
}

// ListSystemAdmins returns every system-admin row. Authorization:
// authz.SystemAdmin — only system admins can enumerate their peers.
//
// See: docs/services/Membership.md
func (r *RoleReader) ListSystemAdmins(
	ctx context.Context,
	caller authz.Caller,
	q ListQuery,
) (SystemAdminListResult, error) {
	if err := q.normalize(); err != nil {
		return SystemAdminListResult{}, err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return SystemAdminListResult{}, err
	}
	sqlBuf := `SELECT zitadel_user_id, created_at
		  FROM projections.system_admins
		 WHERE 1=1`
	args := make([]any, 0, 3)
	if q.After != nil {
		t, zitadelID, err := cursor.Decode(*q.After)
		if err != nil {
			return SystemAdminListResult{}, oops.In("reader.membership.role").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (created_at, zitadel_user_id) < (?, ?)`
		args = append(args, t, zitadelID)
	}
	sqlBuf += ` ORDER BY created_at DESC, zitadel_user_id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return SystemAdminListResult{}, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("role", "system_admin").
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]SystemAdminView, 0, q.Limit+1)
	for rows.Next() {
		var v SystemAdminView
		if err := rows.Scan(&v.ZitadelUserID, &v.CreatedAt); err != nil {
			return SystemAdminListResult{}, oops.In("reader.membership.role").
				Code(ErrCodeRoleLoadFailed).
				With("role", "system_admin").
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return SystemAdminListResult{}, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("role", "system_admin").
			Wrap(err)
	}

	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.CreatedAt, last.ZitadelUserID)
		nextCursor = &s
	}
	return SystemAdminListResult{Items: out, NextCursor: nextCursor}, nil
}

// oneRoleByParent returns the first role row (by employee_id asc) for
// a given parent aggregate, or ErrRoleVacant when the role is vacant.
// Clinic head and department responsible are schema-level
// one-employee-per-parent but the PK permits multiple rows; taking the
// first by employee_id keeps the read idempotent. The SELECT joins
// projections.employee_cards (INNER for holder, LEFT for deputy) so
// the full card view comes back in one round-trip.
func (r *RoleReader) oneRoleByParent(ctx context.Context, table, parentField string, parentID uuid.UUID) (*RoleAssignmentView, error) {
	var v RoleAssignmentView
	err := scanAssignment(
		r.db.WithContext(ctx).Raw(
			selectRoleAssignment+`
			   FROM `+table+` role
			   JOIN projections.employee_cards holder
			     ON holder.employee_id = role.employee_id
			   LEFT JOIN projections.employee_cards deputy
			     ON deputy.employee_id = role.deputy_employee_id
			  WHERE role.`+parentField+` = ?
			  ORDER BY role.employee_id ASC
			  LIMIT 1`, parentID,
		).Row(),
		&v,
	)
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

// listRolesByParent returns a paginated slice of role rows for a given
// parent aggregate (used by org-level roles which legitimately can have
// multiple holders). Callers must have already called q.normalize().
// The SELECT joins projections.employee_cards for holder (INNER) and
// deputy (LEFT) so each returned view carries the full card.
// Ordered by role.updated_at DESC, role.employee_id DESC for stable keyset pagination.
func (r *RoleReader) listRolesByParent(ctx context.Context, table, parentField string, parentID uuid.UUID, q ListQuery) (RoleListResult, error) {
	sqlBuf := selectRoleAssignment + `
		   FROM ` + table + ` role
		   JOIN projections.employee_cards holder
		     ON holder.employee_id = role.employee_id
		   LEFT JOIN projections.employee_cards deputy
		     ON deputy.employee_id = role.deputy_employee_id
		  WHERE role.` + parentField + ` = ?`
	args := make([]any, 0, 4)
	args = append(args, parentID)
	if q.After != nil {
		t, idStr, err := cursor.Decode(*q.After)
		if err != nil {
			return RoleListResult{}, oops.In("reader.membership.role").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		empID, _ := uuid.Parse(idStr)
		sqlBuf += ` AND (role.updated_at, role.employee_id) < (?, ?)`
		args = append(args, t, empID)
	}
	sqlBuf += ` ORDER BY role.updated_at DESC, role.employee_id DESC LIMIT ?`
	args = append(args, q.Limit+1)

	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return RoleListResult{}, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("table", table).
			With(parentField, parentID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]RoleAssignmentView, 0, q.Limit+1)
	for rows.Next() {
		var v RoleAssignmentView
		if err := scanAssignment(rows, &v); err != nil {
			return RoleListResult{}, oops.In("reader.membership.role").
				Code(ErrCodeRoleLoadFailed).
				With("table", table).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return RoleListResult{}, oops.In("reader.membership.role").
			Code(ErrCodeRoleLoadFailed).
			With("table", table).
			Wrap(err)
	}

	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.Holder.EmployeeID.String())
		nextCursor = &s
	}
	return RoleListResult{Items: out, NextCursor: nextCursor}, nil
}
