package membership

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/cursor"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query"
	"github.com/medincident/medincident-backend/internal/util/like"
)

// Error codes emitted by CandidateReader.
const (
	ErrCodeCandidateSearchQueryTooLong = "candidate_search_query_too_long"
	ErrCodeCandidateLoadFailed         = "candidate_load_failed"
)

const candidateSearchMaxQueryLength = 256

// ZitadelUserView is a lightweight projection of projections.users used
// by ForHire and ForSystemAdmin.
type ZitadelUserView struct {
	ZitadelUserID string
	FirstName     string
	LastName      string
	DisplayName   string
	Email         string
	UpdatedAt     time.Time
}

// CandidateReader queries eligible users/employees for hire and role assignment.
// All methods follow: validate → authz → SQL.
//
// See: docs/services/Membership.md
type CandidateReader struct {
	db    *gorm.DB
	authz *authz.Authz
	log   *zerolog.Logger
}

// NewCandidateReader wires the reader with DB, authz, and logger.
func NewCandidateReader(db *gorm.DB, az *authz.Authz, log *zerolog.Logger) *CandidateReader {
	return &CandidateReader{db: db, authz: az, log: log}
}

func validateCandidateQuery(q string) error {
	if len(q) > candidateSearchMaxQueryLength {
		return oops.In("reader.membership.candidate").
			Code(ErrCodeCandidateSearchQueryTooLong).
			Public("Search query is too long.").
			With("max_length", candidateSearchMaxQueryLength).
			With("actual_length", len(q)).
			Errorf("search query too long")
	}
	return nil
}

func normalizeCandidateLimit(limit int) (int, error) {
	if limit == 0 {
		return query.DefaultLimit, nil
	}
	if limit < query.MinLimit || limit > query.MaxLimit {
		return 0, oops.In("reader.membership.candidate").
			Code(ErrCodeListLimitOutOfRange).
			Public("List limit is out of range.").
			With("actual_value", limit).
			With("min_value", query.MinLimit).
			With("max_value", query.MaxLimit).
			Errorf("limit out of range")
	}
	return limit, nil
}

const selectUser = `
	SELECT u.id, u.first_name, u.last_name, u.display_name, u.email, u.updated_at
	  FROM projections.users u`

// selectEmployeeCardEC selects all employee_card columns plus updated_at,
// using the "ec." alias prefix, for cursor generation.
const selectEmployeeCardEC = `
	SELECT ec.employee_id, ec.zitadel_user_id,
	       ec.first_name, ec.last_name, ec.display_name, ec.email,
	       ec.organization_id, ec.organization_name,
	       ec.clinic_id, ec.clinic_name,
	       ec.department_id, ec.department_name,
	       ec.position, ec.terminated_at,
	       ec.current_vacation_ends_at, ec.next_vacation_starts_at,
	       ec.updated_at
	  FROM projections.employee_cards ec`

func scanUserRow(rows interface{ Scan(dest ...any) error }, out *ZitadelUserView) error {
	return rows.Scan(
		&out.ZitadelUserID, &out.FirstName, &out.LastName,
		&out.DisplayName, &out.Email, &out.UpdatedAt,
	)
}

// employeeWithUpdatedAt wraps EmployeeCardView to capture updated_at for cursor generation.
type employeeWithUpdatedAt struct {
	EmployeeCardView
	UpdatedAt time.Time
}

func scanEmployeeWithUpdatedAt(rows interface{ Scan(dest ...any) error }, out *employeeWithUpdatedAt) error {
	return rows.Scan(
		&out.EmployeeID, &out.ZitadelUserID,
		&out.FirstName, &out.LastName, &out.DisplayName, &out.Email,
		&out.OrganizationID, &out.OrganizationName,
		&out.ClinicID, &out.ClinicName,
		&out.DepartmentID, &out.DepartmentName,
		&out.Position, &out.TerminatedAt,
		&out.CurrentVacationEndsAt, &out.NextVacationStartsAt,
		&out.UpdatedAt,
	)
}

