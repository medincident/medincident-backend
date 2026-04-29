// Package analytics is the query-side analytics reader. It computes
// incident, service-request, and patient-buffer statistics directly
// from domain.* tables (not projections) because projections omit
// analytics-critical fields such as priority, source, and reopened_from.
//
// See: docs/services/analytics/Analytics.md
package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

const (
	ErrCodeAnalyticsOrgNotFound    = "analytics_org_not_found"
	ErrCodeAnalyticsClinicNotFound = "analytics_clinic_not_found"
	ErrCodeAnalyticsDeptNotFound   = "analytics_dept_not_found"
	ErrCodeAnalyticsPeriodInvalid  = "analytics_period_invalid"
	ErrCodeAnalyticsPeriodTooLarge = "analytics_period_too_large"
	ErrCodeAnalyticsQueryFailed    = "analytics_query_failed"
)

const (
	analyticsMaxSnapshotPeriodDays = 366
	analyticsTopNDistributions     = 10
	analyticsResolutionP50         = 0.50
	analyticsResolutionP90         = 0.90
	analyticsResolutionP95         = 0.95
)

// Reader exposes the analytics query surface.
type Reader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given gorm DB and authorization service.
func NewReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, authz: az, logger: logger}
}

// filter holds the parsed, validated scope parameters shared by all three methods.
type filter struct {
	orgID    uuid.UUID
	from     time.Time
	to       time.Time
	clinicID *uuid.UUID
	deptID   *uuid.UUID
}

// parseFilter validates and parses the common filter fields.
func parseFilter(orgIDStr, fromStr, toStr string, clinicIDStr, deptIDStr *string) (*filter, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil || orgID == uuid.Nil {
		return nil, oops.In("analytics").
			Code(ErrCodeAnalyticsOrgNotFound).
			Public("organization_id is not a valid UUID.").
			Wrap(fmt.Errorf("invalid organization_id: %s", orgIDStr))
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return nil, oops.In("analytics").
			Code(ErrCodeAnalyticsPeriodInvalid).
			Public("from is not a valid RFC3339 timestamp.").
			Wrap(err)
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return nil, oops.In("analytics").
			Code(ErrCodeAnalyticsPeriodInvalid).
			Public("to is not a valid RFC3339 timestamp.").
			Wrap(err)
	}
	if !from.Before(to) {
		return nil, oops.In("analytics").
			Code(ErrCodeAnalyticsPeriodInvalid).
			Public("from must be before to.").
			Wrap(fmt.Errorf("from >= to"))
	}

	var clinicID *uuid.UUID
	if clinicIDStr != nil && *clinicIDStr != "" {
		id, err := uuid.Parse(*clinicIDStr)
		if err != nil {
			return nil, oops.In("analytics").
				Code(ErrCodeAnalyticsClinicNotFound).
				Public("clinic_id is not a valid UUID.").
				Wrap(err)
		}
		clinicID = &id
	}
	var deptID *uuid.UUID
	if deptIDStr != nil && *deptIDStr != "" {
		id, err := uuid.Parse(*deptIDStr)
		if err != nil {
			return nil, oops.In("analytics").
				Code(ErrCodeAnalyticsDeptNotFound).
				Public("department_id is not a valid UUID.").
				Wrap(err)
		}
		deptID = &id
	}

	return &filter{orgID: orgID, from: from, to: to, clinicID: clinicID, deptID: deptID}, nil
}

// requireAuth checks that the caller holds a management role appropriate
// for the requested scope.
func (r *Reader) requireAuth(ctx context.Context, caller authz.Caller, f *filter) error {
	var policy authz.Policy
	switch {
	case f.deptID != nil:
		policies := []authz.Policy{
			authz.AdminOf.Organization(f.orgID),
			authz.OrgHeadOf.Organization(f.orgID),
			authz.OrgDispatcherOf.Organization(f.orgID),
			authz.DeptResponsibleOf.Department(*f.deptID),
		}
		if f.clinicID != nil {
			policies = append(policies, authz.ClinicHeadOf.Clinic(*f.clinicID))
		}
		policy = authz.AnyOf(policies...)
	case f.clinicID != nil:
		policy = authz.AnyOf(
			authz.AdminOf.Organization(f.orgID),
			authz.OrgHeadOf.Organization(f.orgID),
			authz.OrgDispatcherOf.Organization(f.orgID),
			authz.ClinicHeadOf.Clinic(*f.clinicID),
		)
	default:
		policy = authz.AnyOf(
			authz.AdminOf.Organization(f.orgID),
			authz.OrgHeadOf.Organization(f.orgID),
			authz.OrgDispatcherOf.Organization(f.orgID),
		)
	}
	return r.authz.Require(ctx, caller.ZitadelUserID, policy)
}

