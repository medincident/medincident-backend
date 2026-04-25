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
	ErrCodeEmployeeNotFound           = "employee_card_not_found"
	ErrCodeEmployeeLoadFailed         = "employee_card_load_failed"
	ErrCodeEmployeeCountFailed        = "employee_card_count_failed"
	ErrCodeEmployeeSearchQueryTooLong = "employee_search_query_too_long"
	ErrCodeVacationLoadFailed         = "vacation_load_failed"
	ErrCodeVacationCountFailed        = "vacation_count_failed"
)

// employeeSearchMaxQueryLength caps the user-supplied query for
// SearchByOrganization. Over-long inputs are rejected before any
// authz round-trip.
const employeeSearchMaxQueryLength = 256

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

// EmployeeFilter is the optional filter shared by every List / Count
// over employee_cards. Zero value = no restriction beyond the scope.
// IncludeTerminated=false (default) hides rows whose terminated_at is
// set; OnVacation=true restricts to employees currently on an active
// vacation; Position, when non-empty, is an exact match on the
// employee_cards.position column (the handler layer is responsible for
// trimming whitespace so the zero value here is "no filter").
type EmployeeFilter struct {
	IncludeTerminated bool
	OnVacation        bool
	Position          string
}

// buildFilterClause returns a SQL fragment starting with " AND ..."
// (or the empty string when no filters apply) together with the
// positional args to bind. Position is always bound as a parameter —
// never concatenated — so an attacker cannot inject SQL through it.
func (f EmployeeFilter) buildFilterClause() (clause string, args []any) {
	args = make([]any, 0, 3)
	if !f.IncludeTerminated {
		clause += ` AND terminated_at IS NULL`
	}
	if f.OnVacation {
		clause += ` AND current_vacation_ends_at IS NOT NULL AND current_vacation_ends_at > now()`
	}
	if f.Position != "" {
		clause += ` AND position = ?`
		args = append(args, f.Position)
	}
	return clause, args
}

// listByField is shared by the three ListEmployeesByX methods.
func (r *EmployeeReader) listByField(
	ctx context.Context,
	field string,
	value uuid.UUID,
	q ListQuery,
	filter EmployeeFilter,
) ([]EmployeeCardView, error) {
	if err := q.normalize(); err != nil {
		return nil, err
	}
	filterClause, filterArgs := filter.buildFilterClause()
	args := make([]any, 0, 3+len(filterArgs))
	args = append(args, value)
	args = append(args, filterArgs...)
	args = append(args, q.Limit, q.Offset)
	rows, err := r.db.WithContext(ctx).Raw(selectEmployeeCard+
		` WHERE `+field+` = ?`+filterClause+
		` ORDER BY updated_at DESC, employee_id DESC
		 LIMIT ? OFFSET ?`, args...,
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
	filter EmployeeFilter,
) ([]EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Department(deptID)); err != nil {
		return nil, err
	}
	return r.listByField(ctx, "department_id", deptID, q, filter)
}

// ListByClinic returns cards under a clinic. Authorization:
// authz.ReaderOf.Clinic(clinicID).
func (r *EmployeeReader) ListByClinic(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
	q ListQuery,
	filter EmployeeFilter,
) ([]EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID)); err != nil {
		return nil, err
	}
	return r.listByField(ctx, "clinic_id", clinicID, q, filter)
}

// ListByOrganization returns cards under an organization. Authorization:
// authz.ReaderOf.Organization(orgID).
func (r *EmployeeReader) ListByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
	filter EmployeeFilter,
) ([]EmployeeCardView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}
	return r.listByField(ctx, "organization_id", orgID, q, filter)
}

