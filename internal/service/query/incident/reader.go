// Package incident is the query-side reader for projections.incidents
// and the two history projections. Visibility is enforced inside SQL
// based on the caller's role context — there is no Authz.Require call
// for read paths.
package incident

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/cursor"
	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query"
)

const (
	ErrCodeIncidentReadFailed = "incident_query_read_failed"
	ErrCodeIncidentNotFound   = "incident_query_not_found"
	ErrCodeListBadCursor      = "incident_bad_cursor"
)

// IncidentListResult is returned by paginated incident list methods.
type IncidentListResult struct {
	Items      []IncidentView
	NextCursor *string
}

const scope = "services.query.incident"

// Reader is the query-side service. No interfaces.
type Reader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given db and logger.
func NewReader(db *gorm.DB, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, logger: logger}
}

// callerContext describes which slice of the corpus the caller may see.
type callerContext struct {
	zitadelID     string
	isSystemAdmin bool
	// Per-org capabilities. The maps are keyed by organization_id.
	orgAdmin      map[uuid.UUID]bool
	orgHead       map[uuid.UUID]bool
	orgDispatcher map[uuid.UUID]bool
	clinicHead    map[uuid.UUID]bool // by clinic_id
	deptResp      map[uuid.UUID]bool // by department_id
	// employee_id (if any). Used for "their own" filtering.
	employeeID uuid.NullUUID
}

// CallerContext is the public alias of the resolved caller used by the
// buffer reader. It exposes only the predicates the buffer reader needs.
type CallerContext = callerContext

// ResolveCaller is the public wrapper for resolveCaller.
//
// See: docs/services/incident/Incidents.md
func (r *Reader) ResolveCaller(ctx context.Context, callerID string) (*CallerContext, error) {
	return r.resolveCaller(ctx, callerID)
}

// IsSystemAdmin reports whether the caller is a system admin.
func (cc *callerContext) IsSystemAdmin() bool { return cc.isSystemAdmin }

// ZitadelID returns the caller's zitadel id.
func (cc *callerContext) ZitadelID() string { return cc.zitadelID }

// IsOrgAdminOf reports whether the caller is an OrgAdmin of orgID.
func (cc *callerContext) IsOrgAdminOf(orgID uuid.UUID) bool { return cc.orgAdmin[orgID] }

// IsOrgDispatcherOf reports whether the caller is an OrgDispatcher of orgID.
func (cc *callerContext) IsOrgDispatcherOf(orgID uuid.UUID) bool { return cc.orgDispatcher[orgID] }

// CanSeeBufferForOrg reports whether the caller may list buffer entries for orgID.
// SystemAdmin, OrgAdmin, OrgDispatcher only.
func (cc *callerContext) CanSeeBufferForOrg(orgID uuid.UUID) bool {
	return cc.isSystemAdmin || cc.orgAdmin[orgID] || cc.orgDispatcher[orgID]
}

// IsPatient reports whether the caller is treated as a patient — i.e.
// has no employee row at all.
func (cc *callerContext) IsPatient() bool {
	return !cc.isSystemAdmin && !cc.employeeID.Valid &&
		len(cc.orgAdmin) == 0 && len(cc.orgHead) == 0 &&
		len(cc.orgDispatcher) == 0 && len(cc.clinicHead) == 0 && len(cc.deptResp) == 0
}