// validateScope verifies that clinic_id and department_id belong to the
// stated organization_id. Runs after auth — prevents cross-tenant leakage.
func (r *Reader) validateScope(ctx context.Context, f *filter) error {
	if f.clinicID != nil {
		var exists bool
		row := r.db.WithContext(ctx).Raw(
			`SELECT EXISTS(SELECT 1 FROM domain.clinics WHERE id = ? AND organization_id = ?)`,
			*f.clinicID, f.orgID,
		).Row()
		if err := row.Scan(&exists); err != nil {
			return oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		if !exists {
			return oops.In("analytics").
				Code(ErrCodeAnalyticsClinicNotFound).
				Public("clinic not found in this organization.").
				With("clinic_id", f.clinicID).
				With("organization_id", f.orgID).
				Wrap(fmt.Errorf("clinic not in org"))
		}
	}
	if f.deptID != nil {
		var exists bool
		row := r.db.WithContext(ctx).Raw(
			`SELECT EXISTS(
				SELECT 1 FROM domain.departments d
				JOIN domain.clinics c ON c.id = d.clinic_id
				WHERE d.id = ? AND c.organization_id = ?
			)`,
			*f.deptID, f.orgID,
		).Row()
		if err := row.Scan(&exists); err != nil {
			return oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		if !exists {
			return oops.In("analytics").
				Code(ErrCodeAnalyticsDeptNotFound).
				Public("department not found in this organization.").
				With("department_id", f.deptID).
				With("organization_id", f.orgID).
				Wrap(fmt.Errorf("department not in org"))
		}
	}
	return nil
}

// whereClause builds the WHERE fragment and args for the common filter.
// tableAlias is prepended to each column (e.g. "i" → "i.organization_id = ?").
// Pass "" for single-table queries that need no prefix.
func whereClause(f *filter, tableAlias string) (clause string, filterArgs []any) {
	prefix := ""
	if tableAlias != "" {
		prefix = tableAlias + "."
	}
	clauses := []string{
		prefix + "organization_id = ?",
		prefix + "created_at >= ?",
		prefix + "created_at <= ?",
	}
	filterArgs = []any{f.orgID, f.from, f.to}
	if f.clinicID != nil {
		clauses = append(clauses, prefix+"clinic_id = ?")
		filterArgs = append(filterArgs, *f.clinicID)
	}
	if f.deptID != nil {
		clauses = append(clauses, prefix+"department_id = ?")
		filterArgs = append(filterArgs, *f.deptID)
	}
	return strings.Join(clauses, " AND "), filterArgs
}

// ── Snapshot types ────────────────────────────────────────────────────────────

// SnapshotIncident is a minimal analytics record for one incident.
type SnapshotIncident struct {
	CreatedAt           time.Time
	OccurredAt          time.Time
	ClosedAt            *time.Time
	Status              string
	Priority            string
	CategoryID          uuid.UUID
	CategoryName        string
	TypeID              uuid.UUID
	TypeName            string
	ClinicID            uuid.UUID
	DepartmentID        uuid.UUID
	IsPatientSource     bool
	IsReopened          bool
	LinkedRequestsCount int32
}

// SnapshotRequest is a minimal analytics record for one service request.
type SnapshotRequest struct {
	CreatedAt         time.Time
	CompletedAt       *time.Time
	Status            string
	TypeID            uuid.UUID
	TypeName          string
	DepartmentID      uuid.UUID
	HasLinkedIncident bool
}

// SnapshotPatientBuffer is a minimal analytics record for one buffer entry.
type SnapshotPatientBuffer struct {
	CreatedAt  time.Time
	Status     string
	CategoryID *uuid.UUID
}

// SnapshotResult is the aggregate result for GetSnapshot.
type SnapshotResult struct {
	Incidents     []SnapshotIncident
	Requests      []SnapshotRequest
	PatientBuffer []SnapshotPatientBuffer
}

// GetSnapshot returns raw analytics records for the given filter.
//
// See: docs/services/analytics/Analytics.md
func (r *Reader) GetSnapshot(
	ctx context.Context,
	caller authz.Caller,
	orgIDStr, fromStr, toStr string,
	clinicIDStr, deptIDStr *string,
	includePatientBuffer bool,
) (*SnapshotResult, error) {
	f, err := parseFilter(orgIDStr, fromStr, toStr, clinicIDStr, deptIDStr)
	if err != nil {
		return nil, err
	}
	periodDays := f.to.Sub(f.from).Hours() / 24
	if periodDays > analyticsMaxSnapshotPeriodDays {
		return nil, oops.In("analytics").
			Code(ErrCodeAnalyticsPeriodTooLarge).
			Public(fmt.Sprintf("snapshot period must not exceed %d days.", analyticsMaxSnapshotPeriodDays)).
			With("days", periodDays).
			Wrap(fmt.Errorf("period too large"))
	}
	if err := r.requireAuth(ctx, caller, f); err != nil {
		return nil, err
	}
	if err := r.validateScope(ctx, f); err != nil {
		return nil, err
	}

	result := &SnapshotResult{}
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		incidents, err := r.snapshotIncidents(egCtx, f)
		if err != nil {
			return err
		}
		result.Incidents = incidents
		return nil
	})
	eg.Go(func() error {
		reqs, err := r.snapshotRequests(egCtx, f)
		if err != nil {
			return err
		}
		result.Requests = reqs
		return nil
	})
	if includePatientBuffer {
		eg.Go(func() error {
			buf, err := r.snapshotBuffer(egCtx, f)
			if err != nil {
				return err
			}
			result.PatientBuffer = buf
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Reader) snapshotIncidents(ctx context.Context, f *filter) ([]SnapshotIncident, error) {
	where, args := whereClause(f, "i")
	query := fmt.Sprintf(`
		SELECT
		  i.created_at,
		  i.occurred_at,
		  CASE WHEN i.status IN ('done','rejected','cancelled') THEN i.updated_at END AS closed_at,
		  i.status,
		  i.priority,
		  i.category_id,
		  ic.name AS category_name,
		  i.type_id,
		  it.name AS type_name,
		  i.clinic_id,
		  i.department_id,
		  (i.source_patient_zitadel_user_id IS NOT NULL) AS is_patient_source,
		  (i.reopened_from_incident_id IS NOT NULL)      AS is_reopened,
		  (SELECT COUNT(*)::int FROM domain.service_requests sr WHERE sr.incident_id = i.id) AS linked_requests_count
		FROM domain.incidents i
		JOIN domain.incident_categories ic ON ic.id = i.category_id
		JOIN domain.incident_types it ON it.id = i.type_id
		WHERE %s
		ORDER BY i.created_at`, where)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	defer rows.Close()

	var out []SnapshotIncident
	for rows.Next() {
		var si SnapshotIncident
		if err := rows.Scan(
			&si.CreatedAt, &si.OccurredAt, &si.ClosedAt,
			&si.Status, &si.Priority,
			&si.CategoryID, &si.CategoryName,
			&si.TypeID, &si.TypeName,
			&si.ClinicID, &si.DepartmentID,
			&si.IsPatientSource, &si.IsReopened, &si.LinkedRequestsCount,
		); err != nil {
			return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		out = append(out, si)
	}
	return out, rows.Err()
}

func (r *Reader) snapshotRequests(ctx context.Context, f *filter) ([]SnapshotRequest, error) {
	where, args := whereClause(f, "r")
	query := fmt.Sprintf(`
		SELECT
		  r.created_at,
		  CASE WHEN r.status = 'completed' THEN r.updated_at END AS completed_at,
		  r.status,
		  r.type_id,
		  rt.name AS type_name,
		  r.department_id,
		  (r.incident_id IS NOT NULL) AS has_linked_incident
		FROM domain.service_requests r
		JOIN domain.request_types rt ON rt.id = r.type_id
		WHERE %s
		ORDER BY r.created_at`, where)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	defer rows.Close()

	var out []SnapshotRequest
	for rows.Next() {
		var sr SnapshotRequest
		if err := rows.Scan(
			&sr.CreatedAt, &sr.CompletedAt,
			&sr.Status, &sr.TypeID, &sr.TypeName,
			&sr.DepartmentID, &sr.HasLinkedIncident,
		); err != nil {
			return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}

func (r *Reader) snapshotBuffer(ctx context.Context, f *filter) ([]SnapshotPatientBuffer, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT created_at, status, category_id
		FROM domain.patient_incident_buffer
		WHERE organization_id = ? AND created_at >= ? AND created_at <= ?
		ORDER BY created_at`,
		f.orgID, f.from, f.to,
	).Rows()
	if err != nil {
		return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	defer rows.Close()

	var out []SnapshotPatientBuffer
	for rows.Next() {
		var sb SnapshotPatientBuffer
		if err := rows.Scan(&sb.CreatedAt, &sb.Status, &sb.CategoryID); err != nil {
			return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		out = append(out, sb)
	}
	return out, rows.Err()
}

// ── Summary types ─────────────────────────────────────────────────────────────

// IncidentAggRow holds the flat aggregated incident counts from a single SQL pass.
type IncidentAggRow struct {
	Total              int64
	Pending            int64
	InProgress         int64
	Done               int64
	Rejected           int64
	Cancelled          int64
	Low                int64
	Normal             int64
	High               int64
	Critical           int64
	PatientSource      int64
	StaffSource        int64
	Reopened           int64
	WithLinkedRequests int64
	ResolutionAvg      sql.NullFloat64
	ResolutionMin      sql.NullFloat64
	ResolutionMax      sql.NullFloat64
}

// IncidentPercentiles holds P50/P90/P95 resolution minutes.
type IncidentPercentiles struct {
	P50 sql.NullFloat64
	P90 sql.NullFloat64
	P95 sql.NullFloat64
}

// DistributionRow is one row from a top-N distribution query.
type DistributionRow struct {
	ID    uuid.UUID
	Name  string
	Count int64
}

// RequestAggRow holds the flat aggregated request counts.
type RequestAggRow struct {
	Total         int64
	Created       int64
	InWork        int64
	OnHold        int64
	PendingReview int64
	Completed     int64
	Cancelled     int64
	Linked        int64
	Unlinked      int64
	CompletionAvg sql.NullFloat64
	CompletionMin sql.NullFloat64
	CompletionMax sql.NullFloat64
}

// RequestPercentiles holds P50/P90/P95 completion minutes.
type RequestPercentiles struct {
	P50 sql.NullFloat64
	P90 sql.NullFloat64
	P95 sql.NullFloat64
}

// BufferAggRow holds the flat aggregated buffer counts.
type BufferAggRow struct {
	Total     int64
	Pending   int64
	Published int64
	Rejected  int64
	Cancelled int64
}

// SummaryResult is the aggregate result for GetSummary.
type SummaryResult struct {
	IncidentAgg         IncidentAggRow
	IncidentPercentiles IncidentPercentiles
	IncidentCategories  []DistributionRow
	IncidentTypes       []DistributionRow
	IncidentDepartments []DistributionRow
	RequestAgg          RequestAggRow
	RequestPercentiles  RequestPercentiles
	RequestTypes        []DistributionRow
	RequestDepartments  []DistributionRow
	BufferAgg           BufferAggRow
}

// GetSummary returns aggregated KPIs and distributions.
//
// See: docs/services/analytics/Analytics.md
func (r *Reader) GetSummary(
	ctx context.Context,
	caller authz.Caller,
	orgIDStr, fromStr, toStr string,
	clinicIDStr, deptIDStr *string,
) (*SummaryResult, error) {
	f, err := parseFilter(orgIDStr, fromStr, toStr, clinicIDStr, deptIDStr)
	if err != nil {
		return nil, err
	}
	if err := r.requireAuth(ctx, caller, f); err != nil {
		return nil, err
	}
	if err := r.validateScope(ctx, f); err != nil {
		return nil, err
	}

	result := &SummaryResult{}
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		row, err := r.incidentAgg(egCtx, f)
		if err != nil {
			return err
		}
		result.IncidentAgg = row
		return nil
	})
	eg.Go(func() error {
		p, err := r.incidentPercentiles(egCtx, f)
		if err != nil {
			return err
		}
		result.IncidentPercentiles = p
		return nil
	})
	eg.Go(func() error {
		rows, err := r.topIncidentCategories(egCtx, f)
		if err != nil {
			return err
		}
		result.IncidentCategories = rows
		return nil
	})
	eg.Go(func() error {
		rows, err := r.topIncidentTypes(egCtx, f)
		if err != nil {
			return err
		}
		result.IncidentTypes = rows
		return nil
	})
	eg.Go(func() error {
		rows, err := r.topIncidentDepartments(egCtx, f)
		if err != nil {
			return err
		}
		result.IncidentDepartments = rows
		return nil
	})
	eg.Go(func() error {
		row, err := r.requestAgg(egCtx, f)
		if err != nil {
			return err
		}
		result.RequestAgg = row
		return nil
	})
	eg.Go(func() error {
		p, err := r.requestPercentiles(egCtx, f)
		if err != nil {
			return err
		}
		result.RequestPercentiles = p
		return nil
	})
	eg.Go(func() error {
		rows, err := r.topRequestTypes(egCtx, f)
		if err != nil {
			return err
		}
		result.RequestTypes = rows
		return nil
	})
	eg.Go(func() error {
		rows, err := r.topRequestDepartments(egCtx, f)
		if err != nil {
			return err
		}
		result.RequestDepartments = rows
		return nil
	})
	eg.Go(func() error {
		row, err := r.bufferAgg(egCtx, f)
		if err != nil {
			return err
		}
		result.BufferAgg = row
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Reader) incidentAgg(ctx context.Context, f *filter) (IncidentAggRow, error) {
	where, args := whereClause(f, "i")
	query := fmt.Sprintf(`
		SELECT
		  COUNT(*)                                                                 AS total,
		  COUNT(*) FILTER (WHERE i.status = 'pending')                            AS pending,
		  COUNT(*) FILTER (WHERE i.status = 'in_progress')                        AS in_progress,
		  COUNT(*) FILTER (WHERE i.status = 'done')                               AS done,
		  COUNT(*) FILTER (WHERE i.status = 'rejected')                           AS rejected,
		  COUNT(*) FILTER (WHERE i.status = 'cancelled')                          AS cancelled,
		  COUNT(*) FILTER (WHERE i.priority = 'low')                              AS low,
		  COUNT(*) FILTER (WHERE i.priority = 'normal')                           AS normal,
		  COUNT(*) FILTER (WHERE i.priority = 'high')                             AS high,
		  COUNT(*) FILTER (WHERE i.priority = 'critical')                         AS critical,
		  COUNT(*) FILTER (WHERE i.source_patient_zitadel_user_id IS NOT NULL)    AS patient_source,
		  COUNT(*) FILTER (WHERE i.source_patient_zitadel_user_id IS NULL)        AS staff_source,
		  COUNT(*) FILTER (WHERE i.reopened_from_incident_id IS NOT NULL)         AS reopened,
		  COUNT(*) FILTER (WHERE EXISTS (
		    SELECT 1 FROM domain.service_requests sr WHERE sr.incident_id = i.id
		  ))                                                                      AS with_linked_requests,
		  AVG(CASE WHEN i.status = 'done'
		    THEN EXTRACT(EPOCH FROM (i.updated_at - i.created_at))/60.0 END)     AS resolution_avg,
		  MIN(CASE WHEN i.status = 'done'
		    THEN EXTRACT(EPOCH FROM (i.updated_at - i.created_at))/60.0 END)     AS resolution_min,
		  MAX(CASE WHEN i.status = 'done'
		    THEN EXTRACT(EPOCH FROM (i.updated_at - i.created_at))/60.0 END)     AS resolution_max
		FROM domain.incidents i
		WHERE %s`, where)

	var row IncidentAggRow
	r2 := r.db.WithContext(ctx).Raw(query, args...).Row()
	if err := r2.Scan(
		&row.Total, &row.Pending, &row.InProgress, &row.Done, &row.Rejected, &row.Cancelled,
		&row.Low, &row.Normal, &row.High, &row.Critical,
		&row.PatientSource, &row.StaffSource, &row.Reopened, &row.WithLinkedRequests,
		&row.ResolutionAvg, &row.ResolutionMin, &row.ResolutionMax,
	); err != nil {
		return IncidentAggRow{}, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	return row, nil
}

func (r *Reader) incidentPercentiles(ctx context.Context, f *filter) (IncidentPercentiles, error) {
	where, args := whereClause(f, "")
	query := fmt.Sprintf(`
		SELECT
		  PERCENTILE_CONT(%f) WITHIN GROUP (ORDER BY mins) AS p50,
		  PERCENTILE_CONT(%f) WITHIN GROUP (ORDER BY mins) AS p90,
		  PERCENTILE_CONT(%f) WITHIN GROUP (ORDER BY mins) AS p95
		FROM (
		  SELECT EXTRACT(EPOCH FROM (updated_at - created_at))/60.0 AS mins
		  FROM domain.incidents
		  WHERE %s AND status = 'done'
		) t`,
		analyticsResolutionP50, analyticsResolutionP90, analyticsResolutionP95, where)

	var p IncidentPercentiles
	row := r.db.WithContext(ctx).Raw(query, args...).Row()
	if err := row.Scan(&p.P50, &p.P90, &p.P95); err != nil {
		return IncidentPercentiles{}, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	return p, nil
}

func (r *Reader) topIncidentCategories(ctx context.Context, f *filter) ([]DistributionRow, error) {
	where, args := whereClause(f, "i")
	args = append(args, analyticsTopNDistributions)
	query := fmt.Sprintf(`
		SELECT i.category_id, ic.name, COUNT(*) AS cnt
		FROM domain.incidents i
		JOIN domain.incident_categories ic ON ic.id = i.category_id
		WHERE %s
		GROUP BY i.category_id, ic.name
		ORDER BY cnt DESC LIMIT ?`, where)
	return r.scanDistribution(ctx, query, args)
}

func (r *Reader) topIncidentTypes(ctx context.Context, f *filter) ([]DistributionRow, error) {
	where, args := whereClause(f, "i")
	args = append(args, analyticsTopNDistributions)
	query := fmt.Sprintf(`
		SELECT i.type_id, it.name, COUNT(*) AS cnt
		FROM domain.incidents i
		JOIN domain.incident_types it ON it.id = i.type_id
		WHERE %s
		GROUP BY i.type_id, it.name
		ORDER BY cnt DESC LIMIT ?`, where)
	return r.scanDistribution(ctx, query, args)
}

func (r *Reader) topIncidentDepartments(ctx context.Context, f *filter) ([]DistributionRow, error) {
	where, args := whereClause(f, "i")
	args = append(args, analyticsTopNDistributions)
	query := fmt.Sprintf(`
		SELECT i.department_id, d.name, COUNT(*) AS cnt
		FROM domain.incidents i
		JOIN domain.departments d ON d.id = i.department_id
		WHERE %s
		GROUP BY i.department_id, d.name
		ORDER BY cnt DESC LIMIT ?`, where)
	return r.scanDistribution(ctx, query, args)
}

func (r *Reader) requestAgg(ctx context.Context, f *filter) (RequestAggRow, error) {
	where, args := whereClause(f, "")
	query := fmt.Sprintf(`
		SELECT
		  COUNT(*)                                                          AS total,
		  COUNT(*) FILTER (WHERE status = 'created')                       AS created,
		  COUNT(*) FILTER (WHERE status = 'in_work')                       AS in_work,
		  COUNT(*) FILTER (WHERE status = 'on_hold')                       AS on_hold,
		  COUNT(*) FILTER (WHERE status = 'pending_review')                AS pending_review,
		  COUNT(*) FILTER (WHERE status = 'completed')                     AS completed,
		  COUNT(*) FILTER (WHERE status = 'cancelled')                     AS cancelled,
		  COUNT(*) FILTER (WHERE incident_id IS NOT NULL)                  AS linked,
		  COUNT(*) FILTER (WHERE incident_id IS NULL)                      AS unlinked,
		  AVG(CASE WHEN status = 'completed'
		    THEN EXTRACT(EPOCH FROM (updated_at - created_at))/60.0 END)  AS completion_avg,
		  MIN(CASE WHEN status = 'completed'
		    THEN EXTRACT(EPOCH FROM (updated_at - created_at))/60.0 END)  AS completion_min,
		  MAX(CASE WHEN status = 'completed'
		    THEN EXTRACT(EPOCH FROM (updated_at - created_at))/60.0 END)  AS completion_max
		FROM domain.service_requests
		WHERE %s`, where)

	var row RequestAggRow
	r2 := r.db.WithContext(ctx).Raw(query, args...).Row()
	if err := r2.Scan(
		&row.Total, &row.Created, &row.InWork, &row.OnHold, &row.PendingReview,
		&row.Completed, &row.Cancelled, &row.Linked, &row.Unlinked,
		&row.CompletionAvg, &row.CompletionMin, &row.CompletionMax,
	); err != nil {
		return RequestAggRow{}, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	return row, nil
}

func (r *Reader) requestPercentiles(ctx context.Context, f *filter) (RequestPercentiles, error) {
	where, args := whereClause(f, "")
	query := fmt.Sprintf(`
		SELECT
		  PERCENTILE_CONT(%f) WITHIN GROUP (ORDER BY mins) AS p50,
		  PERCENTILE_CONT(%f) WITHIN GROUP (ORDER BY mins) AS p90,
		  PERCENTILE_CONT(%f) WITHIN GROUP (ORDER BY mins) AS p95
		FROM (
		  SELECT EXTRACT(EPOCH FROM (updated_at - created_at))/60.0 AS mins
		  FROM domain.service_requests
		  WHERE %s AND status = 'completed'
		) t`,
		analyticsResolutionP50, analyticsResolutionP90, analyticsResolutionP95, where)

	var p RequestPercentiles
	row := r.db.WithContext(ctx).Raw(query, args...).Row()
	if err := row.Scan(&p.P50, &p.P90, &p.P95); err != nil {
		return RequestPercentiles{}, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	return p, nil
}

func (r *Reader) topRequestTypes(ctx context.Context, f *filter) ([]DistributionRow, error) {
	where, args := whereClause(f, "r")
	args = append(args, analyticsTopNDistributions)
	query := fmt.Sprintf(`
		SELECT r.type_id, rt.name, COUNT(*) AS cnt
		FROM domain.service_requests r
		JOIN domain.request_types rt ON rt.id = r.type_id
		WHERE %s
		GROUP BY r.type_id, rt.name
		ORDER BY cnt DESC LIMIT ?`, where)
	return r.scanDistribution(ctx, query, args)
}

func (r *Reader) topRequestDepartments(ctx context.Context, f *filter) ([]DistributionRow, error) {
	where, args := whereClause(f, "r")
	args = append(args, analyticsTopNDistributions)
	query := fmt.Sprintf(`
		SELECT r.department_id, d.name, COUNT(*) AS cnt
		FROM domain.service_requests r
		JOIN domain.departments d ON d.id = r.department_id
		WHERE %s
		GROUP BY r.department_id, d.name
		ORDER BY cnt DESC LIMIT ?`, where)
	return r.scanDistribution(ctx, query, args)
}

func (r *Reader) bufferAgg(ctx context.Context, f *filter) (BufferAggRow, error) {
	row := r.db.WithContext(ctx).Raw(`
		SELECT
		  COUNT(*)                                           AS total,
		  COUNT(*) FILTER (WHERE status = 'pending')        AS pending,
		  COUNT(*) FILTER (WHERE status = 'published')      AS published,
		  COUNT(*) FILTER (WHERE status = 'rejected')       AS rejected,
		  COUNT(*) FILTER (WHERE status = 'cancelled')      AS cancelled
		FROM domain.patient_incident_buffer
		WHERE organization_id = ? AND created_at >= ? AND created_at <= ?`,
		f.orgID, f.from, f.to,
	).Row()

	var a BufferAggRow
	if err := row.Scan(&a.Total, &a.Pending, &a.Published, &a.Rejected, &a.Cancelled); err != nil {
		return BufferAggRow{}, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	return a, nil
}

func (r *Reader) scanDistribution(ctx context.Context, query string, args []any) ([]DistributionRow, error) {
	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
	}
	defer rows.Close()

	var out []DistributionRow
	for rows.Next() {
		var d DistributionRow
		if err := rows.Scan(&d.ID, &d.Name, &d.Count); err != nil {
			return nil, oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// ── TimeSeries types ──────────────────────────────────────────────────────────

// TimeSeriesGranularity is the bucket width for GetTimeSeries.
type TimeSeriesGranularity int

const (
	GranularityDay TimeSeriesGranularity = iota + 1
	GranularityWeek
	GranularityMonth
)

// TimeSeriesBucket holds counts for a single time bucket.
type TimeSeriesBucket struct {
	BucketStart    time.Time
	BucketEnd      time.Time
	IncidentTotal  int64
	IPending       int64
	IInProgress    int64
	IDone          int64
	IRejected      int64
	ICancelled     int64
	IHighCritical  int64
	IPatientSource int64
	IReopened      int64
	ReqTotal       int64
	ReqCompleted   int64
	ReqCancelled   int64
	ReqLinked      int64
}

// GetTimeSeries returns bucketed counts for the given filter and granularity.
//
// See: docs/services/analytics/Analytics.md
func (r *Reader) GetTimeSeries(
	ctx context.Context,
	caller authz.Caller,
	orgIDStr, fromStr, toStr string,
	clinicIDStr, deptIDStr *string,
	gran TimeSeriesGranularity,
) ([]TimeSeriesBucket, error) {
	f, err := parseFilter(orgIDStr, fromStr, toStr, clinicIDStr, deptIDStr)
	if err != nil {
		return nil, err
	}
	if err := r.requireAuth(ctx, caller, f); err != nil {
		return nil, err
	}
	if err := r.validateScope(ctx, f); err != nil {
		return nil, err
	}

	interval := granularityInterval(gran)

	type incidentBucketRow struct {
		BucketStart time.Time
		BucketEnd   time.Time
		Total       int64
		Pending     int64
		InProgress  int64
		Done        int64
		Rejected    int64
		Cancelled   int64
		HighCrit    int64
		Patient     int64
		Reopened    int64
	}
	type requestBucketRow struct {
		BucketStart time.Time
		BucketEnd   time.Time
		Total       int64
		Completed   int64
		Cancelled   int64
		Linked      int64
	}

	var iRows []incidentBucketRow
	var rRows []requestBucketRow

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		where, args := whereClause(f, "i")
		query := fmt.Sprintf(`
			WITH buckets AS (
			  SELECT gs AS bucket_start, gs + INTERVAL '%s' AS bucket_end
			  FROM generate_series(?::timestamptz, ?::timestamptz, INTERVAL '%s') gs
			)
			SELECT
			  b.bucket_start, b.bucket_end,
			  COUNT(i.id)                                                              AS total,
			  COUNT(i.id) FILTER (WHERE i.status = 'pending')                         AS pending,
			  COUNT(i.id) FILTER (WHERE i.status = 'in_progress')                     AS in_progress,
			  COUNT(i.id) FILTER (WHERE i.status = 'done')                            AS done,
			  COUNT(i.id) FILTER (WHERE i.status = 'rejected')                        AS rejected,
			  COUNT(i.id) FILTER (WHERE i.status = 'cancelled')                       AS cancelled,
			  COUNT(i.id) FILTER (WHERE i.priority IN ('high','critical'))            AS high_critical,
			  COUNT(i.id) FILTER (WHERE i.source_patient_zitadel_user_id IS NOT NULL) AS patient_source,
			  COUNT(i.id) FILTER (WHERE i.reopened_from_incident_id IS NOT NULL)      AS reopened
			FROM buckets b
			LEFT JOIN domain.incidents i
			  ON i.created_at >= b.bucket_start
			 AND i.created_at <  b.bucket_end
			 AND %s
			GROUP BY b.bucket_start, b.bucket_end
			ORDER BY b.bucket_start`, interval, interval, where)

		tsArgs := append([]any{f.from, f.to}, args...)
		rows, err := r.db.WithContext(egCtx).Raw(query, tsArgs...).Rows()
		if err != nil {
			return oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		defer rows.Close()
		for rows.Next() {
			var row incidentBucketRow
			if err := rows.Scan(
				&row.BucketStart, &row.BucketEnd,
				&row.Total, &row.Pending, &row.InProgress, &row.Done,
				&row.Rejected, &row.Cancelled, &row.HighCrit, &row.Patient, &row.Reopened,
			); err != nil {
				return oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
			}
			iRows = append(iRows, row)
		}
		return rows.Err()
	})

	eg.Go(func() error {
		where, args := whereClause(f, "r")
		query := fmt.Sprintf(`
			WITH buckets AS (
			  SELECT gs AS bucket_start, gs + INTERVAL '%s' AS bucket_end
			  FROM generate_series(?::timestamptz, ?::timestamptz, INTERVAL '%s') gs
			)
			SELECT
			  b.bucket_start, b.bucket_end,
			  COUNT(r.id)                                             AS total,
			  COUNT(r.id) FILTER (WHERE r.status = 'completed')     AS completed,
			  COUNT(r.id) FILTER (WHERE r.status = 'cancelled')     AS cancelled,
			  COUNT(r.id) FILTER (WHERE r.incident_id IS NOT NULL)  AS linked
			FROM buckets b
			LEFT JOIN domain.service_requests r
			  ON r.created_at >= b.bucket_start
			 AND r.created_at <  b.bucket_end
			 AND %s
			GROUP BY b.bucket_start, b.bucket_end
			ORDER BY b.bucket_start`, interval, interval, where)

		tsArgs := append([]any{f.from, f.to}, args...)
		rows, err := r.db.WithContext(egCtx).Raw(query, tsArgs...).Rows()
		if err != nil {
			return oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
		}
		defer rows.Close()
		for rows.Next() {
			var row requestBucketRow
			if err := rows.Scan(
				&row.BucketStart, &row.BucketEnd,
				&row.Total, &row.Completed, &row.Cancelled, &row.Linked,
			); err != nil {
				return oops.In("analytics").Code(ErrCodeAnalyticsQueryFailed).Wrap(err)
			}
			rRows = append(rRows, row)
		}
		return rows.Err()
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	buckets := make([]TimeSeriesBucket, len(iRows))
	for i, ir := range iRows {
		buckets[i] = TimeSeriesBucket{
			BucketStart:    ir.BucketStart,
			BucketEnd:      ir.BucketEnd,
			IncidentTotal:  ir.Total,
			IPending:       ir.Pending,
			IInProgress:    ir.InProgress,
			IDone:          ir.Done,
			IRejected:      ir.Rejected,
			ICancelled:     ir.Cancelled,
			IHighCritical:  ir.HighCrit,
			IPatientSource: ir.Patient,
			IReopened:      ir.Reopened,
		}
		if i < len(rRows) {
			rr := rRows[i]
			buckets[i].ReqTotal = rr.Total
			buckets[i].ReqCompleted = rr.Completed
			buckets[i].ReqCancelled = rr.Cancelled
			buckets[i].ReqLinked = rr.Linked
		}
	}
	return buckets, nil
}

func granularityInterval(g TimeSeriesGranularity) string {
	switch g {
	case GranularityWeek:
		return "1 week"
	case GranularityMonth:
		return "1 month"
	default:
		return "1 day"
	}
}