// appendUserSearchAndCursor appends ILIKE search and keyset cursor conditions for
// projections.users queries (alias "u"). hasCursor reports whether cur is valid.
// Returns the extra SQL fragment and args.
func appendUserSearchAndCursor(searchQuery string, cur cursor.Cursor, hasCursor bool) (clause string, args []any) {
	if searchQuery != "" {
		clause += ` AND (
			u.first_name   ILIKE ? OR
			u.last_name    ILIKE ? OR
			u.display_name ILIKE ? OR
			u.email        ILIKE ?
		)`
		p := "%" + like.EscapePattern(searchQuery) + "%"
		args = append(args, p, p, p, p)
	}
	if hasCursor {
		clause += ` AND (u.updated_at < ? OR (u.updated_at = ? AND u.id < ?))`
		t := cur.Time()
		args = append(args, t, t, cur.I)
	}
	return clause, args
}

func appendEmployeeSearchAndCursor(searchQuery string, cur cursor.Cursor, hasCursor bool) (clause string, args []any) {
	if searchQuery != "" {
		clause += ` AND (
			COALESCE(ec.first_name, '')    ILIKE ? OR
			COALESCE(ec.last_name, '')     ILIKE ? OR
			COALESCE(ec.display_name, '')  ILIKE ? OR
			COALESCE(ec.email, '')         ILIKE ?
		)`
		p := "%" + like.EscapePattern(searchQuery) + "%"
		args = append(args, p, p, p, p)
	}
	if hasCursor {
		clause += ` AND (ec.updated_at < ? OR (ec.updated_at = ? AND ec.employee_id < ?))`
		t := cur.Time()
		args = append(args, t, t, cur.I)
	}
	return clause, args
}

// decodeCursor parses after into a Cursor. The bool return value reports
// whether a cursor was present: false means "no cursor" (empty after),
// true means a cursor was decoded successfully. An error is returned only
// when after is non-empty but malformed.
func decodeCursor(after string) (cursor.Cursor, bool, error) {
	if after == "" {
		return cursor.Cursor{}, false, nil
	}
	c, err := cursor.Decode(after)
	if err != nil {
		return cursor.Cursor{}, false, err
	}
	return c, true, nil
}

// ForHire returns Zitadel users not yet active employees of the given org.
// Authorization: AdminOf.Organization(orgID).
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForHire(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	searchQuery, after string,
	limit int,
) ([]ZitadelUserView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	cur, hasCursor, err := decodeCursor(after)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return nil, "", err
	}
	extra, extraArgs := appendUserSearchAndCursor(searchQuery, cur, hasCursor)
	sql := selectUser + `
	WHERE NOT EXISTS (
		SELECT 1 FROM projections.employee_cards ec
		WHERE ec.zitadel_user_id = u.id
		  AND ec.organization_id = ?
		  AND ec.terminated_at IS NULL
	)` + extra + `
	ORDER BY u.updated_at DESC, u.id DESC
	LIMIT ?`
	args := make([]any, 0, 2+len(extraArgs)+1)
	args = append(args, orgID)
	args = append(args, extraArgs...)
	args = append(args, limit)
	return r.queryUsers(ctx, sql, args, limit)
}

// ForSystemAdmin returns Zitadel users not yet system admins.
// Authorization: SystemAdmin.
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForSystemAdmin(
	ctx context.Context,
	caller authz.Caller,
	searchQuery, after string,
	limit int,
) ([]ZitadelUserView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	cur, hasCursor, err := decodeCursor(after)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return nil, "", err
	}
	extra, extraArgs := appendUserSearchAndCursor(searchQuery, cur, hasCursor)
	sql := selectUser + `
	WHERE NOT EXISTS (
		SELECT 1 FROM projections.system_admins sa
		WHERE sa.zitadel_user_id = u.id
	)` + extra + `
	ORDER BY u.updated_at DESC, u.id DESC
	LIMIT ?`
	args := make([]any, 0, len(extraArgs)+1)
	args = append(args, extraArgs...)
	args = append(args, limit)
	return r.queryUsers(ctx, sql, args, limit)
}

