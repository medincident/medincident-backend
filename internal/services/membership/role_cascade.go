package membership

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
)

// cascadeRevokeDepartmentResponsible revokes every DR role where the
// employee is the holder on the given department_id. For each revoked
// row:
//   - if it carries a deputy, publish DepartmentResponsibleDeputyRemoved first
//   - delete the row
//   - publish DepartmentResponsibleRevoked
func cascadeRevokeDepartmentResponsible(tx *gorm.DB, employeeID, departmentID uuid.UUID, now time.Time) error {
	var rows []model.DepartmentResponsible
	err := tx.Where("employee_id = ? AND department_id = ?", employeeID, departmentID).
		Find(&rows).Error
	if err != nil {
		return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			row.DepartmentID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}
		if err := publishDepartmentResponsibleRevoked(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeDepartmentResponsibleAll is the same as
// cascadeRevokeDepartmentResponsible, but unconstrained by
// department_id. Used by TerminateEmployee.
func cascadeRevokeDepartmentResponsibleAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.DepartmentResponsible
	err := tx.Where("employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishDepartmentResponsibleDeputyRemoved(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.DepartmentResponsible{}, "department_id = ? AND employee_id = ?",
			row.DepartmentID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleDeleteFailed).Wrap(err)
		}
		if err := publishDepartmentResponsibleRevoked(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearDepartmentResponsibleDeputy clears the deputy slot on
// every DR role where the given employee is the deputy. Used by
// TerminateEmployee to handle rows where the terminating employee was
// the deputy of someone else's role.
func cascadeClearDepartmentResponsibleDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.DepartmentResponsible
	err := tx.Where("deputy_employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if err := tx.Exec(
			`UPDATE domain.department_responsibles SET deputy_employee_id = NULL, updated_at = now() WHERE department_id = ? AND employee_id = ?`,
			row.DepartmentID, row.EmployeeID,
		).Error; err != nil {
			return oops.In(scopeDepartmentResponsible).Code(ErrCodeDepartmentResponsibleSaveFailed).Wrap(err)
		}
		if err := publishDepartmentResponsibleDeputyRemoved(tx, row.DepartmentID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeClinicHead revokes every CH role where the employee is
// the holder on the given clinic_id. For each revoked row:
//   - if it carries a deputy, publish ClinicHeadDeputyRemoved first
//   - delete the row
//   - publish ClinicHeadRevoked
func cascadeRevokeClinicHead(tx *gorm.DB, employeeID, clinicID uuid.UUID, now time.Time) error {
	var rows []model.ClinicHead
	err := tx.Where("employee_id = ? AND clinic_id = ?", employeeID, clinicID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishClinicHeadDeputyRemoved(tx, row.ClinicID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.ClinicHead{}, "clinic_id = ? AND employee_id = ?",
			row.ClinicID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadDeleteFailed).Wrap(err)
		}
		if err := publishClinicHeadRevoked(tx, row.ClinicID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeClinicHeadAll is the same as cascadeRevokeClinicHead,
// but unconstrained by clinic_id. Used by TerminateEmployee.
func cascadeRevokeClinicHeadAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.ClinicHead
	err := tx.Where("employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishClinicHeadDeputyRemoved(tx, row.ClinicID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.ClinicHead{}, "clinic_id = ? AND employee_id = ?",
			row.ClinicID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadDeleteFailed).Wrap(err)
		}
		if err := publishClinicHeadRevoked(tx, row.ClinicID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearClinicHeadDeputy clears the deputy slot on every CH
// role where the given employee is the deputy. Used by TerminateEmployee
// to handle rows where the terminating employee was the deputy of
// someone else's role.
func cascadeClearClinicHeadDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.ClinicHead
	err := tx.Where("deputy_employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if err := tx.Exec(
			`UPDATE domain.clinic_heads SET deputy_employee_id = NULL, updated_at = now() WHERE clinic_id = ? AND employee_id = ?`,
			row.ClinicID, row.EmployeeID,
		).Error; err != nil {
			return oops.In(scopeClinicHead).Code(ErrCodeClinicHeadSaveFailed).Wrap(err)
		}
		if err := publishClinicHeadDeputyRemoved(tx, row.ClinicID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeOrgAdminAll revokes every OrgAdmin role where the
// employee is the holder. Used by TerminateEmployee.
func cascadeRevokeOrgAdminAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.OrgAdmin
	err := tx.Where("employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishOrgAdminDeputyRemoved(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.OrgAdmin{}, "organization_id = ? AND employee_id = ?",
			row.OrganizationID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminDeleteFailed).Wrap(err)
		}
		if err := publishOrgAdminRevoked(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearOrgAdminDeputy clears the deputy slot on every OrgAdmin
// role where the given employee is the deputy. Used by TerminateEmployee
// to handle rows where the terminating employee was the deputy of
// someone else's role.
func cascadeClearOrgAdminDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.OrgAdmin
	err := tx.Where("deputy_employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if err := tx.Exec(
			`UPDATE domain.org_admins SET deputy_employee_id = NULL, updated_at = now() WHERE organization_id = ? AND employee_id = ?`,
			row.OrganizationID, row.EmployeeID,
		).Error; err != nil {
			return oops.In(scopeOrgAdmin).Code(ErrCodeOrganizationAdminSaveFailed).Wrap(err)
		}
		if err := publishOrgAdminDeputyRemoved(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeOrgHeadAll revokes every OrgHead role where the
// employee is the holder. Used by TerminateEmployee.
func cascadeRevokeOrgHeadAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.OrgHead
	err := tx.Where("employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishOrgHeadDeputyRemoved(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.OrgHead{}, "organization_id = ? AND employee_id = ?",
			row.OrganizationID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadDeleteFailed).Wrap(err)
		}
		if err := publishOrgHeadRevoked(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearOrgHeadDeputy clears the deputy slot on every OrgHead
// role where the given employee is the deputy. Used by TerminateEmployee
// to handle rows where the terminating employee was the deputy of
// someone else's role.
func cascadeClearOrgHeadDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.OrgHead
	err := tx.Where("deputy_employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if err := tx.Exec(
			`UPDATE domain.org_heads SET deputy_employee_id = NULL, updated_at = now() WHERE organization_id = ? AND employee_id = ?`,
			row.OrganizationID, row.EmployeeID,
		).Error; err != nil {
			return oops.In(scopeOrgHead).Code(ErrCodeOrganizationHeadSaveFailed).Wrap(err)
		}
		if err := publishOrgHeadDeputyRemoved(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeRevokeOrgDispatcherAll revokes every OrgDispatcher role where the
// employee is the holder. Used by TerminateEmployee.
func cascadeRevokeOrgDispatcherAll(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.OrgDispatcher
	err := tx.Where("employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if row.DeputyEmployeeID != nil {
			if err := publishOrgDispatcherDeputyRemoved(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.OrgDispatcher{}, "organization_id = ? AND employee_id = ?",
			row.OrganizationID, row.EmployeeID).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherDeleteFailed).Wrap(err)
		}
		if err := publishOrgDispatcherRevoked(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}

// cascadeClearOrgDispatcherDeputy clears the deputy slot on every OrgDispatcher
// role where the given employee is the deputy. Used by TerminateEmployee
// to handle rows where the terminating employee was the deputy of
// someone else's role.
func cascadeClearOrgDispatcherDeputy(tx *gorm.DB, employeeID uuid.UUID, now time.Time) error {
	var rows []model.OrgDispatcher
	err := tx.Where("deputy_employee_id = ?", employeeID).Find(&rows).Error
	if err != nil {
		return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherLoadFailed).Wrap(err)
	}
	for _, row := range rows {
		if err := tx.Exec(
			`UPDATE domain.org_dispatchers SET deputy_employee_id = NULL, updated_at = now() WHERE organization_id = ? AND employee_id = ?`,
			row.OrganizationID, row.EmployeeID,
		).Error; err != nil {
			return oops.In(scopeOrgDispatcher).Code(ErrCodeOrganizationDispatcherSaveFailed).Wrap(err)
		}
		if err := publishOrgDispatcherDeputyRemoved(tx, row.OrganizationID, row.EmployeeID, now); err != nil {
			return err
		}
	}
	return nil
}
