package projector

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
)

// ── Clinic Head ───────────────────────────────────────────────────────────

func ClinicHeadAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *clinicv1.ClinicHeadAssigned) error {
	clinicID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	assignedAt := ev.GetAssignedAt().AsTime()
	if err := tx.Exec(`INSERT INTO projections.clinic_heads (clinic_id, employee_id, deputy_employee_id, created_at, updated_at) VALUES (?, ?, NULL, ?, ?) ON CONFLICT (clinic_id) DO UPDATE SET employee_id = EXCLUDED.employee_id, updated_at = EXCLUDED.updated_at`,
		clinicID, empID, assignedAt, assignedAt).Error; err != nil {
		return wrapRole(err, "clinic_head", empID)
	}
	return nil
}

func ClinicHeadDeputyAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *clinicv1.ClinicHeadDeputyAssigned) error {
	clinicID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	deputyID := uuid.MustParse(ev.GetDeputyEmployeeId())
	updatedAt := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.clinic_heads SET deputy_employee_id = ?, updated_at = ? WHERE clinic_id = ? AND employee_id = ?`,
		deputyID, updatedAt, clinicID, empID).Error; err != nil {
		return wrapRole(err, "clinic_head", empID)
	}
	return nil
}

func ClinicHeadDeputyRemoved(tx *gorm.DB, aggregateID string, _ time.Time, ev *clinicv1.ClinicHeadDeputyRemoved) error {
	clinicID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	updatedAt := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.clinic_heads SET deputy_employee_id = NULL, updated_at = ? WHERE clinic_id = ? AND employee_id = ?`,
		updatedAt, clinicID, empID).Error; err != nil {
		return wrapRole(err, "clinic_head", empID)
	}
	return nil
}

func ClinicHeadRevoked(tx *gorm.DB, aggregateID string, _ time.Time, ev *clinicv1.ClinicHeadRevoked) error {
	clinicID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	if err := tx.Exec(`DELETE FROM projections.clinic_heads WHERE clinic_id = ? AND employee_id = ?`,
		clinicID, empID).Error; err != nil {
		return wrapRole(err, "clinic_head", empID)
	}
	return nil
}

// ── Department Responsible ────────────────────────────────────────────────

func DeptResponsibleAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *deptv1.DeptResponsibleAssigned) error {
	deptID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetAssignedAt().AsTime()
	if err := tx.Exec(`INSERT INTO projections.department_responsibles (department_id, employee_id, deputy_employee_id, created_at, updated_at) VALUES (?, ?, NULL, ?, ?) ON CONFLICT (department_id) DO UPDATE SET employee_id = EXCLUDED.employee_id, updated_at = EXCLUDED.updated_at`,
		deptID, empID, at, at).Error; err != nil {
		return wrapRole(err, "dept_responsible", empID)
	}
	return nil
}

func DeptResponsibleDeputyAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *deptv1.DeptResponsibleDeputyAssigned) error {
	deptID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	depID := uuid.MustParse(ev.GetDeputyEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.department_responsibles SET deputy_employee_id = ?, updated_at = ? WHERE department_id = ? AND employee_id = ?`,
		depID, at, deptID, empID).Error; err != nil {
		return wrapRole(err, "dept_responsible", empID)
	}
	return nil
}

func DeptResponsibleDeputyRemoved(tx *gorm.DB, aggregateID string, _ time.Time, ev *deptv1.DeptResponsibleDeputyRemoved) error {
	deptID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.department_responsibles SET deputy_employee_id = NULL, updated_at = ? WHERE department_id = ? AND employee_id = ?`,
		at, deptID, empID).Error; err != nil {
		return wrapRole(err, "dept_responsible", empID)
	}
	return nil
}

func DeptResponsibleRevoked(tx *gorm.DB, aggregateID string, _ time.Time, ev *deptv1.DeptResponsibleRevoked) error {
	deptID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	if err := tx.Exec(`DELETE FROM projections.department_responsibles WHERE department_id = ? AND employee_id = ?`,
		deptID, empID).Error; err != nil {
		return wrapRole(err, "dept_responsible", empID)
	}
	return nil
}

// ── Org Admin ─────────────────────────────────────────────────────────────

func OrgAdminAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgAdminAssigned) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetAssignedAt().AsTime()
	if err := tx.Exec(`INSERT INTO projections.org_admins (organization_id, employee_id, deputy_employee_id, created_at, updated_at) VALUES (?, ?, NULL, ?, ?) ON CONFLICT (organization_id) DO UPDATE SET employee_id = EXCLUDED.employee_id, updated_at = EXCLUDED.updated_at`,
		orgID, empID, at, at).Error; err != nil {
		return wrapRole(err, "org_admin", empID)
	}
	return nil
}

func OrgAdminDeputyAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgAdminDeputyAssigned) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	depID := uuid.MustParse(ev.GetDeputyEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.org_admins SET deputy_employee_id = ?, updated_at = ? WHERE organization_id = ? AND employee_id = ?`,
		depID, at, orgID, empID).Error; err != nil {
		return wrapRole(err, "org_admin", empID)
	}
	return nil
}

func OrgAdminDeputyRemoved(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgAdminDeputyRemoved) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.org_admins SET deputy_employee_id = NULL, updated_at = ? WHERE organization_id = ? AND employee_id = ?`,
		at, orgID, empID).Error; err != nil {
		return wrapRole(err, "org_admin", empID)
	}
	return nil
}

