package membership

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
)

// Cascade helpers for the five role types (DepartmentResponsible,
// ClinicHead, OrgAdmin, OrgHead, OrgDispatcher). Every role has the
// same shape of lifecycle cascade:
//
//  1. Revoke a role holder — optionally scoped by the role's "where"
//     column (department_id / clinic_id / organization_id). For each
//     affected row: if the deputy slot is set, publish DeputyRemoved
//     first (Rule 2: cleanup before terminate), then DELETE the row,
//     then publish Revoked.
//
//  2. Clear a deputy slot — for rows where the employee is the deputy
//     of someone else's role. For each affected row: NULL the deputy
//     slot, then publish DeputyRemoved.
//
// To keep the code surface small the two bodies above are written
// once as generic loops; a per-role cascadeSpec[R] carries the
// type-specific bits (table name, key columns, error codes, scope
// string, the publish helpers and the field accessors).

// cascadeSpec bundles the per-role knobs needed by the generic
// cascade loops. R is the GORM model type for the role table; K is
// the type of the role's scope column (always uuid.UUID today).
type cascadeSpec[R any, K comparable] struct {
	// SQL column name for the scope ID (e.g. "clinic_id").
	scopeCol string
	// Raw SQL table name used when NULLing the deputy slot. The DELETE
	// path uses gorm's model-based Delete, so only the clear path
	// needs a table name.
	table string

	// Error codes + scope string for oops.
	loadErrCode   string
	saveErrCode   string
	deleteErrCode string
	scope         string

	// Field accessors — closures so the model types stay method-free.
	scopeOf     func(*R) K
	empOf       func(*R) uuid.UUID
	deputyValid func(*R) bool

	// Event publishers for the role.
	publishRevoked       func(tx *gorm.DB, scopeID K, employeeID uuid.UUID, now time.Time) error
	publishDeputyRemoved func(tx *gorm.DB, scopeID K, employeeID uuid.UUID, now time.Time) error
}

