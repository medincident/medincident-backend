package authz

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

func (a *Authz) resolveOrgByClinic(ctx context.Context, clinicID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := a.db.WithContext(ctx).
		Raw(`SELECT organization_id FROM domain.clinics WHERE id = ?`, clinicID).
		Scan(&orgID).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			With("clinic_id", clinicID).Wrap(err)
	}
	if orgID == uuid.Nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public("Clinic not found.").
			With("clinic_id", clinicID).
			Errorf("clinic not found")
	}
	return orgID, nil
}

func (a *Authz) resolveOrgByDepartment(ctx context.Context, deptID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := a.db.WithContext(ctx).
		Raw(`SELECT c.organization_id FROM domain.departments d JOIN domain.clinics c ON c.id = d.clinic_id WHERE d.id = ?`, deptID).
		Scan(&orgID).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			With("department_id", deptID).Wrap(err)
	}
	if orgID == uuid.Nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public("Department not found.").
			With("department_id", deptID).
			Errorf("department not found")
	}
	return orgID, nil
}

func (a *Authz) resolveOrgByEmployee(ctx context.Context, empID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := a.db.WithContext(ctx).
		Raw(`SELECT organization_id FROM domain.employees WHERE id = ?`, empID).
		Scan(&orgID).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			With("employee_id", empID).Wrap(err)
	}
	if orgID == uuid.Nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public("Employee not found.").
			With("employee_id", empID).
			Errorf("employee not found")
	}
	return orgID, nil
}

func (a *Authz) resolveOrgByVacation(ctx context.Context, vacID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := a.db.WithContext(ctx).
		Raw(`SELECT e.organization_id FROM domain.employee_vacations v JOIN domain.employees e ON e.id = v.employee_id WHERE v.id = ?`, vacID).
		Scan(&orgID).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			With("vacation_id", vacID).Wrap(err)
	}
	if orgID == uuid.Nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public("Vacation not found.").
			With("vacation_id", vacID).
			Errorf("vacation not found")
	}
	return orgID, nil
}

func (a *Authz) resolveOrgByCategory(ctx context.Context, catID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := a.db.WithContext(ctx).
		Raw(`SELECT organization_id FROM domain.incident_categories WHERE id = ?`, catID).
		Scan(&orgID).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			With("category_id", catID).Wrap(err)
	}
	if orgID == uuid.Nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public("Category not found.").
			With("category_id", catID).
			Errorf("category not found")
	}
	return orgID, nil
}

func (a *Authz) resolveOrgByType(ctx context.Context, typeID uuid.UUID) (uuid.UUID, error) {
	var orgID uuid.UUID
	if err := a.db.WithContext(ctx).
		Raw(`SELECT organization_id FROM domain.incident_types WHERE id = ?`, typeID).
		Scan(&orgID).Error; err != nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			With("type_id", typeID).Wrap(err)
	}
	if orgID == uuid.Nil {
		return uuid.Nil, oops.In("service.authz").Code(ErrCodeScopeResolveFailed).
			Public("Type not found.").
			With("type_id", typeID).
			Errorf("type not found")
	}
	return orgID, nil
}