func (r *CandidateReader) queryUsers(ctx context.Context, sql string, args []any, limit int) ([]ZitadelUserView, string, error) {
	// Fetch limit+1 to distinguish "has next page" from "exact multiple of limit".
	args[len(args)-1] = limit + 1
	rows, err := r.db.WithContext(ctx).Raw(sql, args...).Rows()
	if err != nil {
		return nil, "", oops.In("reader.membership.candidate").
			Code(ErrCodeCandidateLoadFailed).Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]ZitadelUserView, 0, limit+1)
	for rows.Next() {
		var v ZitadelUserView
		if err := scanUserRow(rows, &v); err != nil {
			return nil, "", oops.In("reader.membership.candidate").
				Code(ErrCodeCandidateLoadFailed).Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, "", oops.In("reader.membership.candidate").
			Code(ErrCodeCandidateLoadFailed).Wrap(err)
	}
	next := ""
	if len(out) > limit {
		out = out[:limit]
		last := out[len(out)-1]
		next = cursor.Encode(last.UpdatedAt, last.ZitadelUserID)
	}
	return out, next, nil
}

// NOT EXISTS SQL fragments — compile-time constants, never user input.
const notExistsOrgAdmin = `NOT EXISTS (
	SELECT 1 FROM projections.org_admins x
	WHERE x.employee_id = ec.employee_id
	  AND x.organization_id = ec.organization_id
)`

const notExistsOrgHead = `NOT EXISTS (
	SELECT 1 FROM projections.org_heads x
	WHERE x.employee_id = ec.employee_id
	  AND x.organization_id = ec.organization_id
)`

const notExistsOrgDispatcher = `NOT EXISTS (
	SELECT 1 FROM projections.org_dispatchers x
	WHERE x.employee_id = ec.employee_id
	  AND x.organization_id = ec.organization_id
)`

const notExistsClinicHead = `NOT EXISTS (
	SELECT 1 FROM projections.clinic_heads x
	WHERE x.employee_id = ec.employee_id
	  AND x.clinic_id = ec.clinic_id
)`

const notExistsDeptResponsible = `NOT EXISTS (
	SELECT 1 FROM projections.department_responsibles x
	WHERE x.employee_id = ec.employee_id
	  AND x.department_id = ec.department_id
)`

