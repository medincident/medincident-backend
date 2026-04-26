// Package stats exposes aggregate statistics against the rich
// projections. Counter columns are pulled from the matching
// projections.*_counters row; vacation counts are computed at read
// time from projections.employee_vacations so an "employees on
// vacation now" number never drifts relative to the source of truth.
//
// Authorization model: every Get*Stats method is gated by
// authz.ReaderOf.{Organization,Clinic,Department}. Cross-org reads
// fail permission_denied with no distinction from "scope missing",
// matching the rest of the query surface.
package stats

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// ErrCodeStatsLoadFailed is returned when the underlying DB read for
// any of the Get*Stats methods fails.
const ErrCodeStatsLoadFailed = "stats_load_failed"

// OrganizationStats is the aggregate snapshot the dashboards render
// for a single organization. Unknown orgs yield a zero-filled struct
// and nil error — dashboards show "0", not 404.
type OrganizationStats struct {
	OrganizationID      uuid.UUID
	EmployeesTotal      int64
	ClinicsTotal        int64
	DepartmentsTotal    int64
	EmployeesOnVacation int64
	VacationsScheduled  int64
}

// ClinicStats is the aggregate snapshot for a single clinic.
type ClinicStats struct {
	ClinicID            uuid.UUID
	OrganizationID      uuid.UUID
	EmployeesTotal      int64
	DepartmentsTotal    int64
	EmployeesOnVacation int64
}

// DepartmentStats is the aggregate snapshot for a single department.
type DepartmentStats struct {
	DepartmentID        uuid.UUID
	ClinicID            *uuid.UUID
	OrganizationID      uuid.UUID
	EmployeesTotal      int64
	EmployeesOnVacation int64
}

// Reader exposes the stats aggregate query surface.
type Reader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewReader returns a Reader bound to the given gorm DB and
// authorization service.
func NewReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *Reader {
	return &Reader{db: db, authz: az, logger: logger}
}

// GetOrganizationStats returns the aggregate snapshot for an
// organization. A missing counter row is soft-missed and returns a
// zero-filled struct rather than an error. Authorization:
// authz.ReaderOf.Organization(orgID).
//
// See: docs/services/Stats.md
func (r *Reader) GetOrganizationStats(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
) (*OrganizationStats, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}
	const query = `
		SELECT COALESCE(oc.employees_total, 0)   AS employees_total,
		       COALESCE(oc.clinics_total, 0)     AS clinics_total,
		       COALESCE(oc.departments_total, 0) AS departments_total,
		       (SELECT COUNT(*)::bigint
		          FROM projections.employee_vacations v
		          JOIN projections.employees e ON e.id = v.employee_id
		         WHERE e.organization_id = ?
		           AND e.terminated_at IS NULL
		           AND v.state = 'active')       AS employees_on_vacation,
		       (SELECT COUNT(*)::bigint
		          FROM projections.employee_vacations v
		          JOIN projections.employees e ON e.id = v.employee_id
		         WHERE e.organization_id = ?
		           AND e.terminated_at IS NULL
		           AND v.state = 'scheduled')    AS vacations_scheduled
		  FROM (SELECT ?::uuid AS organization_id) filter
		  LEFT JOIN projections.organization_counters oc
		    ON oc.organization_id = filter.organization_id
	`
	out := &OrganizationStats{OrganizationID: orgID}
	row := r.db.WithContext(ctx).Raw(query, orgID, orgID, orgID).Row()
	if err := row.Scan(
		&out.EmployeesTotal,
		&out.ClinicsTotal,
		&out.DepartmentsTotal,
		&out.EmployeesOnVacation,
		&out.VacationsScheduled,
	); err != nil {
		return nil, oops.In("reader.stats").
			Code(ErrCodeStatsLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	return out, nil
}

// GetClinicStats returns the aggregate snapshot for a clinic. A
// missing counter row yields zero-filled counts. Authorization:
// authz.ReaderOf.Clinic(clinicID).
//
// See: docs/services/Stats.md
func (r *Reader) GetClinicStats(
	ctx context.Context,
	caller authz.Caller,
	clinicID uuid.UUID,
) (*ClinicStats, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Clinic(clinicID)); err != nil {
		return nil, err
	}
	const query = `
		SELECT COALESCE(cc.organization_id, '00000000-0000-0000-0000-000000000000'::uuid)
		                                          AS organization_id,
		       COALESCE(cc.employees_total, 0)    AS employees_total,
		       COALESCE(cc.departments_total, 0)  AS departments_total,
		       (SELECT COUNT(*)::bigint
		          FROM projections.employee_vacations v
		          JOIN projections.employees e ON e.id = v.employee_id
		         WHERE e.clinic_id = ?
		           AND e.terminated_at IS NULL
		           AND v.state = 'active')        AS employees_on_vacation
		  FROM (SELECT ?::uuid AS clinic_id) filter
		  LEFT JOIN projections.clinic_counters cc
		    ON cc.clinic_id = filter.clinic_id
	`
	out := &ClinicStats{ClinicID: clinicID}
	row := r.db.WithContext(ctx).Raw(query, clinicID, clinicID).Row()
	if err := row.Scan(
		&out.OrganizationID,
		&out.EmployeesTotal,
		&out.DepartmentsTotal,
		&out.EmployeesOnVacation,
	); err != nil {
		return nil, oops.In("reader.stats").
			Code(ErrCodeStatsLoadFailed).
			With("clinic_id", clinicID).
			Wrap(err)
	}
	return out, nil
}

// GetDepartmentStats returns the aggregate snapshot for a department.
// A missing counter row yields zero-filled counts; clinic_id stays
// nil and organization_id stays zero in that soft-miss case.
// Authorization: authz.ReaderOf.Department(deptID).
//
// See: docs/services/Stats.md
func (r *Reader) GetDepartmentStats(
	ctx context.Context,
	caller authz.Caller,
	deptID uuid.UUID,
) (*DepartmentStats, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Department(deptID)); err != nil {
		return nil, err
	}
	const query = `
		SELECT dc.clinic_id                                                AS clinic_id,
		       COALESCE(dc.organization_id, '00000000-0000-0000-0000-000000000000'::uuid)
		                                                                   AS organization_id,
		       COALESCE(dc.employees_total, 0)                             AS employees_total,
		       (SELECT COUNT(*)::bigint
		          FROM projections.employee_vacations v
		          JOIN projections.employees e ON e.id = v.employee_id
		         WHERE e.department_id = ?
		           AND e.terminated_at IS NULL
		           AND v.state = 'active')                                 AS employees_on_vacation
		  FROM (SELECT ?::uuid AS department_id) filter
		  LEFT JOIN projections.department_counters dc
		    ON dc.department_id = filter.department_id
	`
	out := &DepartmentStats{DepartmentID: deptID}
	row := r.db.WithContext(ctx).Raw(query, deptID, deptID).Row()
	if err := row.Scan(
		&out.ClinicID,
		&out.OrganizationID,
		&out.EmployeesTotal,
		&out.EmployeesOnVacation,
	); err != nil {
		return nil, oops.In("reader.stats").
			Code(ErrCodeStatsLoadFailed).
			With("department_id", deptID).
			Wrap(err)
	}
	return out, nil
}