// cascadeRevokeRolesBy revokes every role row matching the given
// WHERE clause. Order per row: DeputyRemoved (if slot valid) → DELETE →
// Revoked.
func cascadeRevokeRolesBy[R any, K comparable](
	tx *gorm.DB,
	spec cascadeSpec[R, K],
	where string,
	args []any,
	now time.Time,
) error {
	var rows []R
	if err := tx.Where(where, args...).Find(&rows).Error; err != nil {
		return oops.In(spec.scope).Code(spec.loadErrCode).Wrap(err)
	}
	deleteWhere := spec.scopeCol + " = ? AND employee_id = ?"
	for i := range rows {
		row := &rows[i]
		scopeID := spec.scopeOf(row)
		empID := spec.empOf(row)
		if spec.deputyValid(row) {
			if err := spec.publishDeputyRemoved(tx, scopeID, empID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(new(R), deleteWhere, scopeID, empID).Error; err != nil {
			return oops.In(spec.scope).Code(spec.deleteErrCode).Wrap(err)
		}
		if err := spec.publishRevoked(tx, scopeID, empID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearRoleDeputyBy clears the deputy slot on every role row
// matching `where`, then publishes DeputyRemoved for each.
func cascadeClearRoleDeputyBy[R any, K comparable](
	tx *gorm.DB,
	spec cascadeSpec[R, K],
	where string,
	args []any,
	now time.Time,
) error {
	var rows []R
	if err := tx.Where(where, args...).Find(&rows).Error; err != nil {
		return oops.In(spec.scope).Code(spec.loadErrCode).Wrap(err)
	}
	updateSQL := `UPDATE ` + spec.table +
		` SET deputy_employee_id = NULL, updated_at = now() WHERE ` +
		spec.scopeCol + ` = ? AND employee_id = ?`
	for i := range rows {
		row := &rows[i]
		scopeID := spec.scopeOf(row)
		empID := spec.empOf(row)
		if err := tx.Exec(updateSQL, scopeID, empID).Error; err != nil {
			return oops.In(spec.scope).Code(spec.saveErrCode).Wrap(err)
		}
		if err := spec.publishDeputyRemoved(tx, scopeID, empID, now); err != nil {
			return err
		}
	}
	return nil
}

// Per-role specs. Each is a package-level singleton — behaviour is
// static and the closures hold no state.

var drCascade = cascadeSpec[model.DepartmentResponsible, uuid.UUID]{
	scopeCol:             "department_id",
	table:                "domain.department_responsibles",
	loadErrCode:          ErrCodeDepartmentResponsibleLoadFailed,
	saveErrCode:          ErrCodeDepartmentResponsibleSaveFailed,
	deleteErrCode:        ErrCodeDepartmentResponsibleDeleteFailed,
	scope:                scopeDepartmentResponsible,
	scopeOf:              func(r *model.DepartmentResponsible) uuid.UUID { return r.DepartmentID },
	empOf:                func(r *model.DepartmentResponsible) uuid.UUID { return r.EmployeeID },
	deputyValid:          func(r *model.DepartmentResponsible) bool { return r.DeputyEmployeeID.Valid },
	publishRevoked:       publishDepartmentResponsibleRevoked,
	publishDeputyRemoved: publishDepartmentResponsibleDeputyRemoved,
}

var chCascade = cascadeSpec[model.ClinicHead, uuid.UUID]{
	scopeCol:             "clinic_id",
	table:                "domain.clinic_heads",
	loadErrCode:          ErrCodeClinicHeadLoadFailed,
	saveErrCode:          ErrCodeClinicHeadSaveFailed,
	deleteErrCode:        ErrCodeClinicHeadDeleteFailed,
	scope:                scopeClinicHead,
	scopeOf:              func(r *model.ClinicHead) uuid.UUID { return r.ClinicID },
	empOf:                func(r *model.ClinicHead) uuid.UUID { return r.EmployeeID },
	deputyValid:          func(r *model.ClinicHead) bool { return r.DeputyEmployeeID.Valid },
	publishRevoked:       publishClinicHeadRevoked,
	publishDeputyRemoved: publishClinicHeadDeputyRemoved,
}

var orgAdminCascade = cascadeSpec[model.OrgAdmin, uuid.UUID]{
	scopeCol:             "organization_id",
	table:                "domain.org_admins",
	loadErrCode:          ErrCodeOrganizationAdminLoadFailed,
	saveErrCode:          ErrCodeOrganizationAdminSaveFailed,
	deleteErrCode:        ErrCodeOrganizationAdminDeleteFailed,
	scope:                scopeOrgAdmin,
	scopeOf:              func(r *model.OrgAdmin) uuid.UUID { return r.OrganizationID },
	empOf:                func(r *model.OrgAdmin) uuid.UUID { return r.EmployeeID },
	deputyValid:          func(r *model.OrgAdmin) bool { return r.DeputyEmployeeID.Valid },
	publishRevoked:       publishOrgAdminRevoked,
	publishDeputyRemoved: publishOrgAdminDeputyRemoved,
}

var orgHeadCascade = cascadeSpec[model.OrgHead, uuid.UUID]{
	scopeCol:             "organization_id",
	table:                "domain.org_heads",
	loadErrCode:          ErrCodeOrganizationHeadLoadFailed,
	saveErrCode:          ErrCodeOrganizationHeadSaveFailed,
	deleteErrCode:        ErrCodeOrganizationHeadDeleteFailed,
	scope:                scopeOrgHead,
	scopeOf:              func(r *model.OrgHead) uuid.UUID { return r.OrganizationID },
	empOf:                func(r *model.OrgHead) uuid.UUID { return r.EmployeeID },
	deputyValid:          func(r *model.OrgHead) bool { return r.DeputyEmployeeID.Valid },
	publishRevoked:       publishOrgHeadRevoked,
	publishDeputyRemoved: publishOrgHeadDeputyRemoved,
}

var orgDispatcherCascade = cascadeSpec[model.OrgDispatcher, uuid.UUID]{
	scopeCol:             "organization_id",
	table:                "domain.org_dispatchers",
	loadErrCode:          ErrCodeOrganizationDispatcherLoadFailed,
	saveErrCode:          ErrCodeOrganizationDispatcherSaveFailed,
	deleteErrCode:        ErrCodeOrganizationDispatcherDeleteFailed,
	scope:                scopeOrgDispatcher,
	scopeOf:              func(r *model.OrgDispatcher) uuid.UUID { return r.OrganizationID },
	empOf:                func(r *model.OrgDispatcher) uuid.UUID { return r.EmployeeID },
	deputyValid:          func(r *model.OrgDispatcher) bool { return r.DeputyEmployeeID.Valid },
	publishRevoked:       publishOrgDispatcherRevoked,
	publishDeputyRemoved: publishOrgDispatcherDeputyRemoved,
}

// Call-site wrappers preserve the existing API so the rest of the
// service code keeps working without change.

func cascadeRevokeDepartmentResponsible(tx *gorm.DB, employeeID, departmentID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, drCascade, "employee_id = ? AND department_id = ?", []any{employeeID, departmentID}, now)
}

func cascadeRevokeDepartmentResponsibleAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, drCascade, "employee_id = ?", []any{employeeID}, now)
}

func cascadeClearDepartmentResponsibleDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, drCascade, "deputy_employee_id = ?", []any{employeeID}, now)
}

func cascadeClearDepartmentResponsibleDeputyInDepartment(tx *gorm.DB, employeeID, departmentID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, drCascade, "deputy_employee_id = ? AND department_id = ?", []any{employeeID, departmentID}, now)
}

func cascadeRevokeClinicHead(tx *gorm.DB, employeeID, clinicID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, chCascade, "employee_id = ? AND clinic_id = ?", []any{employeeID, clinicID}, now)
}

func cascadeRevokeClinicHeadAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, chCascade, "employee_id = ?", []any{employeeID}, now)
}

func cascadeClearClinicHeadDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, chCascade, "deputy_employee_id = ?", []any{employeeID}, now)
}

func cascadeClearClinicHeadDeputyInClinic(tx *gorm.DB, employeeID, clinicID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, chCascade, "deputy_employee_id = ? AND clinic_id = ?", []any{employeeID, clinicID}, now)
}

func cascadeRevokeOrgAdminAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, orgAdminCascade, "employee_id = ?", []any{employeeID}, now)
}

func cascadeClearOrgAdminDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, orgAdminCascade, "deputy_employee_id = ?", []any{employeeID}, now)
}

func cascadeRevokeOrgHeadAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, orgHeadCascade, "employee_id = ?", []any{employeeID}, now)
}

func cascadeClearOrgHeadDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, orgHeadCascade, "deputy_employee_id = ?", []any{employeeID}, now)
}

func cascadeRevokeOrgDispatcherAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeRevokeRolesBy(tx, orgDispatcherCascade, "employee_id = ?", []any{employeeID}, now)
}

func cascadeClearOrgDispatcherDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	return cascadeClearRoleDeputyBy(tx, orgDispatcherCascade, "deputy_employee_id = ?", []any{employeeID}, now)
}