// resolveCaller assembles the role context with one query per role
// table. All checks are vacation-aware.
func (r *Reader) resolveCaller(ctx context.Context, callerID string) (*callerContext, error) {
	cc := &callerContext{
		zitadelID:     callerID,
		orgAdmin:      map[uuid.UUID]bool{},
		orgHead:       map[uuid.UUID]bool{},
		orgDispatcher: map[uuid.UUID]bool{},
		clinicHead:    map[uuid.UUID]bool{},
		deptResp:      map[uuid.UUID]bool{},
	}
	tx := r.db.WithContext(ctx)

	// SystemAdmin
	var saCount int64
	if err := tx.Raw(
		`SELECT COUNT(*) FROM domain.system_admins WHERE zitadel_user_id = ?`, callerID,
	).Scan(&saCount).Error; err != nil {
		return nil, wrapRead(err, "system_admin lookup")
	}
	cc.isSystemAdmin = saCount > 0

	// Employee identity (any org). Used for ListMyIncidents & cancelled-self filter.
	var emp model.Employee
	if err := tx.Where("zitadel_user_id = ?", callerID).Limit(1).First(&emp).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, wrapRead(err, "employee lookup")
		}
	} else {
		cc.employeeID = uuid.NullUUID{UUID: emp.ID, Valid: true}
	}

	// Per-role queries. Each branch hard-codes the role's table + scope
	// column so the SQL is never built from interpolated identifiers.
	if err := r.collectScopes(tx, callerID, roleSrcOrgAdmin, cc.orgAdmin); err != nil {
		return nil, err
	}
	if err := r.collectScopes(tx, callerID, roleSrcOrgHead, cc.orgHead); err != nil {
		return nil, err
	}
	if err := r.collectScopes(tx, callerID, roleSrcOrgDispatcher, cc.orgDispatcher); err != nil {
		return nil, err
	}
	if err := r.collectScopes(tx, callerID, roleSrcClinicHead, cc.clinicHead); err != nil {
		return nil, err
	}
	if err := r.collectScopes(tx, callerID, roleSrcDeptResp, cc.deptResp); err != nil {
		return nil, err
	}
	return cc, nil
}

// roleScopeSource identifies a (table, scope-column) pair for
// collectScopes. Using a typed enum instead of raw strings means the
// SQL is hard-coded per branch — there is no path by which a caller
// can inject identifiers, even by accident.
type roleScopeSource int

const (
	roleSrcOrgAdmin roleScopeSource = iota
	roleSrcOrgHead
	roleSrcOrgDispatcher
	roleSrcClinicHead
	roleSrcDeptResp
)

// roleScopeQuery returns the fully-baked UNION SQL and a label for
// error wrapping for the given role source. All identifiers are
// compile-time literals; only the caller id is bound as a parameter.
func roleScopeQuery(src roleScopeSource) (rawSQL, label string) {
	switch src {
	case roleSrcOrgAdmin:
		return `SELECT t.organization_id FROM domain.org_admins t
			JOIN domain.employees e ON e.id = t.employee_id
			WHERE e.zitadel_user_id = ?
			UNION
			SELECT t.organization_id FROM domain.org_admins t
			JOIN domain.employees e ON e.id = t.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
			WHERE e.zitadel_user_id = ?
			  AND v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`, "domain.org_admins"
	case roleSrcOrgHead:
		return `SELECT t.organization_id FROM domain.org_heads t
			JOIN domain.employees e ON e.id = t.employee_id
			WHERE e.zitadel_user_id = ?
			UNION
			SELECT t.organization_id FROM domain.org_heads t
			JOIN domain.employees e ON e.id = t.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
			WHERE e.zitadel_user_id = ?
			  AND v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`, "domain.org_heads"
	case roleSrcOrgDispatcher:
		return `SELECT t.organization_id FROM domain.org_dispatchers t
			JOIN domain.employees e ON e.id = t.employee_id
			WHERE e.zitadel_user_id = ?
			UNION
			SELECT t.organization_id FROM domain.org_dispatchers t
			JOIN domain.employees e ON e.id = t.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
			WHERE e.zitadel_user_id = ?
			  AND v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`, "domain.org_dispatchers"
	case roleSrcClinicHead:
		return `SELECT t.clinic_id FROM domain.clinic_heads t
			JOIN domain.employees e ON e.id = t.employee_id
			WHERE e.zitadel_user_id = ?
			UNION
			SELECT t.clinic_id FROM domain.clinic_heads t
			JOIN domain.employees e ON e.id = t.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
			WHERE e.zitadel_user_id = ?
			  AND v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`, "domain.clinic_heads"
	case roleSrcDeptResp:
		return `SELECT t.department_id FROM domain.department_responsibles t
			JOIN domain.employees e ON e.id = t.employee_id
			WHERE e.zitadel_user_id = ?
			UNION
			SELECT t.department_id FROM domain.department_responsibles t
			JOIN domain.employees e ON e.id = t.deputy_employee_id
			JOIN domain.employee_vacations v ON v.employee_id = t.employee_id
			WHERE e.zitadel_user_id = ?
			  AND v.starts_at <= now() AND (v.ends_at IS NULL OR v.ends_at > now())`, "domain.department_responsibles"
	}
	// Unreachable: the enum is closed at compile time.
	return "", ""
}