func OrgAdminRevoked(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgAdminRevoked) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	if err := tx.Exec(`DELETE FROM projections.org_admins WHERE organization_id = ? AND employee_id = ?`,
		orgID, empID).Error; err != nil {
		return wrapRole(err, "org_admin", empID)
	}
	return nil
}

// ── Org Head ──────────────────────────────────────────────────────────────

func OrgHeadAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgHeadAssigned) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetAssignedAt().AsTime()
	if err := tx.Exec(`INSERT INTO projections.org_heads (organization_id, employee_id, deputy_employee_id, created_at, updated_at) VALUES (?, ?, NULL, ?, ?) ON CONFLICT (organization_id) DO UPDATE SET employee_id = EXCLUDED.employee_id, updated_at = EXCLUDED.updated_at`,
		orgID, empID, at, at).Error; err != nil {
		return wrapRole(err, "org_head", empID)
	}
	return nil
}

func OrgHeadDeputyAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgHeadDeputyAssigned) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	depID := uuid.MustParse(ev.GetDeputyEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.org_heads SET deputy_employee_id = ?, updated_at = ? WHERE organization_id = ? AND employee_id = ?`,
		depID, at, orgID, empID).Error; err != nil {
		return wrapRole(err, "org_head", empID)
	}
	return nil
}

func OrgHeadDeputyRemoved(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgHeadDeputyRemoved) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.org_heads SET deputy_employee_id = NULL, updated_at = ? WHERE organization_id = ? AND employee_id = ?`,
		at, orgID, empID).Error; err != nil {
		return wrapRole(err, "org_head", empID)
	}
	return nil
}

func OrgHeadRevoked(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgHeadRevoked) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	if err := tx.Exec(`DELETE FROM projections.org_heads WHERE organization_id = ? AND employee_id = ?`,
		orgID, empID).Error; err != nil {
		return wrapRole(err, "org_head", empID)
	}
	return nil
}

// ── Org Dispatcher ────────────────────────────────────────────────────────

func OrgDispatcherAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgDispatcherAssigned) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetAssignedAt().AsTime()
	if err := tx.Exec(`INSERT INTO projections.org_dispatchers (organization_id, employee_id, deputy_employee_id, created_at, updated_at) VALUES (?, ?, NULL, ?, ?) ON CONFLICT (organization_id) DO UPDATE SET employee_id = EXCLUDED.employee_id, updated_at = EXCLUDED.updated_at`,
		orgID, empID, at, at).Error; err != nil {
		return wrapRole(err, "org_dispatcher", empID)
	}
	return nil
}