// SearchByOrganization returns cards under an organization whose
// first_name, last_name, display_name, or email contain the given
// query substring (case-insensitive, ILIKE %query%). An empty query
// behaves like ListByOrganization. The query is length-capped and
// validated BEFORE the authz round-trip; authorization matches
// ListByOrganization so that cross-org callers see the same
// unauthorized behaviour for List and Search. Query is bound
// positionally — never concatenated — so users cannot inject SQL.
func (r *EmployeeReader) SearchByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	query string,
	q ListQuery,
	filter EmployeeFilter,
) ([]EmployeeCardView, error) {
	if len(query) > employeeSearchMaxQueryLength {
		return nil, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeSearchQueryTooLong).
			Public("Search query is too long.").
			With("max_length", employeeSearchMaxQueryLength).
			With("actual_length", len(query)).
			Errorf("search query too long")
	}
	if err := q.normalize(); err != nil {
		return nil, err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}
	filterClause, filterArgs := filter.buildFilterClause()
	sqlBuf := selectEmployeeCard + ` WHERE organization_id = ?` + filterClause
	args := make([]any, 0, 3+len(filterArgs)+4)
	args = append(args, orgID)
	args = append(args, filterArgs...)
	if query != "" {
		sqlBuf += ` AND (
			COALESCE(first_name, '')   ILIKE ? OR
			COALESCE(last_name, '')    ILIKE ? OR
			COALESCE(display_name, '') ILIKE ? OR
			COALESCE(email, '')        ILIKE ?
		)`
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}
	sqlBuf += ` ORDER BY updated_at DESC, employee_id DESC
		LIMIT ? OFFSET ?`
	args = append(args, q.Limit, q.Offset)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return nil, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]EmployeeCardView, 0, q.Limit)
	for rows.Next() {
		var v EmployeeCardView
		if err := scanEmployeeCard(rows, &v); err != nil {
			return nil, oops.In("reader.membership.employee").
				Code(ErrCodeEmployeeLoadFailed).
				With("organization_id", orgID).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	return out, nil
}

// countByField is shared by the three CountEmployeesByX methods.
func (r *EmployeeReader) countByField(
	ctx context.Context,
	field string,
	value uuid.UUID,
	filter EmployeeFilter,
) (int64, error) {
	filterClause, filterArgs := filter.buildFilterClause()
	args := make([]any, 0, 1+len(filterArgs))
	args = append(args, value)
	args = append(args, filterArgs...)
	var total int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT count(*) FROM projections.employee_cards WHERE `+field+` = ?`+filterClause,
		args...,
	).Row().Scan(&total); err != nil {
		return 0, oops.In("reader.membership.employee").
			Code(ErrCodeEmployeeCountFailed).
			With(field, value).
			Wrap(err)
	}
	return total, nil
}

// CountByDepartment returns the total employees under a department.
// Authorization: authz.ReaderOf.Department(deptID).
func (r *EmployeeReader) CountByDepartment(
	ctx context.Context,
	caller authz.Caller,
	deptID uuid.UUID,
	filter EmployeeFilter,
) (int64, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Department(deptID)); err != nil {
		return 0, err
	}
	return r.countByField(ctx, "department_id", deptID, filter)
}

// CountByClinic returns the total employees under a clinic.
// Authorization: authz.ReaderOf.Clinic(clinicID).
func (r *EmployeeReader) CountByClinic(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
	filter EmployeeFilter,
) (int64, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID)); err != nil {
		return 0, err
	}
	return r.countByField(ctx, "clinic_id", clinicID, filter)
}

// CountByOrganization returns the total employees under an organization.
// Authorization: authz.ReaderOf.Organization(orgID).
func (r *EmployeeReader) CountByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	filter EmployeeFilter,
) (int64, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return 0, err
	}
	return r.countByField(ctx, "organization_id", orgID, filter)
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

// vacationAuthzPolicy composes the "system admin OR org admin of the
// employee OR the employee themselves" gate shared by read and count
// over employee vacations.
func vacationAuthzPolicy(employeeID uuid.UUID) authz.Policy {
	return authz.AnyOf(
		authz.SystemAdmin,
		authz.OrgAdminOf.Employee(employeeID),
		authz.SelfEmployee(employeeID),
	)
}

// CountVacationsByEmployee returns the total vacation rows for an
// employee. If state is non-empty, it is used as an exact-match
// filter. Authorization matches ListVacationsByEmployee — only system
// admins, admins of the owning organization, and the employee
// themselves may count.
func (r *EmployeeReader) CountVacationsByEmployee(
	ctx context.Context,
	caller authz.Caller,
	employeeID uuid.UUID,
	state string,
) (int64, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, vacationAuthzPolicy(employeeID)); err != nil {
		return 0, err
	}
	var (
		total int64
		row   *sql.Row
	)
	if state == "" {
		row = r.db.WithContext(ctx).Raw(
			`SELECT count(*) FROM projections.employee_vacations WHERE employee_id = ?`,
			employeeID,
		).Row()
	} else {
		row = r.db.WithContext(ctx).Raw(
			`SELECT count(*) FROM projections.employee_vacations
			  WHERE employee_id = ? AND state = ?`,
			employeeID, state,
		).Row()
	}
	if err := row.Scan(&total); err != nil {
		return 0, oops.In("reader.membership.vacation").
			Code(ErrCodeVacationCountFailed).
			With("employee_id", employeeID).
			With("state", state).
			Wrap(err)
	}
	return total, nil
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
	q ListQuery,
) ([]VacationView, error) {
	if err := q.normalize(); err != nil {
		return nil, err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, vacationAuthzPolicy(employeeID)); err != nil {
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
			 ORDER BY starts_at DESC, id DESC
			 LIMIT ? OFFSET ?`, employeeID, q.Limit, q.Offset,
		).Rows()
	} else {
		rows, err = r.db.WithContext(ctx).Raw(`
			SELECT id, employee_id, state, starts_at, ends_at, created_at, updated_at
			  FROM projections.employee_vacations
			 WHERE employee_id = ? AND state = ?
			 ORDER BY starts_at DESC, id DESC
			 LIMIT ? OFFSET ?`, employeeID, state, q.Limit, q.Offset,
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

	out := make([]VacationView, 0, q.Limit)
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