// collectScopes runs the pre-baked SQL for the given role source and
// fills `out` with every scope id the caller holds (directly or as
// a deputy with an active vacation on the principal).
func (r *Reader) collectScopes(
	tx *gorm.DB, callerID string, src roleScopeSource, out map[uuid.UUID]bool,
) error {
	q, label := roleScopeQuery(src)
	rows, err := tx.Raw(q, callerID, callerID).Rows()
	if err != nil {
		return wrapRead(err, "scope lookup "+label)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return wrapRead(err, "scope scan "+label)
		}
		out[id] = true
	}
	if err := rows.Err(); err != nil {
		return wrapRead(err, "iterate scope rows "+label)
	}
	return nil
}

func wrapRead(err error, action string) error {
	return oops.In(scope).Code(ErrCodeIncidentReadFailed).
		With("action", action).Wrap(err)
}

// IncidentView is the projection row returned to handlers.
type IncidentView struct {
	ID                         uuid.UUID
	OrganizationID             uuid.UUID
	ClinicID                   uuid.UUID
	DepartmentID               uuid.UUID
	CategoryID                 uuid.UUID
	TypeID                     uuid.UUID
	Status                     model.IncidentStatus
	Priority                   model.IncidentPriority
	Description                null.String
	PatientOriginalDescription null.String
	OccurredAt                 time.Time
	RegistrarEmployeeID        uuid.UUID
	RegistrarDisplayName       string
	RegistrarPosition          null.String
	RegistrarOrganizationID    uuid.UUID
	RegistrarClinicID          uuid.UUID
	RegistrarDepartmentID      uuid.UUID
	SourcePatientZitadelUserID null.String
	SourceBufferID             uuid.NullUUID
	ReopenedFromIncidentID     uuid.NullUUID
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	// PatientPerspective is true when the caller is being served as a patient
	// (non-employee) — handlers redact accordingly.
	PatientPerspective bool
}

const selectColumns = `id, organization_id, clinic_id, department_id, category_id, type_id,
	status, priority, description, patient_original_description, occurred_at,
	registrar_employee_id, registrar_display_name, registrar_position,
	registrar_organization_id, registrar_clinic_id, registrar_department_id,
	source_patient_zitadel_user_id, source_buffer_id, reopened_from_incident_id,
	created_at, updated_at`

func scanIncident(row interface{ Scan(...any) error }, v *IncidentView) error {
	return row.Scan(
		&v.ID, &v.OrganizationID, &v.ClinicID, &v.DepartmentID,
		&v.CategoryID, &v.TypeID, &v.Status, &v.Priority,
		&v.Description, &v.PatientOriginalDescription, &v.OccurredAt,
		&v.RegistrarEmployeeID, &v.RegistrarDisplayName, &v.RegistrarPosition,
		&v.RegistrarOrganizationID, &v.RegistrarClinicID, &v.RegistrarDepartmentID,
		&v.SourcePatientZitadelUserID, &v.SourceBufferID, &v.ReopenedFromIncidentID,
		&v.CreatedAt, &v.UpdatedAt,
	)
}

// GetIncident loads one incident if the caller is allowed to see it.
//
// See: docs/services/incident/Incidents.md
func (r *Reader) GetIncident(ctx context.Context, callerID string, id uuid.UUID) (*IncidentView, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return nil, err
	}
	return r.getIncidentForCaller(ctx, cc, id)
}

