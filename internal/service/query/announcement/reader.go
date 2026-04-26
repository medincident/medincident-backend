// Package announcement is the query-side reader for domain.announcements
// and projections.announcement_views. Visibility is enforced via inline
// SQL — there are no authz.Require calls on the read path.
//
// See: docs/services/Announcements.md
package announcement

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/query"
)

const (
	ErrCodeAnnouncementNotFound   = "announcement_query_not_found"
	ErrCodeAnnouncementReadFailed = "announcement_query_read_failed"
	ErrCodeAnnouncementBadCursor  = "announcement_query_bad_cursor"
)

const readerScope = "services.query.announcement"

// Reader is the query-side service for announcements.
type Reader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewReader wires the reader.
func NewReader(db *gorm.DB, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, logger: logger}
}

// AnnouncementRow is the projection returned from DB queries.
type AnnouncementRow struct {
	ID             uuid.UUID     `gorm:"column:id"`
	OrganizationID uuid.UUID     `gorm:"column:organization_id"`
	ClinicID       uuid.NullUUID `gorm:"column:clinic_id"`
	DepartmentID   uuid.NullUUID `gorm:"column:department_id"`
	AuthorID       string        `gorm:"column:author_id"`
	Title          string        `gorm:"column:title"`
	Content        string        `gorm:"column:content"`
	Priority       string        `gorm:"column:priority"`
	IsArchived     bool          `gorm:"column:is_archived"`
	StartsAt       null.Time     `gorm:"column:starts_at"`
	EndsAt         null.Time     `gorm:"column:ends_at"`
	CreatedAt      time.Time     `gorm:"column:created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at"`
	ViewCount      int64         `gorm:"column:view_count"`
}

// callerCtx carries the minimal resolved identity for announcement visibility checks.
type callerCtx struct {
	zitadelID     string
	isSystemAdmin bool
	// orgAdmin[orgID] = true if caller is OrgAdmin (direct or deputy) of that org.
	orgAdmin map[uuid.UUID]bool
	// clinicHead[clinicID] = true if caller is ClinicHead of that clinic.
	clinicHead map[uuid.UUID]bool
	// deptResp[deptID] = true if caller is DeptResponsible of that dept.
	deptResp map[uuid.UUID]bool
	// employee represents the caller's employee row, if any.
	employeeOrgID  uuid.NullUUID
	employeeClinic uuid.NullUUID
	employeeDept   uuid.NullUUID
}

const activeVacation = `v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`

