package membership

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Error codes emitted by EmployeeReader.
const (
	ErrCodeEmployeeNotFound   = "employee_card_not_found"
	ErrCodeEmployeeLoadFailed = "employee_card_load_failed"
	ErrCodeVacationLoadFailed = "vacation_load_failed"
)

// EmployeeCardView mirrors projections.employee_cards. Optional columns
// are represented as *T so the downstream proto layer can omit them
// when NULL.
type EmployeeCardView struct {
	EmployeeID            uuid.UUID
	ZitadelUserID         string
	FirstName             *string
	LastName              *string
	DisplayName           *string
	Email                 *string
	OrganizationID        uuid.UUID
	OrganizationName      *string
	ClinicID              *uuid.UUID
	ClinicName            *string
	DepartmentID          uuid.UUID
	DepartmentName        *string
	Position              *string
	TerminatedAt          *time.Time
	CurrentVacationEndsAt *time.Time
	NextVacationStartsAt  *time.Time
}

// selectEmployeeCard is the reusable SELECT list matching the scan
// order used by every employee_card query.
const selectEmployeeCard = `
	SELECT employee_id, zitadel_user_id,
	       first_name, last_name, display_name, email,
	       organization_id, organization_name,
	       clinic_id, clinic_name,
	       department_id, department_name,
	       position, terminated_at,
	       current_vacation_ends_at, next_vacation_starts_at
	  FROM projections.employee_cards`

// scanEmployeeCard reads one row from a *sql.Rows cursor into a view.
func scanEmployeeCard(scanner interface {
	Scan(dest ...any) error
}, out *EmployeeCardView,
) error {
	return scanner.Scan(
		&out.EmployeeID, &out.ZitadelUserID,
		&out.FirstName, &out.LastName, &out.DisplayName, &out.Email,
		&out.OrganizationID, &out.OrganizationName,
		&out.ClinicID, &out.ClinicName,
		&out.DepartmentID, &out.DepartmentName,
		&out.Position, &out.TerminatedAt,
		&out.CurrentVacationEndsAt, &out.NextVacationStartsAt,
	)
}

// Get returns the employee_card row for the given id. Authorization:
// authz.ReaderOf.Employee(id) — system admin, organization admin of
// the employee's org, or any employee of the same organization.
func (r *EmployeeReader) Get(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Employee(id)); err != nil {
		return nil, err
	}
	var out EmployeeCardView
	err := scanEmployeeCard(
		r.db.WithContext(ctx).Raw(selectEmployeeCard+` WHERE employee_id = ?`, id).Row(),
		&out,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.membership.employee").
				Code(ErrCodeEmployeeNotFound).
				Public("Employee not found.").
				With("employee_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeLoadFailed).
			With("employee_id", id).
			Wrap(err)
	}
	return &out, nil
}

// listByField is shared by the three ListEmployeesByX methods.
func (r *EmployeeReader) listByField(
	ctx context.Context,
	field string,
	value uuid.UUID,
	q ListQuery,
) ([]EmployeeCardView, error) {
	if err := q.normalize(); err != nil {
		return nil, err
	}
	rows, err := r.db.WithContext(ctx).Raw(selectEmployeeCard+
		` WHERE `+field+` = ?
		 ORDER BY updated_at DESC, employee_id DESC
		 LIMIT ? OFFSET ?`, value, q.Limit, q.Offset,
	).Rows()
	if err != nil {
		return nil, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeLoadFailed).
			With(field, value).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]EmployeeCardView, 0, q.Limit)
	for rows.Next() {
		var v EmployeeCardView
		if err := scanEmployeeCard(rows, &v); err != nil {
			return nil, oops.In("reader.membership.employee").
				Code(ErrCodeEmployeeLoadFailed).
				With(field, value).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeLoadFailed).
			With(field, value).
			Wrap(err)
	}
	return out, nil
}

// ListByDepartment returns cards under a department. Authorization:
// authz.ReaderOf.Department(deptID).
func (r *EmployeeReader) ListByDepartment(
	ctx context.Context,
	caller authz.Caller,
	deptID uuid.UUID,
	q ListQuery,
) ([]EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Department(deptID)); err != nil {
		return nil, err
	}
	return r.listByField(ctx, "department_id", deptID, q)
}

// ListByClinic returns cards under a clinic. Authorization:
// authz.ReaderOf.Clinic(clinicID).
func (r *EmployeeReader) ListByClinic(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
	q ListQuery,
) ([]EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID)); err != nil {
		return nil, err
	}
	return r.listByField(ctx, "clinic_id", clinicID, q)
}

// ListByOrganization returns cards under an organization. Authorization:
// authz.ReaderOf.Organization(orgID).
func (r *EmployeeReader) ListByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) ([]EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}
	return r.listByField(ctx, "organization_id", orgID, q)
}

// VacationView mirrors projections.employee_vacations.
type VacationView struct {
	ID         uuid.UUID
	EmployeeID uuid.UUID
	State      string
	StartsAt   time.Time
	EndsAt     *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ListVacationsByEmployee returns vacation rows for an employee. If
// state is non-empty, it is used as an exact-match filter. Authorization
// is tighter than the employee card: only system admins, admins of the
// owning organization, and the employee themselves may read vacations —
// no sibling-employee access. The policy is composed inline because
// the "self" branch is specific to this read and does not reuse the
// ReaderOf battery.
func (r *EmployeeReader) ListVacationsByEmployee(
	ctx context.Context,
	caller authz.Caller,
	employeeID uuid.UUID,
	state string,
) ([]VacationView, error) {
	policy := authz.AnyOf(
		authz.SystemAdmin,
		authz.OrgAdminOf.Employee(employeeID),
		authz.SelfEmployee(employeeID),
	)
	if err := r.authz.Require(ctx, caller.ZitadelUserID, policy); err != nil {
		return nil, err
	}
	var (
		rows *sql.Rows
		err  error
	)
	if state == "" {
		rows, err = r.db.WithContext(ctx).Raw(`
			SELECT id, employee_id, state, starts_at, ends_at, created_at, updated_at
			  FROM projections.employee_vacations
			 WHERE employee_id = ?
			 ORDER BY starts_at DESC, id DESC`, employeeID,
		).Rows()
	} else {
		rows, err = r.db.WithContext(ctx).Raw(`
			SELECT id, employee_id, state, starts_at, ends_at, created_at, updated_at
			  FROM projections.employee_vacations
			 WHERE employee_id = ? AND state = ?
			 ORDER BY starts_at DESC, id DESC`, employeeID, state,
		).Rows()
	}
	if err != nil {
		return nil, oops.In("reader.membership.vacation").
			Code(ErrCodeVacationLoadFailed).
			With("employee_id", employeeID).
			With("state", state).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]VacationView, 0)
	for rows.Next() {
		var v VacationView
		if err := rows.Scan(&v.ID, &v.EmployeeID, &v.State, &v.StartsAt, &v.EndsAt, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, oops.In("reader.membership.vacation").
				Code(ErrCodeVacationLoadFailed).
				With("employee_id", employeeID).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.membership.vacation").
			Code(ErrCodeVacationLoadFailed).
			With("employee_id", employeeID).
			Wrap(err)
	}
	return out, nil
}