// getIncidentForCaller is the visibility-filtered single-incident
// fetch shared by GetIncident and GetIncidentHistory. Internal callers
// pass a pre-resolved callerContext so the role-resolution queries
// (six SQL round-trips) are not repeated.
func (r *Reader) getIncidentForCaller(
	ctx context.Context, cc *callerContext, id uuid.UUID,
) (*IncidentView, error) {
	where, args := r.visibilityClause(cc, "")
	args = append([]any{id}, args...)
	q := `SELECT ` + selectColumns + ` FROM projections.incidents
		WHERE id = ? AND (` + where + `) LIMIT 1`
	var v IncidentView
	row := r.db.WithContext(ctx).Raw(q, args...).Row()
	if err := scanIncident(row, &v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In(scope).Code(ErrCodeIncidentNotFound).
				Public("Incident not found.").With("incident_id", id).Wrap(err)
		}
		return nil, wrapRead(err, "scan incident")
	}
	v.PatientPerspective = cc.IsPatient()
	return &v, nil
}

// visibilityClause assembles the OR-of-OR WHERE fragment based on
// the caller's role context.
func (r *Reader) visibilityClause(cc *callerContext, prefix string) (clause string, args []any) {
	col := func(c string) string {
		if prefix == "" {
			return c
		}
		return prefix + c
	}
	if cc.isSystemAdmin {
		return "TRUE", nil
	}
	var parts []string

	for orgID := range cc.orgAdmin {
		parts = append(parts, col("organization_id")+" = ?")
		args = append(args, orgID)
	}
	if len(cc.orgHead) > 0 || len(cc.orgDispatcher) > 0 {
		set := map[uuid.UUID]bool{}
		for k := range cc.orgHead {
			set[k] = true
		}
		for k := range cc.orgDispatcher {
			set[k] = true
		}
		ids := make([]uuid.UUID, 0, len(set))
		for k := range set {
			ids = append(ids, k)
		}
		if len(ids) > 0 {
			placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
			parts = append(parts,
				"("+col("organization_id")+" IN ("+placeholders+") AND "+col("status")+" != ?)")
			for _, id := range ids {
				args = append(args, id)
			}
			args = append(args, model.IncidentStatusCancelled)
		}
	}
	for clinicID := range cc.clinicHead {
		parts = append(parts, "("+col("clinic_id")+" = ? AND "+col("status")+" != ?)")
		args = append(args, clinicID, model.IncidentStatusCancelled)
	}
	for deptID := range cc.deptResp {
		parts = append(parts, "("+col("department_id")+" = ? AND "+col("status")+" != ?)")
		args = append(args, deptID, model.IncidentStatusCancelled)
	}
	if cc.employeeID.Valid {
		parts = append(parts, col("registrar_employee_id")+" = ?")
		args = append(args, cc.employeeID.UUID)
	}
	parts = append(parts, col("source_patient_zitadel_user_id")+" = ?")
	args = append(args, cc.zitadelID)

	if len(parts) == 0 {
		return "FALSE", nil
	}
	return "(" + strings.Join(parts, " OR ") + ")", args
}

// ListFilters captures the optional list filters.
type ListFilters struct {
	Statuses     []model.IncidentStatus
	Priorities   []model.IncidentPriority
	ClinicID     uuid.NullUUID
	DepartmentID uuid.NullUUID
	CategoryID   uuid.NullUUID
	TypeID       uuid.NullUUID
	OccurredFrom *time.Time
	OccurredTo   *time.Time
	Limit        int
	After        *string
}