// listEmployeeCandidates executes the shared NOT EXISTS + scope filter query.
// Receives already-validated searchQuery and limit.
func (r *CandidateReader) listEmployeeCandidates(
	ctx context.Context,
	scopeField string, // "organization_id", "clinic_id", or "department_id"
	scopeID uuid.UUID,
	notExistsSQL string, // compile-time constant, never user input
	searchQuery, after string,
	limit int,
) ([]EmployeeCardView, string, error) {
	cur, hasCursor, err := decodeCursor(after)
	if err != nil {
		return nil, "", err
	}
	// employee_id is a UUID column; a non-UUID cursor.I would cause a Postgres
	// cast error reported as 500 instead of 400. Validate early.
	if hasCursor {
		if _, err := uuid.Parse(cur.I); err != nil {
			return nil, "", oops.In("reader.membership.candidate").
				Code(cursor.ErrCodeInvalidCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
	}
	extra, extraArgs := appendEmployeeSearchAndCursor(searchQuery, cur, hasCursor)
	sql := selectEmployeeCardEC + `
	WHERE ec.` + scopeField + ` = ?
	  AND ec.terminated_at IS NULL
	  AND ` + notExistsSQL + extra + `
	ORDER BY ec.updated_at DESC, ec.employee_id DESC
	LIMIT ?`
	args := make([]any, 0, 2+len(extraArgs)+1)
	args = append(args, scopeID)
	args = append(args, extraArgs...)
	// Fetch limit+1 to distinguish "has next page" from "exact multiple of limit".
	args = append(args, limit+1)

	rows, err := r.db.WithContext(ctx).Raw(sql, args...).Rows()
	if err != nil {
		return nil, "", oops.In("reader.membership.candidate").
			Code(ErrCodeCandidateLoadFailed).Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	raw := make([]employeeWithUpdatedAt, 0, limit+1)
	for rows.Next() {
		var v employeeWithUpdatedAt
		if err := scanEmployeeWithUpdatedAt(rows, &v); err != nil {
			return nil, "", oops.In("reader.membership.candidate").
				Code(ErrCodeCandidateLoadFailed).Wrap(err)
		}
		raw = append(raw, v)
	}
	if err := rows.Err(); err != nil {
		return nil, "", oops.In("reader.membership.candidate").
			Code(ErrCodeCandidateLoadFailed).Wrap(err)
	}
	next := ""
	if len(raw) > limit {
		raw = raw[:limit]
		last := raw[len(raw)-1]
		next = cursor.Encode(last.UpdatedAt, last.EmployeeID.String())
	}
	out := make([]EmployeeCardView, len(raw))
	for i := range raw {
		out[i] = raw[i].EmployeeCardView
	}
	return out, next, nil
}

// ForOrgAdmin returns active org employees not yet assigned as org admin.
// Authorization: AdminOf.Organization(orgID).
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForOrgAdmin(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	searchQuery, after string,
	limit int,
) ([]EmployeeCardView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return nil, "", err
	}
	return r.listEmployeeCandidates(ctx, "organization_id", orgID, notExistsOrgAdmin, searchQuery, after, limit)
}

// ForOrgHead returns active org employees not yet assigned as org head.
// Authorization: AdminOf.Organization(orgID).
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForOrgHead(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	searchQuery, after string,
	limit int,
) ([]EmployeeCardView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return nil, "", err
	}
	return r.listEmployeeCandidates(ctx, "organization_id", orgID, notExistsOrgHead, searchQuery, after, limit)
}

// ForOrgDispatcher returns active org employees not yet assigned as dispatcher.
// Authorization: AdminOf.Organization(orgID).
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForOrgDispatcher(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	searchQuery, after string,
	limit int,
) ([]EmployeeCardView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return nil, "", err
	}
	return r.listEmployeeCandidates(ctx, "organization_id", orgID, notExistsOrgDispatcher, searchQuery, after, limit)
}

// ForClinicHead returns active clinic employees not yet assigned as clinic head.
// Authorization: AdminOf.Clinic(clinicID).
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForClinicHead(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
	searchQuery, after string,
	limit int,
) ([]EmployeeCardView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.AdminOf.Clinic(clinicID)); err != nil {
		return nil, "", err
	}
	return r.listEmployeeCandidates(ctx, "clinic_id", clinicID, notExistsClinicHead, searchQuery, after, limit)
}

// ForDeptResponsible returns active dept employees not yet assigned as responsible.
// Authorization: AdminOf.Department(deptID).
//
// See: docs/services/Membership.md
func (r *CandidateReader) ForDeptResponsible(
	ctx context.Context,
	caller authz.Caller,
	deptID uuid.UUID,
	searchQuery, after string,
	limit int,
) ([]EmployeeCardView, string, error) {
	if err := validateCandidateQuery(searchQuery); err != nil {
		return nil, "", err
	}
	limit, err := normalizeCandidateLimit(limit)
	if err != nil {
		return nil, "", err
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.AdminOf.Department(deptID)); err != nil {
		return nil, "", err
	}
	return r.listEmployeeCandidates(ctx, "department_id", deptID, notExistsDeptResponsible, searchQuery, after, limit)
}