func OrgDispatcherDeputyAssigned(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgDispatcherDeputyAssigned) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	depID := uuid.MustParse(ev.GetDeputyEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.org_dispatchers SET deputy_employee_id = ?, updated_at = ? WHERE organization_id = ? AND employee_id = ?`,
		depID, at, orgID, empID).Error; err != nil {
		return wrapRole(err, "org_dispatcher", empID)
	}
	return nil
}

func OrgDispatcherDeputyRemoved(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgDispatcherDeputyRemoved) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	at := ev.GetUpdatedAt().AsTime()
	if err := tx.Exec(`UPDATE projections.org_dispatchers SET deputy_employee_id = NULL, updated_at = ? WHERE organization_id = ? AND employee_id = ?`,
		at, orgID, empID).Error; err != nil {
		return wrapRole(err, "org_dispatcher", empID)
	}
	return nil
}

func OrgDispatcherRevoked(tx *gorm.DB, aggregateID string, _ time.Time, ev *orgv1.OrgDispatcherRevoked) error {
	orgID := uuid.MustParse(aggregateID)
	empID := uuid.MustParse(ev.GetEmployeeId())
	if err := tx.Exec(`DELETE FROM projections.org_dispatchers WHERE organization_id = ? AND employee_id = ?`,
		orgID, empID).Error; err != nil {
		return wrapRole(err, "org_dispatcher", empID)
	}
	return nil
}

// ── System Admin ──────────────────────────────────────────────────────────

func SystemAdminGranted(tx *gorm.DB, aggregateID string, _ time.Time, ev *sav1.SystemAdminGranted) error {
	grantedAt := ev.GetGrantedAt().AsTime()
	if err := tx.Exec(`INSERT INTO projections.system_admins (zitadel_user_id, created_at, updated_at) VALUES (?, ?, ?) ON CONFLICT (zitadel_user_id) DO NOTHING`,
		aggregateID, grantedAt, grantedAt).Error; err != nil {
		return oops.In("projector.system_admin").Code(ErrCodeRoleProjectionFailed).With("zitadel_user_id", aggregateID).Wrap(err)
	}
	return nil
}

func SystemAdminRevoked(tx *gorm.DB, aggregateID string, _ time.Time, _ *sav1.SystemAdminRevoked) error {
	if err := tx.Exec(`DELETE FROM projections.system_admins WHERE zitadel_user_id = ?`, aggregateID).Error; err != nil {
		return oops.In("projector.system_admin").Code(ErrCodeRoleProjectionFailed).With("zitadel_user_id", aggregateID).Wrap(err)
	}
	return nil
}

func wrapRole(err error, kind string, empID uuid.UUID) error {
	return oops.In("projector.role."+kind).
		Code(ErrCodeRoleProjectionFailed).
		With("employee_id", empID).
		Wrap(err)
}

// Forwarding methods on *Projectors

func (p *Projectors) ClinicHeadAssigned(tx *gorm.DB, id string, t time.Time, ev *clinicv1.ClinicHeadAssigned) error {
	return ClinicHeadAssigned(tx, id, t, ev)
}

func (p *Projectors) ClinicHeadDeputyAssigned(tx *gorm.DB, id string, t time.Time, ev *clinicv1.ClinicHeadDeputyAssigned) error {
	return ClinicHeadDeputyAssigned(tx, id, t, ev)
}

func (p *Projectors) ClinicHeadDeputyRemoved(tx *gorm.DB, id string, t time.Time, ev *clinicv1.ClinicHeadDeputyRemoved) error {
	return ClinicHeadDeputyRemoved(tx, id, t, ev)
}

func (p *Projectors) ClinicHeadRevoked(tx *gorm.DB, id string, t time.Time, ev *clinicv1.ClinicHeadRevoked) error {
	return ClinicHeadRevoked(tx, id, t, ev)
}

func (p *Projectors) DeptResponsibleAssigned(tx *gorm.DB, id string, t time.Time, ev *deptv1.DeptResponsibleAssigned) error {
	return DeptResponsibleAssigned(tx, id, t, ev)
}

func (p *Projectors) DeptResponsibleDeputyAssigned(tx *gorm.DB, id string, t time.Time, ev *deptv1.DeptResponsibleDeputyAssigned) error {
	return DeptResponsibleDeputyAssigned(tx, id, t, ev)
}

func (p *Projectors) DeptResponsibleDeputyRemoved(tx *gorm.DB, id string, t time.Time, ev *deptv1.DeptResponsibleDeputyRemoved) error {
	return DeptResponsibleDeputyRemoved(tx, id, t, ev)
}

func (p *Projectors) DeptResponsibleRevoked(tx *gorm.DB, id string, t time.Time, ev *deptv1.DeptResponsibleRevoked) error {
	return DeptResponsibleRevoked(tx, id, t, ev)
}

func (p *Projectors) OrgAdminAssigned(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgAdminAssigned) error {
	return OrgAdminAssigned(tx, id, t, ev)
}

func (p *Projectors) OrgAdminDeputyAssigned(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgAdminDeputyAssigned) error {
	return OrgAdminDeputyAssigned(tx, id, t, ev)
}

func (p *Projectors) OrgAdminDeputyRemoved(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgAdminDeputyRemoved) error {
	return OrgAdminDeputyRemoved(tx, id, t, ev)
}

func (p *Projectors) OrgAdminRevoked(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgAdminRevoked) error {
	return OrgAdminRevoked(tx, id, t, ev)
}

func (p *Projectors) OrgHeadAssigned(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgHeadAssigned) error {
	return OrgHeadAssigned(tx, id, t, ev)
}

func (p *Projectors) OrgHeadDeputyAssigned(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgHeadDeputyAssigned) error {
	return OrgHeadDeputyAssigned(tx, id, t, ev)
}

func (p *Projectors) OrgHeadDeputyRemoved(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgHeadDeputyRemoved) error {
	return OrgHeadDeputyRemoved(tx, id, t, ev)
}

func (p *Projectors) OrgHeadRevoked(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgHeadRevoked) error {
	return OrgHeadRevoked(tx, id, t, ev)
}

func (p *Projectors) OrgDispatcherAssigned(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgDispatcherAssigned) error {
	return OrgDispatcherAssigned(tx, id, t, ev)
}

func (p *Projectors) OrgDispatcherDeputyAssigned(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgDispatcherDeputyAssigned) error {
	return OrgDispatcherDeputyAssigned(tx, id, t, ev)
}

func (p *Projectors) OrgDispatcherDeputyRemoved(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgDispatcherDeputyRemoved) error {
	return OrgDispatcherDeputyRemoved(tx, id, t, ev)
}

func (p *Projectors) OrgDispatcherRevoked(tx *gorm.DB, id string, t time.Time, ev *orgv1.OrgDispatcherRevoked) error {
	return OrgDispatcherRevoked(tx, id, t, ev)
}

func (p *Projectors) SystemAdminGranted(tx *gorm.DB, id string, t time.Time, ev *sav1.SystemAdminGranted) error {
	return SystemAdminGranted(tx, id, t, ev)
}

func (p *Projectors) SystemAdminRevoked(tx *gorm.DB, id string, t time.Time, ev *sav1.SystemAdminRevoked) error {
	return SystemAdminRevoked(tx, id, t, ev)
}