// ListIncidents returns incidents in an organization the caller may see.
//
// See: docs/services/incident/Incidents.md
func (r *Reader) ListIncidents(
	ctx context.Context, callerID string, orgID uuid.UUID, f *ListFilters,
) (IncidentListResult, error) {
	if f == nil {
		f = &ListFilters{}
	}
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return IncidentListResult{}, err
	}
	where, visArgs := r.visibilityClause(cc, "")
	conds := []string{"organization_id = ?", where}
	args := make([]any, 0, len(visArgs)+12)
	args = append(args, orgID)
	args = append(args, visArgs...)

	if len(f.Statuses) > 0 {
		ph := strings.TrimSuffix(strings.Repeat("?,", len(f.Statuses)), ",")
		conds = append(conds, "status IN ("+ph+")")
		for _, s := range f.Statuses {
			args = append(args, s)
		}
	}
	if len(f.Priorities) > 0 {
		ph := strings.TrimSuffix(strings.Repeat("?,", len(f.Priorities)), ",")
		conds = append(conds, "priority IN ("+ph+")")
		for _, p := range f.Priorities {
			args = append(args, p)
		}
	}
	if f.ClinicID.Valid {
		conds = append(conds, "clinic_id = ?")
		args = append(args, f.ClinicID.UUID)
	}
	if f.DepartmentID.Valid {
		conds = append(conds, "department_id = ?")
		args = append(args, f.DepartmentID.UUID)
	}
	if f.CategoryID.Valid {
		conds = append(conds, "category_id = ?")
		args = append(args, f.CategoryID.UUID)
	}
	if f.TypeID.Valid {
		conds = append(conds, "type_id = ?")
		args = append(args, f.TypeID.UUID)
	}
	if f.OccurredFrom != nil {
		conds = append(conds, "occurred_at >= ?")
		args = append(args, *f.OccurredFrom)
	}
	if f.OccurredTo != nil {
		conds = append(conds, "occurred_at <= ?")
		args = append(args, *f.OccurredTo)
	}
	limit := normLimit(f.Limit)
	if f.After != nil {
		t, idStr, err := cursor.Decode(*f.After)
		if err != nil {
			return IncidentListResult{}, oops.In(scope).
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		incID, _ := uuid.Parse(idStr)
		conds = append(conds, "(updated_at, id) < (?, ?)")
		args = append(args, t, incID)
	}
	args = append(args, limit+1)

	q := `SELECT ` + selectColumns + ` FROM projections.incidents WHERE ` +
		strings.Join(conds, " AND ") + ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return IncidentListResult{}, wrapRead(err, "list incidents")
	}
	defer func() { _ = rows.Close() }()
	out := make([]IncidentView, 0, limit+1)
	for rows.Next() {
		var v IncidentView
		if err := scanIncident(rows, &v); err != nil {
			return IncidentListResult{}, wrapRead(err, "scan incident row")
		}
		v.PatientPerspective = cc.IsPatient()
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return IncidentListResult{}, wrapRead(err, "iterate incident rows")
	}
	var nextCursor *string
	if len(out) > limit {
		out = out[:limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return IncidentListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListMyIncidents returns incidents where the caller is the registrar
// (employees) OR linked to the caller's buffer submissions (patients).
//
// See: docs/services/incident/Incidents.md
func (r *Reader) ListMyIncidents(
	ctx context.Context, callerID string, limit int, after *string,
) (IncidentListResult, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return IncidentListResult{}, err
	}
	limit = normLimit(limit)

	conds := []string{}
	args := []any{}
	if cc.employeeID.Valid {
		conds = append(conds, "registrar_employee_id = ?")
		args = append(args, cc.employeeID.UUID)
	}
	conds = append(conds, "source_patient_zitadel_user_id = ?")
	args = append(args, callerID)

	cursorConds := ""
	if after != nil {
		t, idStr, err := cursor.Decode(*after)
		if err != nil {
			return IncidentListResult{}, oops.In(scope).
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		incID, _ := uuid.Parse(idStr)
		cursorConds = ` AND (updated_at, id) < (?, ?)`
		args = append(args, t, incID)
	}
	args = append(args, limit+1)

	q := `SELECT ` + selectColumns + ` FROM projections.incidents WHERE (` +
		strings.Join(conds, " OR ") + `)` + cursorConds +
		` ORDER BY updated_at DESC, id DESC LIMIT ?`

	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return IncidentListResult{}, wrapRead(err, "list my incidents")
	}
	defer func() { _ = rows.Close() }()
	out := make([]IncidentView, 0, limit+1)
	for rows.Next() {
		var v IncidentView
		if err := scanIncident(rows, &v); err != nil {
			return IncidentListResult{}, wrapRead(err, "scan my incident row")
		}
		v.PatientPerspective = cc.IsPatient()
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return IncidentListResult{}, wrapRead(err, "iterate my incident rows")
	}
	var nextCursor *string
	if len(out) > limit {
		out = out[:limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return IncidentListResult{Items: out, NextCursor: nextCursor}, nil
}

// NormLimit clamps a caller-supplied limit to safe bounds.
// Zero or out-of-range values fall back to query.DefaultLimit.
// Exported so the buffer reader and other readers can share the impl.
func NormLimit(limit int) int {
	if limit <= 0 || limit > query.MaxLimit {
		return query.DefaultLimit
	}
	return limit
}

// normLimit is the unexported alias for call sites in this package.
func normLimit(limit int) int { return NormLimit(limit) }

// StatusHistoryEntry mirrors a row of projections.incident_status_history.
type StatusHistoryEntry struct {
	ID               uuid.UUID
	OldStatus        null.String // nullable enum
	NewStatus        model.IncidentStatus
	ActorEmployeeID  uuid.NullUUID
	ActorDisplayName null.String
	ChangedAt        time.Time
}

// PriorityHistoryEntry mirrors a row of projections.incident_priority_history.
type PriorityHistoryEntry struct {
	ID               uuid.UUID
	OldPriority      model.IncidentPriority
	NewPriority      model.IncidentPriority
	ActorEmployeeID  uuid.NullUUID
	ActorDisplayName null.String
	ChangedAt        time.Time
}

// IncidentHistory bundles both timelines.
type IncidentHistory struct {
	Status   []StatusHistoryEntry
	Priority []PriorityHistoryEntry
}

// GetIncidentHistory returns both timelines for one incident, only if
// the caller can see the incident at all. Patients are NOT allowed to
// view history (returns permission_denied).
//
// See: docs/services/incident/Incidents.md
func (r *Reader) GetIncidentHistory(
	ctx context.Context, callerID string, incidentID uuid.UUID,
) (*IncidentHistory, error) {
	cc, err := r.resolveCaller(ctx, callerID)
	if err != nil {
		return nil, err
	}
	if cc.IsPatient() {
		return nil, oops.In(scope).
			Code(authz.ErrCodePermissionDenied).
			Public("History is not available for patients.").
			Errorf("patient denied")
	}
	// Reuse the already-resolved callerContext to avoid a second
	// round of role-lookup queries inside GetIncident.
	if _, err := r.getIncidentForCaller(ctx, cc, incidentID); err != nil {
		return nil, err
	}

	tx := r.db.WithContext(ctx)
	statusRows, err := tx.Raw(
		`SELECT id, old_status, new_status, actor_employee_id, actor_display_name, changed_at
		 FROM projections.incident_status_history
		 WHERE incident_id = ? ORDER BY changed_at ASC`, incidentID,
	).Rows()
	if err != nil {
		return nil, wrapRead(err, "list status history")
	}
	defer func() { _ = statusRows.Close() }()
	statuses := []StatusHistoryEntry{}
	for statusRows.Next() {
		var e StatusHistoryEntry
		var rawNew string
		if err := statusRows.Scan(&e.ID, &e.OldStatus, &rawNew,
			&e.ActorEmployeeID, &e.ActorDisplayName, &e.ChangedAt); err != nil {
			return nil, wrapRead(err, "scan status history row")
		}
		e.NewStatus = model.IncidentStatus(rawNew)
		statuses = append(statuses, e)
	}
	if err := statusRows.Err(); err != nil {
		return nil, wrapRead(err, "iterate status history")
	}

	prioRows, err := tx.Raw(
		`SELECT id, old_priority, new_priority, actor_employee_id, actor_display_name, changed_at
		 FROM projections.incident_priority_history
		 WHERE incident_id = ? ORDER BY changed_at ASC`, incidentID,
	).Rows()
	if err != nil {
		return nil, wrapRead(err, "list priority history")
	}
	defer func() { _ = prioRows.Close() }()
	priorities := []PriorityHistoryEntry{}
	for prioRows.Next() {
		var e PriorityHistoryEntry
		var oldP, newP string
		if err := prioRows.Scan(&e.ID, &oldP, &newP,
			&e.ActorEmployeeID, &e.ActorDisplayName, &e.ChangedAt); err != nil {
			return nil, wrapRead(err, "scan priority history row")
		}
		e.OldPriority = model.IncidentPriority(oldP)
		e.NewPriority = model.IncidentPriority(newP)
		priorities = append(priorities, e)
	}
	if err := prioRows.Err(); err != nil {
		return nil, wrapRead(err, "iterate priority history")
	}
	return &IncidentHistory{Status: statuses, Priority: priorities}, nil
}