// resolveCaller runs one query per role table to build callerCtx.
func (r *Reader) resolveCaller(ctx context.Context, callerID string) (*callerCtx, error) {
	cc := &callerCtx{
		zitadelID:  callerID,
		orgAdmin:   map[uuid.UUID]bool{},
		clinicHead: map[uuid.UUID]bool{},
		deptResp:   map[uuid.UUID]bool{},
	}
	tx := r.db.WithContext(ctx)

	// System admin check.
	var saCount int64
	if err := tx.Raw(`SELECT COUNT(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, callerID).
		Scan(&saCount).Error; err != nil {
		return nil, wrapRead(err, "system_admin lookup")
	}
	cc.isSystemAdmin = saCount > 0

	// Employee row for scope-based visibility.
	type empRow struct {
		OrgID    uuid.NullUUID `gorm:"column:organization_id"`
		ClinicID uuid.NullUUID `gorm:"column:clinic_id"`
		DeptID   uuid.NullUUID `gorm:"column:department_id"`
	}
	var er empRow
	if err := tx.Raw(`
		SELECT e.organization_id, d.clinic_id, e.department_id
		FROM domain.employees e
		JOIN domain.departments d ON d.id = e.department_id
		WHERE e.zitadel_user_id = ?
		LIMIT 1`, callerID).Scan(&er).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, wrapRead(err, "employee lookup")
	}
	cc.employeeOrgID = er.OrgID
	cc.employeeClinic = er.ClinicID
	cc.employeeDept = er.DeptID

	// OrgAdmin scopes.
	if err := r.collectScopes(tx, callerID, `
		SELECT t.organization_id AS scope_id FROM domain.org_admins t
		JOIN domain.employees e ON e.id = t.employee_id
		WHERE e.zitadel_user_id = ?
		UNION
		SELECT t.organization_id FROM domain.org_admins t
		JOIN domain.employees e ON e.id = t.deputy_employee_id
		JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
		WHERE e.zitadel_user_id = ? AND `+activeVacation,
		cc.orgAdmin); err != nil {
		return nil, err
	}

	// ClinicHead scopes.
	if err := r.collectScopes(tx, callerID, `
		SELECT t.clinic_id AS scope_id FROM domain.clinic_heads t
		JOIN domain.employees e ON e.id = t.employee_id
		WHERE e.zitadel_user_id = ?
		UNION
		SELECT t.clinic_id FROM domain.clinic_heads t
		JOIN domain.employees e ON e.id = t.deputy_employee_id
		JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
		WHERE e.zitadel_user_id = ? AND `+activeVacation,
		cc.clinicHead); err != nil {
		return nil, err
	}

	// DeptResponsible scopes.
	if err := r.collectScopes(tx, callerID, `
		SELECT t.department_id AS scope_id FROM domain.department_responsibles t
		JOIN domain.employees e ON e.id = t.employee_id
		WHERE e.zitadel_user_id = ?
		UNION
		SELECT t.department_id FROM domain.department_responsibles t
		JOIN domain.employees e ON e.id = t.deputy_employee_id
		JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
		WHERE e.zitadel_user_id = ? AND `+activeVacation,
		cc.deptResp); err != nil {
		return nil, err
	}

	return cc, nil
}

func (r *Reader) collectScopes(tx *gorm.DB, callerID, rawSQL string, dst map[uuid.UUID]bool) error {
	var rows []struct {
		ScopeID uuid.UUID `gorm:"column:scope_id"`
	}
	if err := tx.Raw(rawSQL, callerID, callerID).Scan(&rows).Error; err != nil {
		return wrapRead(err, "scope collection")
	}
	for _, row := range rows {
		dst[row.ScopeID] = true
	}
	return nil
}

// canManage reports whether the caller may see archived announcements
// and manage the given scope.
func (cc *callerCtx) canManage(orgID uuid.UUID, clinicID, deptID uuid.NullUUID) bool {
	if cc.isSystemAdmin || cc.orgAdmin[orgID] {
		return true
	}
	if clinicID.Valid && cc.clinicHead[clinicID.UUID] {
		return true
	}
	if deptID.Valid && cc.deptResp[deptID.UUID] {
		return true
	}
	return false
}

// ListFilter holds caller-supplied list parameters.
type ListFilter struct {
	Priority        string // "" = all
	IncludeArchived bool
	Limit           int
	Cursor          *string
}

// cursor is the keyset pagination token.
type cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

func encodeCursor(c cursor) string {
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func decodeCursor(s string) (cursor, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return cursor{}, err
	}
	var c cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return cursor{}, err
	}
	return c, nil
}

func normalizeLimit(l int) int {
	if l <= 0 {
		return query.DefaultLimit
	}
	if l > query.MaxLimit {
		return query.MaxLimit
	}
	return l
}

// GetAnnouncement returns an announcement by ID and increments its view counter.
//
// See: docs/services/Announcements.md
func (r *Reader) GetAnnouncement(ctx context.Context, callerID string, id uuid.UUID) (*AnnouncementRow, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return nil, err
	}

	var row AnnouncementRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT a.*, COALESCE(v.view_count, 0) AS view_count
		FROM domain.announcements a
		LEFT JOIN projections.announcement_views v ON v.announcement_id = a.id
		WHERE a.id = ?`, id).Scan(&row).Error; err != nil {
		return nil, wrapRead(err, "get announcement")
	}
	if row.ID == uuid.Nil {
		return nil, oops.In(readerScope).
			Code(ErrCodeAnnouncementNotFound).
			Public("Announcement not found.").
			With("announcement_id", id).
			Errorf("not found")
	}

	// Visibility: admins see everything; employees only see active announcements
	// at their scope.
	if !cc.canManage(row.OrganizationID, row.ClinicID, row.DepartmentID) {
		if !isVisible(&row, cc) {
			return nil, oops.In(readerScope).
				Code(ErrCodeAnnouncementNotFound).
				Public("Announcement not found.").
				With("announcement_id", id).
				Errorf("not visible")
		}
	}

	// Increment view counter.
	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO projections.announcement_views (announcement_id, view_count)
		VALUES (?, 1)
		ON CONFLICT (announcement_id) DO UPDATE SET view_count = projections.announcement_views.view_count + 1`,
		id,
	).Error; err != nil {
		r.logger.Warn().Err(err).Str("announcement_id", id.String()).Msg("failed to increment view count")
	} else {
		row.ViewCount++
	}

	return &row, nil
}

// isVisible checks whether an active (non-admin) caller may see the announcement.
func isVisible(row *AnnouncementRow, cc *callerCtx) bool {
	if row.IsArchived {
		return false
	}
	now := time.Now()
	if row.StartsAt.Valid && row.StartsAt.Time.After(now) {
		return false
	}
	if row.EndsAt.Valid && !row.EndsAt.Time.After(now) {
		return false
	}
	// Scope check: org-level announcement visible to any org member.
	if !row.ClinicID.Valid {
		return cc.employeeOrgID.Valid && cc.employeeOrgID.UUID == row.OrganizationID
	}
	// Clinic-level: visible to clinic members.
	if !row.DepartmentID.Valid {
		return cc.employeeClinic.Valid && cc.employeeClinic.UUID == row.ClinicID.UUID
	}
	// Dept-level: visible to dept members.
	return cc.employeeDept.Valid && cc.employeeDept.UUID == row.DepartmentID.UUID
}

// ListResult is the paginated result of a list query.
type ListResult struct {
	Items      []AnnouncementRow
	NextCursor *string
}

// ListForOrganization returns announcements scoped to an organization.
//
// See: docs/services/Announcements.md
func (r *Reader) ListForOrganization(
	ctx context.Context, callerID string, orgID uuid.UUID, f ListFilter,
) (ListResult, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return ListResult{}, err
	}
	canManage := cc.canManage(orgID, uuid.NullUUID{}, uuid.NullUUID{})
	// Employees must belong to this org.
	if !canManage && (!cc.employeeOrgID.Valid || cc.employeeOrgID.UUID != orgID) {
		return ListResult{Items: []AnnouncementRow{}}, nil
	}
	where := "a.organization_id = ? AND a.clinic_id IS NULL AND a.department_id IS NULL"
	args := []any{orgID}
	return r.listAnnouncements(ctx, canManage, f, where, args)
}

// ListForClinic returns announcements scoped to a clinic (and its org-level announcements).
//
// See: docs/services/Announcements.md
func (r *Reader) ListForClinic(
	ctx context.Context, callerID string, clinicID uuid.UUID, f ListFilter,
) (ListResult, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return ListResult{}, err
	}

	// Resolve org from clinic.
	var orgID uuid.UUID
	if err := r.db.WithContext(ctx).
		Raw(`SELECT organization_id FROM domain.clinics WHERE id = ? LIMIT 1`, clinicID).
		Row().Scan(&orgID); err != nil || orgID == uuid.Nil {
		return ListResult{}, oops.In(readerScope).
			Code(ErrCodeAnnouncementReadFailed).
			Public("Clinic not found.").
			With("clinic_id", clinicID).
			Wrap(err)
	}

	clinicNullUUID := uuid.NullUUID{UUID: clinicID, Valid: true}
	canManage := cc.canManage(orgID, clinicNullUUID, uuid.NullUUID{})
	if !canManage && (!cc.employeeClinic.Valid || cc.employeeClinic.UUID != clinicID) {
		return ListResult{Items: []AnnouncementRow{}}, nil
	}
	// Returns org-level + clinic-level announcements for this clinic.
	where := `(a.organization_id = ? AND a.clinic_id IS NULL AND a.department_id IS NULL)
		OR (a.clinic_id = ? AND a.department_id IS NULL)`
	args := []any{orgID, clinicID}
	return r.listAnnouncements(ctx, canManage, f, where, args)
}

// ListForDepartment returns announcements scoped to a department (plus clinic + org).
//
// See: docs/services/Announcements.md
func (r *Reader) ListForDepartment(
	ctx context.Context, callerID string, deptID uuid.UUID, f ListFilter,
) (ListResult, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return ListResult{}, err
	}

	type deptRow struct {
		ClinicID uuid.UUID `gorm:"column:clinic_id"`
	}
	var dr deptRow
	if err := r.db.WithContext(ctx).
		Raw(`SELECT clinic_id FROM domain.departments WHERE id = ? LIMIT 1`, deptID).
		Scan(&dr).Error; err != nil || dr.ClinicID == uuid.Nil {
		return ListResult{}, oops.In(readerScope).
			Code(ErrCodeAnnouncementReadFailed).
			Public("Department not found.").
			With("department_id", deptID).
			Wrap(err)
	}
	clinicID := dr.ClinicID

	var orgID uuid.UUID
	if err := r.db.WithContext(ctx).
		Raw(`SELECT organization_id FROM domain.clinics WHERE id = ? LIMIT 1`, clinicID).
		Row().Scan(&orgID); err != nil || orgID == uuid.Nil {
		return ListResult{}, oops.In(readerScope).Code(ErrCodeAnnouncementReadFailed).Wrap(err)
	}

	deptNullUUID := uuid.NullUUID{UUID: deptID, Valid: true}
	clinicNullUUID := uuid.NullUUID{UUID: clinicID, Valid: true}
	canManage := cc.canManage(orgID, clinicNullUUID, deptNullUUID)
	if !canManage && (!cc.employeeDept.Valid || cc.employeeDept.UUID != deptID) {
		return ListResult{Items: []AnnouncementRow{}}, nil
	}
	// Returns org-level + clinic-level + dept-level announcements.
	where := `(a.organization_id = ? AND a.clinic_id IS NULL AND a.department_id IS NULL)
		OR (a.clinic_id = ? AND a.department_id IS NULL)
		OR a.department_id = ?`
	args := []any{orgID, clinicID, deptID}
	return r.listAnnouncements(ctx, canManage, f, where, args)
}

// listAnnouncements is the shared list implementation.
func (r *Reader) listAnnouncements(
	ctx context.Context, canManage bool, f ListFilter, scopeWhere string, scopeArgs []any,
) (ListResult, error) {
	limit := normalizeLimit(f.Limit)

	clauses := []string{"(" + scopeWhere + ")"}
	args := append([]any(nil), scopeArgs...)

	if !canManage || !f.IncludeArchived {
		clauses = append(clauses,
			"a.is_archived = false",
			"(a.starts_at IS NULL OR a.starts_at <= now())",
			"(a.ends_at IS NULL OR a.ends_at >= now())",
		)
	}

	if f.Priority != "" {
		clauses = append(clauses, "a.priority = ?")
		args = append(args, f.Priority)
	}

	if f.Cursor != nil {
		c, err := decodeCursor(*f.Cursor)
		if err != nil {
			return ListResult{}, oops.In(readerScope).
				Code(ErrCodeAnnouncementBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		clauses = append(clauses, "(a.created_at, a.id) < (?, ?)")
		args = append(args, c.CreatedAt, c.ID)
	}

	whereSQL := ""
	for i, c := range clauses {
		if i == 0 {
			whereSQL = "WHERE " + c
		} else {
			whereSQL += " AND " + c
		}
	}

	rawSQL := fmt.Sprintf(`
		SELECT a.*, COALESCE(v.view_count, 0) AS view_count
		FROM domain.announcements a
		LEFT JOIN projections.announcement_views v ON v.announcement_id = a.id
		%s
		ORDER BY a.created_at DESC, a.id DESC
		LIMIT ?`, whereSQL)

	args = append(args, limit+1)

	var rows []AnnouncementRow
	if err := r.db.WithContext(ctx).Raw(rawSQL, args...).Scan(&rows).Error; err != nil {
		return ListResult{}, wrapRead(err, "list announcements")
	}

	var nextCursor *string
	if len(rows) > limit {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		s := encodeCursor(cursor{CreatedAt: last.CreatedAt, ID: last.ID})
		nextCursor = &s
	}

	return ListResult{Items: rows, NextCursor: nextCursor}, nil
}

func wrapRead(err error, label string) error {
	return oops.In(readerScope).
		Code(ErrCodeAnnouncementReadFailed).
		With("label", label).
		Wrap(err)
}

// NullTimePtr converts null.Time to *string (RFC3339Nano) for proto mapping.
func NullTimePtr(t null.Time) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.UTC().Format(time.RFC3339Nano)
	return &s
}
